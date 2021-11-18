package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
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

const gasLimit = 300000

//Contract contains the abi and bin files of contract
type Contract struct {
	ABI             string
	BIN             string
	parsedAbi       abi.ABI
	contractAddress common.Address
}

//ETH the client of eth
type ETH struct {
	*base.BlockchainBase
	ethClient  *ethclient.Client
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	auth       *bind.TransactOpts
	startBlock uint64
	contract   *Contract
	Accounts   map[string]*ecdsa.PrivateKey
	round      uint64
	nonce      uint64
	engineCap  uint64
	workerNum  uint64
	wkIdx      uint64
	vmIdx      uint64
}

//Msg contains message of context
type Msg struct {
	Contract *Contract
}

// New use given blockchainBase create ETH.
func New(blockchainBase *base.BlockchainBase) (client *ETH, err error) {
	log := fcom.GetLogger("eth")
	ethConfig, err := os.Open(viper.GetString(fcom.ClientConfigPath) + "/eth.toml")
	if err != nil {
		log.Errorf("load eth configuration fialed: %v", err)
		return nil, err
	}
	viper.MergeConfig(ethConfig)
	ethClient, err := ethclient.Dial("http://" + viper.GetString("rpc.node") + ":" + viper.GetString("rpc.port"))
	if err != nil {
		log.Errorf("ethClient initiate fialed: %v", err)
		return nil, err
	}
	files, err := ioutil.ReadDir(viper.GetString(fcom.ClientConfigPath) + "/keystore")
	if err != nil {
		log.Errorf("access keystore failed:%v", err)
		return nil, err
	}
	var (
		PublicK  *ecdsa.PublicKey
		PrivateK *ecdsa.PrivateKey
	)
	accounts := make(map[string]*ecdsa.PrivateKey)
	for i, file := range files {
		fileName := file.Name()
		account := fileName[strings.LastIndex(fileName, "-")+1:]
		privKey, _, err := KeystoreToPrivateKey(viper.GetString(fcom.ClientConfigPath)+"/keystore/"+fileName, "")
		if err != nil {
			log.Errorf("access account file failed: %v", err)
			return nil, err
		}

		privateKey, err := crypto.HexToECDSA(privKey)
		if err != nil {
			log.Errorf("privatekey encode failed %v ", err)
			return nil, err
		}
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
			return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		accounts[account] = privateKey
		if i == 0 {
			PublicK = publicKeyECDSA
			PrivateK = privateKey
		}
	}

	fromAddress := crypto.PubkeyToAddress(*PublicK)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Errorf("pending nonce failed: %v", err)
		return nil, err
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Errorf("generate gasprice failed: %v", err)
		return nil, err
	}
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		log.Errorf("get chainID failed: %v", err)
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(PrivateK, chainID)
	if err != nil {
		log.Errorf("generate transaction options failed: %v", err)
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(gasLimit) // in units
	auth.GasPrice = gasPrice
	startBlock, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Errorf("get number of headerblock failed: %v", err)
		return nil, err
	}
	workerNum := uint64(len(viper.GetStringSlice(fcom.EngineURLsPath)))
	if workerNum == 0 {
		workerNum = 1
	}
	vmIdx := uint64(blockchainBase.Options["vmIdx"].(int64))
	wkIdx := uint64(blockchainBase.Options["wkIdx"].(int64))
	client = &ETH{
		BlockchainBase: blockchainBase,
		ethClient:      ethClient,
		privateKey:     PrivateK,
		publicKey:      PublicK,
		auth:           auth,
		startBlock:     startBlock.Number.Uint64(),
		Accounts:       accounts,
		round:          0,
		nonce:          nonce,
		engineCap:      viper.GetUint64(fcom.EngineCapPath),
		workerNum:      workerNum,
		vmIdx:          vmIdx,
		wkIdx:          wkIdx,
	}
	return
}
func (e *ETH) DeployContract() error {
	contractPath := viper.GetString(fcom.ClientContractPath)
	if contractPath != "" {
		var er error
		e.contract, er = newContract()
		if er != nil {
			e.Logger.Errorf("initiate contract failed: %v", er)
			return er
		}
	} else {
		return nil
	}
	parsed, err := abi.JSON(strings.NewReader(e.contract.ABI))
	if err != nil {
		e.Logger.Errorf("decode abi of contract failed: %v", err)
		return err
	}
	e.contract.parsedAbi = parsed
	input := "1.0"
	contractAddress, tx, _, err := bind.DeployContract(e.auth, parsed, common.FromHex(e.contract.BIN), e.ethClient, input)
	if err != nil {
		e.Logger.Errorf("deploycontract failed: %v", err)
	}
	e.contract.contractAddress = contractAddress
	e.Logger.Info("contractAddress:" + contractAddress.Hex())
	e.Logger.Info("txHash:" + tx.Hash().Hex())
	return nil
}

