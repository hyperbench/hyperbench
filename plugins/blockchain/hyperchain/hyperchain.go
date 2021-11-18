package hyperchain

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/meshplus/gosdk/abi"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/hvm"
	"github.com/meshplus/gosdk/rpc"
	"github.com/meshplus/gosdk/utils/java"
	fcom "github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Client the implementation of  client.Blockchain
////based on hyperchain/flato network
type Client struct {
	*base.BlockchainBase
	rpcClient *rpc.RPC
	am        *AccountManager
	op        option
	contract  *Contract
}

// option means the the options of hyperchain client
type option struct {
	poll           bool
	simulate       bool
	defaultAccount string
	fakeSign       bool
	nonce          int64
	extraIDStr     []string
	extraIDInt64   []int64
}

// NewClient use given blockchainBase create Client.
func NewClient(blockchainBase *base.BlockchainBase) (client *Client, err error) {

	keystorePath := cast.ToString(blockchainBase.Options["keystore"])
	keystoreType := cast.ToString(blockchainBase.Options["sign"])
	poll := cast.ToBool(blockchainBase.Options["poll"])
	am := NewAccountManager(keystorePath, keystoreType, blockchainBase.Logger)
	rpcClient := rpc.NewRPCWithPath(blockchainBase.ConfigPath)

	client = &Client{
		BlockchainBase: blockchainBase,
		am:             am,
		rpcClient:      rpcClient,
		op: option{
			nonce: -1,
			poll:  poll,
		},
	}
	return
}

func convert(m map[interface{}]interface{}) []interface{} {
	ret := make([]interface{}, 0, len(m))
	// hint that lua index starts from 1
	for i := 1; i <= len(m); i++ {
		val, exist := m[float64(i)]
		if !exist {
			break
		}
		switch o := val.(type) {
		case map[interface{}]interface{}:
			ret = append(ret, convert(o))
		case string:
			ret = append(ret, val)
		}
	}
	return ret
}

