package hyperchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/meshplus/gosdk/rpc"
	fcom "github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	t.Skip()
	op := make(map[string]interface{})
	op["keystore"] = "./../../../benchmark/evmType/keystore"
	op["sign"] = "sm2"
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "hyperchain",
		ConfigPath:   "./../../../benchmark/evmType/hyperchain",
		ContractPath: "./../../../benchmark/evmType/hyperchain",
		Args:         nil,
		Options:      op,
	})
	hpc, _ := NewClient(b)
	acc, err := hpc.am.GetAccountJSON("")
	assert.NotNil(t, acc)
	assert.NoError(t, err)
	a, _ := hpc.am.GetAccount("1")
	hpc.sign(&rpc.Transaction{}, a)

	ac, err := hpc.am.GetAccount("111")
	assert.NotNil(t, ac)
	assert.NoError(t, err)
	ac, err = hpc.am.GetAccount("111")
	assert.NotNil(t, ac)
	assert.NoError(t, err)

	acc, err = hpc.am.GetAccountJSON("111")
	assert.NotNil(t, acc)
	assert.NoError(t, err)

	Acc, err := hpc.am.SetAccount("", "", "")
	assert.Nil(t, Acc)
	assert.Error(t, err)

	b = base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "hyperchain",
		ConfigPath:   "./../../../benchmark/evmType/hyperchain",
		ContractPath: "./../../../benchmark/evmType/hyperchain",
		Args:         nil,
		Options:      nil,
	})
	hpc, _ = NewClient(b)
	acc, err = hpc.am.GetAccountJSON("111")
	assert.NotNil(t, acc)
	assert.NoError(t, err)

	Acc, err = hpc.am.SetAccount("", "", "")
	assert.Nil(t, Acc)
	assert.Error(t, err)

	op["keystore"] = "/"
	op["sign"] = "sm2"
	b = base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "hyperchain",
		ConfigPath:   "./../../../benchmark/evmType/hyperchain",
		ContractPath: "./../../../benchmark/evmType/hyperchain",
		Args:         nil,
		Options:      op,
	})
	hpc, _ = NewClient(b)

	hpc.am.AccountType = 2

	Acc, err = hpc.am.SetAccount("", "", "")
	assert.Nil(t, Acc)
	assert.Error(t, err)

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	acc, err = hpc.am.GetAccountJSON("11")
	assert.NotNil(t, acc)
	assert.NoError(t, err)
}

