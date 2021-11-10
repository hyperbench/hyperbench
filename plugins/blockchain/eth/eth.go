package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"io/ioutil"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	fcom "github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"

	"github.com/spf13/viper"
)

//Contract contains the abi and bin files of contract
type Contract struct {
	ABI string
	BIN string
}

//ETH the client of eth
type ETH struct {
	*base.BlockchainBase
	ethClient       *ethclient.Client
	privateKey      *ecdsa.PrivateKey
	auth            *bind.TransactOpts
	contractAddress common.Address
	startBlock      uint64
	contract        *Contract
}

//Msg contains message of context
type Msg struct {
	ContractAddress common.Address
	StartBlock      uint64
}

// New use given blockchainBase create ETH.
func New(blockchainBase *base.BlockchainBase) (client *ETH, err error) {
	ethClient, err := ethclient.Dial(viper.GetString(fcom.ClientConfigPath) + "/geth.ipc")
	if err != nil {
		return nil, err
	}

	privKey, _, err := KeystoreToPrivateKey(viper.GetString(fcom.ClientConfigPath)+"/keystore/"+viper.GetString(fcom.ClientAccount), "")
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	chainID, _ := ethClient.NetworkID(context.Background())
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	startBlock, err := ethClient.HeaderByNumber(context.Background(), nil)
	contract, _ := newContract()
	client = &ETH{
		BlockchainBase: blockchainBase,
		ethClient:      ethClient,
		privateKey:     privateKey,
		auth:           auth,
		startBlock:     startBlock.Number.Uint64(),
		contract:       contract,
	}

	return
}
func (e *ETH) DeployContract() error {
	parsed, err := abi.JSON(strings.NewReader(e.contract.ABI))
	if err != nil {
		return err
	}
	input := "1.0"
	contractAddress, tx, instance, err := bind.DeployContract(e.auth, parsed, common.FromHex(e.contract.BIN), e.ethClient, input)
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.contractAddress = contractAddress
	e.Logger.Info("contractAddress:" + contractAddress.Hex())
	e.Logger.Info("txHash:" + tx.Hash().Hex())
	_ = instance

	return nil
}

//Invoke invoke contract with funcName and args in eth network
func (e *ETH) Invoke(invoke bcom.Invoke, ops ...bcom.Option) *fcom.Result {
	parsed, err := abi.JSON(strings.NewReader(e.contract.ABI))
	if err != nil {
		e.Logger.Error(err)
		return nil
	}
	instance := bind.NewBoundContract(e.contractAddress, parsed, e.ethClient, e.ethClient, e.ethClient)
	publicKey := e.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		e.Logger.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		e.Logger.Error(err)
	}

	gasPrice, err := e.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		e.Logger.Error(err)
	}
	chainID, _ := e.ethClient.NetworkID(context.Background())
	auth, _ := bind.NewKeyedTransactorWithChainID(e.privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	startTime := time.Now().UnixNano()
	tx, err := instance.Transact(auth, invoke.Func, invoke.Args...)
	if err != nil {
		e.Logger.Error(err)
	}
	endTime := time.Now().UnixNano()
	if err != nil {
		return &fcom.Result{
			Label:     invoke.Func,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: startTime,
			SendTime:  endTime,
		}
	}
	ret := &fcom.Result{
		Label:     invoke.Func,
		UID:       tx.Hash(),
		Ret:       []interface{}{tx.Data()},
		Status:    fcom.Success,
		BuildTime: startTime,
		SendTime:  endTime,
	}

	return ret

}

// Confirm check the result of `Invoke` or `Transfer`
func (e *ETH) Confirm(result *fcom.Result, ops ...bcom.Option) *fcom.Result {
	if result.UID == "" ||
		result.UID == fcom.InvalidUID ||
		result.Status != fcom.Success ||
		result.Label == fcom.InvalidLabel {
		return result
	}
	tx, _, err := e.ethClient.TransactionByHash(context.Background(), result.UID.(common.Hash))
	result.ConfirmTime = time.Now().UnixNano()
	if err != nil || tx == nil {
		e.Logger.Error("invoke failed: %v", err)
		result.Status = fcom.Unknown
		return result
	}
	result.Status = fcom.Confirm
	return result
}