//Invoke invoke contract with funcName and args in hyperchain network
func (c *Client) Invoke(invoke bcom.Invoke, ops ...bcom.Option) *fcom.Result {
	funcName, args := invoke.Func, invoke.Args
	for idx, arg := range args {
		if m, ok := arg.(map[interface{}]interface{}); ok {
			args[idx] = convert(m)
		}
	}
	var (
		payload []byte
		err     error
	)

	if c.contract == nil {
		return &fcom.Result{}
	}
	buildTime := time.Now().UnixNano()

	switch c.contract.VM {
	case rpc.EVM:
		c.Logger.Debugf("invoke evm contract funcName: %v, param: %v", funcName, args)

		payload, err = c.contract.ABI.Encode(funcName, args...)
		if err != nil {
			c.Logger.Errorf("abi %v can not pack param: %v", c.contract.ABI, err)
			return &fcom.Result{
				Label:     funcName,
				UID:       fcom.InvalidUID,
				Ret:       []interface{}{},
				Status:    fcom.Failure,
				BuildTime: buildTime,
			}
		}
	case rpc.JVM:
		var argStrings = make([]string, len(args))
		for idx, arg := range args {
			argStrings[idx] = fmt.Sprint(arg)
		}
		c.Logger.Debugf("invoke evm contract funcName: %v, param: %v", funcName, argStrings)
		payload = java.EncodeJavaFunc(funcName, argStrings...)
	case rpc.HVM:
		var beanAbi *hvm.BeanAbi
		beanAbi, err = c.contract.hvmABI.GetBeanAbi(funcName)
		if err != nil {
			c.Logger.Info(err)
			return &fcom.Result{
				Label:     funcName,
				UID:       fcom.InvalidUID,
				Ret:       []interface{}{},
				Status:    fcom.Failure,
				BuildTime: buildTime,
			}
		}
		payload, err = hvm.GenPayload(beanAbi, args...)
		if err != nil {
			c.Logger.Info(err)
			return &fcom.Result{
				Label:     funcName,
				UID:       fcom.InvalidUID,
				Ret:       []interface{}{},
				Status:    fcom.Failure,
				BuildTime: buildTime,
			}
		}
	}

	// invoke
	ac, err := c.am.GetAccount(c.op.defaultAccount)
	if err != nil {
		return &fcom.Result{
			Label:     funcName,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}

	tranInvoke := rpc.NewTransaction(ac.GetAddress().Hex()).Invoke(c.contract.Addr, payload).VMType(c.contract.VM).Simulate(c.op.simulate)
	if c.op.nonce >= 0 {
		tranInvoke.SetNonce(c.op.nonce)
	}

	c.sign(tranInvoke, ac)

	// just send tx after sending tx
	hash, stdErr := c.rpcClient.InvokeContractReturnHash(tranInvoke)
	sendTime := time.Now().UnixNano()
	if stdErr != nil {
		c.Logger.Infof("invoke error: %v", stdErr)
		return &fcom.Result{
			Label:     funcName,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
			SendTime:  sendTime,
		}
	}

	ret := &fcom.Result{
		Label:     funcName,
		UID:       hash,
		Ret:       []interface{}{},
		Status:    fcom.Success,
		BuildTime: buildTime,
		SendTime:  sendTime,
	}
	if !c.op.poll {
		return ret
	}
	return c.Confirm(ret)

}

// Confirm check the result of `Invoke` or `Transfer`
func (c *Client) Confirm(result *fcom.Result, ops ...bcom.Option) *fcom.Result {

	if result.UID == "" ||
		result.UID == fcom.InvalidUID ||
		result.Status != fcom.Success ||
		result.Label == fcom.InvalidLabel {
		return result
	}

	// poll
	txReceipt, stdErr, got := c.rpcClient.GetTxReceiptByPolling(result.UID, false)
	result.ConfirmTime = time.Now().UnixNano()
	if stdErr != nil || !got {
		c.Logger.Errorf("invoke failed: %v", stdErr)
		result.Status = fcom.Unknown
		return result
	}

	result.Status = fcom.Confirm
	var results []interface{}
	if result.Label == fcom.BuiltinTransferLabel {
		result.Ret = []interface{}{txReceipt.Ret}
		return result
	}
	// decode result
	switch c.contract.VM {
	case rpc.EVM:
		c.Logger.Debugf("error: %v", txReceipt)
		decodeResult, err := c.contract.ABI.Decode(result.Label, common.FromHex(txReceipt.Ret))
		if err != nil {
			c.Logger.Noticef("decode error: %v, result hex: %v,result: %v", err, txReceipt.Ret, common.FromHex(txReceipt.Ret))
			return result
		}
		if array, ok := decodeResult.([]interface{}); ok { // multiple return value
			results = array
		} else { // single return value
			results = append(results, decodeResult)
		}

	case rpc.JVM, rpc.HVM:
		results = append(results, java.DecodeJavaResult(txReceipt.Ret))
	default:
		results = append(results, txReceipt.Ret)
	}

	result.Ret = results
	info, stdErr := c.rpcClient.GetTransactionByHash(txReceipt.TxHash)
	if stdErr != nil {
		c.Logger.Infof("get transaction by hash error: %v", stdErr)
		return result
	}
	result.WriteTime = info.BlockWriteTime
	return result
}

func (c *Client) sign(tx *rpc.Transaction, acc Account) {
	if c.op.fakeSign {
		tx.SetSignature(fakeSign())
	} else {
		switch c.am.AccountType {
		case ECDSA:
			tx.SignWithClang(acc)
		case SM2:
			tx.Sign(acc)
		}
	}
}

//Transfer transfer a amount of money from a account to the other one
func (c *Client) Transfer(args bcom.Transfer, ops ...bcom.Option) (result *fcom.Result) {

	//c.Logger.Notice("transfer")
	//defer func() {
	//	_ = recover()
	//	var buf [4096]byte
	//	n := runtime.Stack(buf[:], false)
	//	c.Logger.Critical("==> %s\n", string(buf[:n]))
	//}()

	from, to, amount, extra := args.From, args.To, args.Amount, args.Extra
	buildTime := time.Now().UnixNano()
	fromAcc, err := c.am.GetAccount(from)
	if err != nil {
		return &fcom.Result{
			Label:     fcom.BuiltinTransferLabel,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}
	toAcc, err := c.am.GetAccount(to)
	if err != nil {
		return &fcom.Result{
			Label:     fcom.BuiltinTransferLabel,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}

	tx := rpc.NewTransaction(fromAcc.GetAddress().Hex()).Transfer(toAcc.GetAddress().Hex(), amount).Extra(extra).Simulate(c.op.simulate)
	if len(c.op.extraIDStr) > 0 {
		tx.SetExtraIDString(c.op.extraIDStr...)
	}
	if len(c.op.extraIDInt64) > 0 {
		tx.SetExtraIDInt64(c.op.extraIDInt64...)
	}

	if c.op.nonce >= 0 {
		tx.SetNonce(c.op.nonce)
	}

	c.sign(tx, fromAcc)
	hash, stdErr := c.rpcClient.SendTxReturnHash(tx)
	sendTime := time.Now().UnixNano()
	if stdErr != nil {
		c.Logger.Infof("transfer error: %v", stdErr)
		return &fcom.Result{
			Label:     fcom.BuiltinTransferLabel,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
			SendTime:  sendTime,
		}
	}
	ret := &fcom.Result{
		Label:     fcom.BuiltinTransferLabel,
		UID:       hash,
		Ret:       []interface{}{},
		Status:    fcom.Success,
		BuildTime: buildTime,
		SendTime:  sendTime,
	}

	if !c.op.poll {
		return ret
	}
	return c.Confirm(result)
}

//SetContext set test group context in go client
func (c *Client) SetContext(context string) error {
	c.Logger.Debugf("prepare msg: %v", context)
	msg := &Msg{}
	var (
		err error
	)

	if context == "" {
		c.Logger.Infof("Prepare nothing")
		return nil
	}

	err = json.Unmarshal([]byte(context), msg)
	if err != nil {
		c.Logger.Errorf("can not unmarshal msg: %v \n err: %v", context, err)
		return err
	}

	// set contract context
	contract := &Contract{
		ContractRaw: msg.Contract,
	}
	switch msg.Contract.VM {
	case rpc.EVM:
		a, err := abi.JSON(strings.NewReader(msg.Contract.ABIRaw))
		if err != nil {
			c.Logger.Errorf("can not parse abi: %v \n err: %v", contract.ABIRaw, err)
			return err
		}
		contract.ABI = a
	case rpc.JVM:
	case rpc.HVM:
		a, err := hvm.GenAbi(msg.Contract.ABIRaw)
		if err != nil {
			return err
		}
		contract.hvmABI = a
	default:
	}
	c.contract = contract

	// set account context
	for acName, ac := range msg.Accounts {
		_, _ = c.am.SetAccount(acName, ac, PASSWORD)
	}

	return nil
}

//ResetContext reset test group context in go client
func (c *Client) ResetContext() error {
	return nil
}

//GetContext generate TxContext
func (c *Client) GetContext() (string, error) {
	var (
		bytes []byte
		err   error
	)
	if c.contract == nil || c.am == nil {
		return "", nil
	}

	msg := Msg{
		Contract: c.contract.ContractRaw,
	}

	bytes, err = json.Marshal(msg)

	return string(bytes), err
}

//Statistic statistic remote node performance
func (c *Client) Statistic(statistic bcom.Statistic) (*fcom.RemoteStatistic, error) {
	from, to := statistic.From, statistic.To
	tps, err := c.rpcClient.QueryTPS(uint64(from), uint64(to))
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	ret := &fcom.RemoteStatistic{
		Start:    from,
		End:      to,
		BlockNum: int(tps.TotalBlockNum),
		TxNum:    int(tps.Tps) * int(to-from) / int(time.Second),
	}

	return ret, nil
}

// Option hyperchain receive options to change the config to client.
// Supported Options:
// 1. key: confirm
//    valueType: bool
//    effect: set confirm true will let client poll for receipt after sending transaction
//            set confirm false will let client return immediately after sending transaction
//    default: default value is setting by the `benchmark.confirm` in testplan
// 2. key: simulate
//    valueType: bool
//    effect: set simulate true will let client send simulate transaction
//            set simulate false will let client send common transaction
//    default: false
// 3. key: account
//    value: account
//    effect: use the account to invoke contract
//    default:  account aliased as '0'
// 4. key: nonce
//    value: float64
//    effect: if nonce is non-negative, it will be set to transaction's `nonce` field
//    default: -1
func (c *Client) Option(options bcom.Option) error {
	for key, value := range options {
		switch key {
		case "confirm":
			if poll, ok := value.(bool); ok {
				c.op.poll = poll
			} else {
				return errors.Errorf("option `confirm` type error: %v", reflect.TypeOf(value).Name())
			}
		case "simulate":
			if simulate, ok := value.(bool); ok {
				c.op.simulate = simulate
			} else {
				return errors.Errorf("option `simulate` type error: %v", reflect.TypeOf(value).Name())
			}
		case "account":
			if a, ok := value.(string); ok {
				c.op.defaultAccount = a
			} else {
				return errors.Errorf("option `account` type error: %v", reflect.TypeOf(value).Name())
			}
		case "nonce":
			if n, ok := value.(float64); ok {
				c.op.nonce = int64(n)
			}
		case "extraid":
			if n, ok := value.([]interface{}); ok {
				var strs = make([]string, 0, len(n))
				var ints = make([]int64, 0, len(n))
				for _, v := range n {
					switch o := v.(type) {
					case string:
						strs = append(strs, o)
					case float64:
						ints = append(ints, int64(o))
					}
				}

				c.op.extraIDStr = strs
				c.op.extraIDInt64 = ints
			}
		}
	}
	return nil
}