func TestHyperchain(t *testing.T) {
	t.Skip()

	//evm
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "hyperchain",
		ConfigPath:   "./../../../benchmark/evmType/hyperchain",
		ContractPath: "./../../../benchmark/evmType",
		Args:         nil,
		Options:      nil,
	})
	hpc, err := NewClient(b)
	assert.NotNil(t, hpc)
	assert.NoError(t, err)
	err = hpc.DeployContract()
	assert.NoError(t, err)

	args, a := make(map[interface{}]interface{}), make(map[interface{}]interface{})
	args[1] = 1
	res := hpc.Invoke(bcom.Invoke{
		Func: "typeUint8",
		Args: []interface{}{args},
	})
	fmt.Println(res)
	assert.Equal(t, res.Status, fcom.Status(""))

	b.ContractPath = "./../../../benchmark/evmType/contract"
	hpc, _ = NewClient(b)
	hpc.DeployContract()

	// evm invoke, confirm
	args[float64(1)] = "111"
	args[float64(2)] = a
	res = hpc.Invoke(bcom.Invoke{
		Func: "typeUint8",
		Args: []interface{}{args},
	})
	assert.Equal(t, res.Status, fcom.Failure)

	res = hpc.Invoke(bcom.Invoke{
		Func: "typeString",
		Args: []interface{}{"test string"},
	})
	assert.Equal(t, res.Status, fcom.Success)

	res = hpc.Invoke(bcom.Invoke{
		Func: "typeString",
		Args: []interface{}{"test string"},
	})
	assert.Equal(t, res.Status, fcom.Success)
	res.Label = "ttt"
	hpc.Confirm(res)
	assert.Equal(t, res.Status, fcom.Confirm)

	hpc.op.poll = true
	res = hpc.Invoke(bcom.Invoke{
		Func: "typeString",
		Args: []interface{}{"test string"},
	})
	assert.Equal(t, res.Status, fcom.Confirm)

	res = hpc.Invoke(bcom.Invoke{
		Func: "typeBool",
		Args: []interface{}{"true", []string{"false"}, []string{"false", "true", "false"}},
	})
	assert.Equal(t, res.Status, fcom.Confirm)

	bytes, err := json.Marshal(Msg{Contract: hpc.contract.ContractRaw})
	assert.NotNil(t, bytes)
	assert.NoError(t, err)

	err = hpc.SetContext(string(bytes))
	assert.NoError(t, err)

	//evm DeployContract
	defer os.RemoveAll("./benchmark")

	os.Mkdir("./benchmark", 0755)

	b.ContractPath = "./benchmark/evm1/contract"

	os.MkdirAll("./benchmark/evm1/contract/evm", 0755)
	ioutil.WriteFile("./benchmark/evm1/contract/evm/test.addr", []byte(""), 0644)
	ioutil.WriteFile("./benchmark/evm1/contract/evm/test.abi", []byte(""), 0644)

	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.Error(t, err)

	ioutil.WriteFile("./benchmark/evm1/contract/evm/test.addr", []byte("0xc6a91501d2ff05467f2336898da266d6de60c4"), 0644)
	err = hpc.DeployContract()
	assert.NoError(t, err)

	bytes, err = json.Marshal(Msg{Contract: hpc.contract.ContractRaw})
	assert.NotNil(t, bytes)
	assert.NoError(t, err)

	hpc.SetContext(string(bytes))
	assert.NoError(t, err)

	os.MkdirAll("./benchmark/evm2/contract/evm", 0755)
	ioutil.WriteFile("./benchmark/evm2/contract/evm/test.solc", []byte(""), 0644)
	b.ContractPath = "./evm/contract"

	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.NoError(t, err)

	b.ContractPath = "./benchmark/evm2/contract"
	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.Error(t, err)

	//jvm DeployContract
	os.MkdirAll("./benchmark/jvm1/contract/jvm", 0755)
	ioutil.WriteFile("./benchmark/jvm1/contract/jvm/test.addr", []byte(""), 0644)

	b.ContractPath = "./benchmark/jvm1/contract"

	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.Error(t, err)

	ioutil.WriteFile("./benchmark/jvm1/contract/jvm/test.addr", []byte("0xc6a91501d2ff05467f2336898da266d6de60c41111"), 0644)
	err = hpc.DeployContract()
	assert.NoError(t, err)

	//jvm invoke
	hpc.op.nonce = 0
	hpc.op.fakeSign = true
	res = hpc.Invoke(bcom.Invoke{
		Func: "typeUint8",
		Args: []interface{}{"1", "2"},
	})
	assert.Equal(t, res.Status, fcom.Failure)

	os.MkdirAll("./benchmark/jvm2/contract/jvm", 0755)
	ioutil.WriteFile("./benchmark/jvm2/contract/jvm/test.java", []byte(""), 0644)

	b.ContractPath = "./benchmark/jvm2/contract"
	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.Error(t, err)

	b.ContractPath = "./../../../benchmark/javaContract/contract"
	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.Error(t, err)

	//hvm DeployContract
	b.ContractPath = "./../../../benchmark/hvmSBank/contract"
	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.NoError(t, err)

	//hvm invoke, confirm
	res = hpc.Invoke(bcom.Invoke{
		Func: "typeUint8",
		Args: []interface{}{args},
	})
	assert.Equal(t, res.Status, fcom.Failure)

	res = hpc.Invoke(bcom.Invoke{
		Func: "com.hpc.sbank.invoke.IssueInvoke",
		Args: []interface{}{args},
	})
	assert.Equal(t, res.Status, fcom.Failure)

	hpc.op.poll = true
	res = hpc.Invoke(bcom.Invoke{
		Func: "com.hpc.sbank.invoke.IssueInvoke",
		Args: []interface{}{"1", "1000000"},
	})
	assert.Equal(t, res.Status, fcom.Confirm)

	bytes, err = json.Marshal(Msg{Contract: hpc.contract.ContractRaw, Accounts: map[string]string{"11": "11"}})
	assert.NotNil(t, bytes)
	assert.NoError(t, err)

	err = hpc.SetContext(string(bytes))
	assert.NoError(t, err)

	os.MkdirAll("./benchmark/hvm2/contract/hvm", 0755)
	ioutil.WriteFile("./benchmark/hvm2/contract/hvm/test.jar", []byte(""), 0644)
	ioutil.WriteFile("./benchmark/hvm2/contract/hvm/test.abi", []byte(""), 0644)
	b.ContractPath = "./benchmark/hvm2/contract"
	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.Error(t, err)

	os.MkdirAll("./benchmark/hvm1/contract/hvm", 0755)
	ioutil.WriteFile("./benchmark/hvm1/contract/hvm/test.addr", []byte(""), 0644)
	ioutil.WriteFile("./benchmark/hvm1/contract/hvm/test.abi", []byte(""), 0644)
	b.ContractPath = "./benchmark/hvm1/contract"
	hpc, _ = NewClient(b)
	err = hpc.DeployContract()
	assert.Error(t, err)

	ioutil.WriteFile("./benchmark/hvm1/contract/hvm/test.addr", []byte("0xc6a91501d2ff05467f2336898da266d6de60c4"), 0644)
	err = hpc.DeployContract()
	assert.NoError(t, err)

	//confirm
	res = hpc.Confirm(&fcom.Result{UID: ""})
	assert.Equal(t, res.Status, fcom.Status(""))

	res = hpc.Confirm(&fcom.Result{UID: "111", Status: fcom.Success, Label: "111"})
	assert.Equal(t, res.Status, fcom.Unknown)

	//transfer
	res = hpc.Transfer(bcom.Transfer{From: "0", To: "1", Amount: 0, Extra: ""})
	assert.Equal(t, res.Status, fcom.Success)

	hpc.op.extraIDStr = []string{"1"}
	hpc.op.extraIDInt64 = []int64{1}
	hpc.op.nonce = int64(1)
	hpc.op.poll = true
	res = hpc.Transfer(bcom.Transfer{From: "0", To: "1", Amount: 0, Extra: ""})
	assert.Equal(t, res.Status, fcom.Confirm)

	//getcontext,setcontext
	contract := hpc.contract.ContractRaw

	msg, err := hpc.GetContext()
	assert.NotNil(t, msg)
	assert.NoError(t, err)
	hpc.am = nil
	msg, err = hpc.GetContext()
	assert.NotNil(t, msg)
	assert.NoError(t, err)

	err = hpc.ResetContext()
	assert.NoError(t, err)

	err = hpc.SetContext("")
	assert.NoError(t, err)

	err = hpc.SetContext("111")
	assert.Error(t, err)

	Msg := Msg{
		Contract: contract,
	}
	bytes, err = json.Marshal(Msg)
	assert.NotNil(t, bytes)
	assert.NoError(t, err)

	err = hpc.SetContext(string(bytes))
	assert.Error(t, err)

	result, err := hpc.Statistic(bcom.Statistic{From: 1, To: 1})
	assert.NotNil(t, result)
	assert.NoError(t, err)

	result, err = hpc.Statistic(bcom.Statistic{From: -1, To: 1})
	assert.Nil(t, result)
	assert.Error(t, err)

	m := make(map[string]interface{})
	m["confirm"] = true
	m["simulate"] = true
	m["account"] = "true"
	m["nonce"] = float64(1)
	m["extraid"] = []interface{}{"11", float64(1)}
	err = hpc.Option(m)
	assert.NoError(t, err)

	m["account"] = true
	err = hpc.Option(m)
	assert.Error(t, err)

	m["account"] = "true"
	m["simulate"] = "true"
	err = hpc.Option(m)
	assert.Error(t, err)

	m["account"] = "true"
	m["simulate"] = true
	m["confirm"] = "true"
	err = hpc.Option(m)
	assert.Error(t, err)

}
