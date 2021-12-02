package fabric

import (
	"encoding/json"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	common2 "github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	//DefaultConf the default config file name
	DefaultConf = "config.yaml"
)

//Fabric the implementation of  client.Blockchain
//based on fabric network
type Fabric struct {
	*base.BlockchainBase
	SDK        *SDK
	ChannelID  string
	CCId       string
	CCPath     string
	StartBlock uint64
	//OrgMSPId 		string
	ShareAccount   int
	Instant        int
	InitArgs       [][]byte
	AccountManager *ClientManager
	MSP            bool
	invoke         bool
}

//Msg contains message of context
type Msg struct {
	CCId       string
	Accounts   map[string]*Client
	StartBlock uint64
}

// New use given blockchainBase create Fabric.
func New(blockchainBase *base.BlockchainBase) (client *Fabric, err error) {
	client = &Fabric{
		BlockchainBase: blockchainBase,
	}

	client.Instant = cast.ToInt(client.Options["instant"])
	client.SDK = NewSDK(blockchainBase, filepath.Join(client.ConfigPath, DefaultConf))
	client.ChannelID = cast.ToString(client.Options["channel"])
	client.CCPath = client.ContractPath

	initArgs := cast.ToStringSlice(client.Args)
	client.InitArgs = make([][]byte, 0, len(initArgs))
	for _, arg := range initArgs {
		client.InitArgs = append(client.InitArgs, []byte(arg))
	}
	client.MSP = cast.ToBool(client.Options["MSP"])
	client.invoke = true
	return
}

////New create a Fabric with config
//func (f *Fabric) New(config string) error {
//	//logger.Debugf("config: %v", config)
//	confParser := viper.New()
//	confParser.SetConfigType("toml")
//	if err := confParser.ReadConfig(bytes.NewBufferString(config)); err != nil {
//		panic(err)
//	}
//	f.CCPath = confParser.GetString("benchmark.contract")
//	f.Polling = confParser.GetBool("benchmark.confirm")
//
//	logger.Debugf("Contract:%v", f.CCPath)
//	configPath := confParser.GetString("network.config")
//	f.Instant = confParser.GetInt("benchmark.instant")
//	f.SDK = NewSDK(configPath + "/" + DefaultConfRelPath)
//	f.ChannelID = confParser.GetString("option.channel")
//	initArgs := confParser.GetStringSlice("option.initArgs")
//	for _, arg := range initArgs {
//		f.InitArgs = append(f.InitArgs, []byte(arg))
//	}
//	f.ShareAccount = confParser.GetInt("benchmark.user")
//	f.MSP = confParser.GetBool("option.MSP")
//	f.invoke = true
//	return nil
//}

// DeployContract deploy contract to fabric network
func (f *Fabric) DeployContract() error {
	//install chaincode
	ccID := strconv.Itoa(int(time.Now().UnixNano()))
	ccVersion := "0"
	_, err := InstallCC(f.CCPath, ccID, ccVersion, f.SDK.GetResmgmtClient())
	if err != nil {
		return err
	}

	//instantiate chaincode
	ccPolicy := cauthdsl.SignedByAnyMember(f.SDK.MspIds)
	_, err = InstantiateCC(f.CCPath, ccID, ccVersion, f.ChannelID, f.InitArgs, ccPolicy, f.SDK.GetResmgmtClient())
	if err != nil {
		return err
	}
	f.CCId = ccID
	return nil
}

// Option Fabric does not need now
func (f *Fabric) Option(option bcom.Option) error {
	if mode, ok := option["mode"]; ok {
		if mode == "query" {
			f.invoke = false
		} else {
			f.invoke = true
		}
	}
	return nil
}

// Invoke invoke contract with funcName and args in fabric network
func (f *Fabric) Invoke(invoke bcom.Invoke, ops ...bcom.Option) *common2.Result {
	funcName := invoke.Func
	args := invoke.Args
	intn := rand.Intn(len(f.AccountManager.Clients))
	account, e := f.AccountManager.GetAccount(strconv.Itoa(intn))
	var channelClient *channel.Client
	if e != nil {
		f.Logger.Error(e)
		channelClient = f.SDK.GetChannelClient(f.ChannelID, f.SDK.OrgAdmin, f.SDK.OrgName)
	} else {
		channelClient = f.SDK.GetChannelClient(f.ChannelID, account.Name, account.OrgName)
	}

	bytesArgs := make([][]byte, len(args))
	for i, arg := range args {
		s := arg.(string)
		bytesArgs[i] = []byte(s)
	}
	startTime := time.Now().UnixNano()
	resp, err := ExecuteCC(channelClient, f.CCId, funcName, bytesArgs, f.SDK.EndPoints, f.invoke)
	endTime := time.Now().UnixNano()
	if err != nil {
		return &common2.Result{
			UID:       common2.InvalidUID,
			Ret:       []interface{}{},
			Status:    common2.Failure,
			BuildTime: startTime,
			SendTime:  endTime,
		}
	}

	result := &common2.Result{
		UID:       string(resp.TransactionID),
		Ret:       []interface{}{resp.Payload},
		Status:    common2.Success,
		BuildTime: startTime,
		SendTime:  endTime,
	}

	return result

}

//GetContext generate context for fabric client
func (f *Fabric) GetContext() (string, error) {
	am, err := NewClientManager(f.SDK, f.MSP, f.Logger)
	if err != nil {
		f.Logger.Error("new client manager error. ", err)
		return "", err
	}
	f.AccountManager = am
	e := f.AccountManager.InitAccount(f.Instant)
	if e != nil {
		return "", e
	}
	info, err := f.SDK.GetLedgerClient(f.ChannelID, f.SDK.OrgAdmin, f.SDK.OrgName).QueryInfo()
	if err != nil {
		return "", err
	}
	msg := &Msg{
		CCId:       f.CCId,
		Accounts:   f.AccountManager.Clients,
		StartBlock: info.BCI.Height,
	}
	marshal, e := json.Marshal(msg)
	return string(marshal), e
}

//SetContext set context to each fabric client in VM
func (f *Fabric) SetContext(context string) error {
	am, err := NewClientManager(f.SDK, f.MSP, f.Logger)
	if err != nil {
		f.Logger.Error("new client manager error. ", err)
		return err
	}
	f.AccountManager = am

	msg := &Msg{}
	err = json.Unmarshal([]byte(context), msg)
	if err != nil {
		f.Logger.Errorf("can not unmarshal msg: %v \n err: %v", context, err)
		return err
	}
	f.AccountManager.Clients = msg.Accounts
	f.CCId = msg.CCId
	f.StartBlock = msg.StartBlock
	return nil
}

//ResetContext reset context
func (f *Fabric) ResetContext() error {
	return nil
}

//Statistic statistic node performance
func (f *Fabric) Statistic(statistic bcom.Statistic) (*common2.RemoteStatistic, error) {
	from, to := statistic.From, statistic.To
	statisticData, err := GetTPS(f.SDK.GetLedgerClient(f.ChannelID, f.SDK.OrgAdmin, f.SDK.OrgName), f.StartBlock, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	return statisticData, nil
}

//String serial fabric to string
func (f *Fabric) String() string {
	marshal, _ := json.Marshal(f)
	return string(marshal)
}