//Invoke invoke contract with funcName and args in eth network
func (e *ETH) Invoke(invoke bcom.Invoke, ops ...bcom.Option) *fcom.Result {
	buildTime := time.Now().UnixNano()
	instance := bind.NewBoundContract(e.contract.contractAddress, e.contract.parsedAbi, e.ethClient, e.ethClient, e.ethClient)
	nonce := e.nonce + (e.wkIdx+e.round*e.workerNum)*(e.engineCap/e.workerNum) + e.vmIdx + 1
	e.round++
	gasPrice, err := e.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return &fcom.Result{
			Label:     invoke.Func,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}
	chainID, err := e.ethClient.NetworkID(context.Background())
	if err != nil {
		return &fcom.Result{
			Label:     invoke.Func,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}
	auth, err := bind.NewKeyedTransactorWithChainID(e.privateKey, chainID)
	if err != nil {
		return &fcom.Result{
			Label:     invoke.Func,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(gasLimit) // in units
	auth.GasPrice = gasPrice

	tx, err := instance.Transact(auth, invoke.Func, invoke.Args...)
	sendTime := time.Now().UnixNano()
	if err != nil {
		e.Logger.Errorf("invoke error: %v", err)
		return &fcom.Result{
			Label:     invoke.Func,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
			SendTime:  sendTime,
		}
	}
	ret := &fcom.Result{
		Label:     invoke.Func,
		UID:       tx.Hash().String(),
		Ret:       []interface{}{tx.Data()},
		Status:    fcom.Success,
		BuildTime: buildTime,
		SendTime:  sendTime,
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
	tx, _, err := e.ethClient.TransactionByHash(context.Background(), common.HexToHash(result.UID))
	result.ConfirmTime = time.Now().UnixNano()
	if err != nil || tx == nil {
		e.Logger.Errorf("query failed: %v", err)
		result.Status = fcom.Unknown
		return result
	}
	result.Status = fcom.Confirm
	return result
}

//Transfer transfer a amount of money from a account to the other one
func (e *ETH) Transfer(args bcom.Transfer, ops ...bcom.Option) (result *fcom.Result) {
	buildTime := time.Now().UnixNano()
	nonce := e.nonce + (e.wkIdx+e.round*e.workerNum)*(e.engineCap/e.workerNum) + e.vmIdx
	e.round++

	value := big.NewInt(args.Amount) // in wei (1 eth)
	gasLimit := uint64(gasLimit)     // in units
	gasPrice, err := e.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return &fcom.Result{
			Label:     fcom.BuiltinTransferLabel,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}

	toAddress := common.HexToAddress(args.To)
	data := []byte(args.Extra)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := e.ethClient.NetworkID(context.Background())
	if err != nil {
		return &fcom.Result{
			Label:     fcom.BuiltinTransferLabel,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.Accounts[args.From])
	if err != nil {
		return &fcom.Result{
			Label:     fcom.BuiltinTransferLabel,
			UID:       fcom.InvalidUID,
			Ret:       []interface{}{},
			Status:    fcom.Failure,
			BuildTime: buildTime,
		}
	}

	err = e.ethClient.SendTransaction(context.Background(), signedTx)
	sendTime := time.Now().UnixNano()
	if err != nil {
		e.Logger.Errorf("transfer error: %v", err)
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
		UID:       signedTx.Hash().String(),
		Ret:       []interface{}{tx.Data()},
		Status:    fcom.Success,
		BuildTime: buildTime,
		SendTime:  sendTime,
	}

	return ret
}

//SetContext set test group context in go client
func (e *ETH) SetContext(context string) error {
	e.Logger.Debugf("prepare msg: %v", context)
	msg := &Msg{}

	if context == "" {
		e.Logger.Infof("Prepare nothing")
		return nil
	}

	err := json.Unmarshal([]byte(context), msg)
	if err != nil {
		e.Logger.Errorf("can not unmarshal msg: %v \n err: %v", context, err)
		return err
	}

	// set contractaddress,abi,publickey
	e.contract = msg.Contract
	if e.contract != nil {
		parsed, err := abi.JSON(strings.NewReader(e.contract.ABI))
		if err != nil {
			e.Logger.Errorf("decode abi of contract failed: %v", err)
			return err
		}
		e.contract.parsedAbi = parsed
	}
	publicKey := e.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		e.Logger.Error("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	e.publicKey = publicKeyECDSA
	return nil
}

//ResetContext reset test group context in go client
func (e *ETH) ResetContext() error {
	return nil
}

//GetContext generate TxContext
func (e *ETH) GetContext() (string, error) {

	msg := &Msg{
		Contract: e.contract,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		e.Logger.Errorf("marshal msg failed: %v", err)
		return "", err
	}

	return string(bytes), err
}

//Statistic statistic remote node performance
func (e *ETH) Statistic(statistic bcom.Statistic) (*fcom.RemoteStatistic, error) {

	from, to := statistic.From, statistic.To

	statisticData, err := GetTPS(e, from, to)
	if err != nil {
		e.Logger.Errorf("getTPS failed: %v", err)
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
	log := fcom.GetLogger("eth")
	keyjson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		log.Errorf("read keyjson file failed：%v", err)
		return "", "", err
	}
	unlockedKey, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		log.Errorf("decryptKey failed: %v", err)
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