//Transfer transfer a amount of money from a account to the other one
func (e *ETH) Transfer(args bcom.Transfer, ops ...bcom.Option) (result *fcom.Result) {

	publicKey := e.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		e.Logger.Error("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		e.Logger.Error(err)
	}

	value := big.NewInt(args.Amount) // in wei (1 eth)
	gasLimit := uint64(21000)        // in units
	gasPrice, err := e.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		e.Logger.Error(err)
	}

	toAddress := common.HexToAddress(args.To)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := e.ethClient.NetworkID(context.Background())
	if err != nil {
		e.Logger.Error(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.privateKey)
	if err != nil {
		e.Logger.Error(err)
	}

	startTime := time.Now().UnixNano()
	err = e.ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		e.Logger.Error(err)
	}
	endTime := time.Now().UnixNano()

	if err != nil {
		e.Logger.Error(err)
		return &fcom.Result{
			Label:     fcom.BuiltinTransferLabel,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: startTime,
			SendTime:  endTime,
		}
	}

	ret := &fcom.Result{
		Label:     fcom.BuiltinTransferLabel,
		UID:       signedTx.Hash(),
		Ret:       []interface{}{tx.Data()},
		Status:    fcom.Success,
		BuildTime: startTime,
		SendTime:  endTime,
	}

	return ret
}

//SetContext set test group context in go client
func (e *ETH) SetContext(context string) error {
	e.Logger.Debugf("prepare msg: %v", context)
	msg := &Msg{}
	var (
		err error
	)

	if context == "" {
		e.Logger.Infof("Prepare nothing")
		return nil
	}

	err = json.Unmarshal([]byte(context), msg)
	if err != nil {
		e.Logger.Errorf("can not unmarshal msg: %v \n err: %v", context, err)
		return err
	}

	// set contract address
	e.contractAddress = msg.ContractAddress
	e.startBlock = msg.StartBlock
	return nil
}

//ResetContext reset test group context in go client
func (e *ETH) ResetContext() error {
	return nil
}

//GetContext generate TxContext
func (e *ETH) GetContext() (string, error) {

	msg := &Msg{
		ContractAddress: e.contractAddress,
		StartBlock:      e.startBlock,
	}

	bytes, error := json.Marshal(msg)
	if error != nil {
		fmt.Println(error)
	}

	return string(bytes), error
}

//Statistic statistic remote node performance
func (e *ETH) Statistic(statistic bcom.Statistic) (*fcom.RemoteStatistic, error) {

	from, to := statistic.From, statistic.To

	statisticData, err := GetTPS(e, from, to)

	if err != nil {
		return &fcom.RemoteStatistic{
			Start: from,
			End:   to,
		}, err
	}
	return statisticData, nil
}

func (e *ETH) Option(options bcom.Option) error {

	return nil
}

func KeystoreToPrivateKey(privateKeyFile, password string) (string, string, error) {
	keyjson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		fmt.Println("read keyjson file failed：", err)
	}
	unlockedKey, err := keystore.DecryptKey(keyjson, password)
	if err != nil {

		return "", "", err

	}
	privKey := hex.EncodeToString(unlockedKey.PrivateKey.D.Bytes())
	addr := crypto.PubkeyToAddress(unlockedKey.PrivateKey.PublicKey)
	return privKey, addr.String(), nil

}

// GetTPS calculates txnum and blocknum of pressure test
func GetTPS(e *ETH, beginTime, endTime int64) (*fcom.RemoteStatistic, error) {

	blockInfo, err := e.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var (
		blockCounter int
		txCounter    int
	)

	height := blockInfo.Number.Uint64()
	for i := e.startBlock; i < height; i++ {
		block, err := e.ethClient.BlockByNumber(context.Background(), new(big.Int).SetUint64(i))
		if err != nil {
			return nil, err
		}
		txCounter += len(block.Transactions())
		blockCounter++
	}

	statistic := &fcom.RemoteStatistic{
		Start:    beginTime,
		End:      endTime,
		BlockNum: blockCounter,
		TxNum:    txCounter,
	}
	return statistic, nil
}

// newContract initiates abi and bin files of contract
func newContract() (contract *Contract, err error) {
	files, err := ioutil.ReadDir(viper.GetString(fcom.ClientContractPath))
	var abiData, binData []byte
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if path.Ext(file.Name()) == ".abi" {
			abiData, err = ioutil.ReadFile(viper.GetString(fcom.ClientContractPath) + "/" + file.Name())
			if err != nil {
				return nil, err
			}
		}
		if path.Ext(file.Name()) == ".bin" {
			binData, err = ioutil.ReadFile(viper.GetString(fcom.ClientContractPath) + "/" + file.Name())
			if err != nil {
				return nil, err
			}
		}
	}

	abi := (string)(abiData)
	bin := (string)(binData)
	contract = &Contract{
		ABI: abi,
		BIN: bin,
	}
	return contract, nil
}
