package eth

import (
	"io/ioutil"
	"os"
	"testing"

	fcom "github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Skip()
	config := `
	[engine]
	rate = 1
	duration = "5s"
	cap = 1
	`
	ethConfig := `
	[rpc]
	node = "localhost"
	port = "a"
	`
	defer os.RemoveAll("./benchmark")

	os.MkdirAll("./benchmark/ethInvoke/eth", 0755)
	ioutil.WriteFile("./benchmark/ethInvoke/config.toml", []byte(config), 0644)
	ioutil.WriteFile("./benchmark/ethInvoke/eth/eth.toml", []byte(ethConfig), 0644)

	viper.AddConfigPath("benchmark/ethInvoke")
	viper.ReadInConfig()
	op := make(map[string]interface{})
	op["wkIdx"] = int64(0)
	op["vmIdx"] = int64(0)
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "eth",
		ConfigPath:   "./../../../benchmark/ethInvoke/eth",
		ContractPath: "./../../../benchmark/ethInvoke",
		Args:         nil,
		Options:      op,
	})
	client, err := New(b)
	assert.NotNil(t, client)
	assert.NoError(t, err)

	b.ConfigPath = ""
	client, err = New(b)
	assert.Nil(t, client)
	assert.Error(t, err)

	b.ConfigPath = "./benchmark/ethInvoke/eth"
	client, err = New(b)
	assert.Nil(t, client)
	assert.Error(t, err)

	viper.Set("rpc.port", "")
	client, err = New(b)
	assert.Nil(t, client)
	assert.Error(t, err)

	key1 := `
	{"address":"74d366e0649a91395bb122c005917644382b9452","crypto":{"cipher":"aes-128-ct","ciphertext":"fc4e8e2c753a98762828fad76697322da6a0143d6bfe223ce8a590637b433b75","cipherparams":{"iv":"9eab2eb01311d078ac7e3325150eecb2"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"0a8bda7b2e61a563a277601e65f6f30a92ab58e6e18f806105bb7218dff4c883"},"mac":"18f543410e2869a6a843166f1c3fb6aae5a5ec0dc6fdd41d1d76d8e8b19c5983"},"id":"98123f84-3855-4f12-b844-8c0d8ac02c09","version":3}
	`
	os.MkdirAll("./benchmark/ethInvoke/eth/keystore", 0755)
	ioutil.WriteFile("./benchmark/ethInvoke/eth/keystore/key1", []byte(key1), 0644)
	client, err = New(b)
	assert.Nil(t, client)
	assert.Error(t, err)

	key1 = `
	{"address":"74d366e0649a91395bb122c005917644382b9452","crypto":{"cipher":"aes-128-ctr","ciphertext":"fc4e8e2c753a98762828fad76697322da6a0143d6bfe223ce8a590637b433b75","cipherparams":{"iv":"9eab2eb01311d078ac7e3325150eecb2"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"0a8bda7b2e61a563a277601e65f6f30a92ab58e6e18f806105bb7218dff4c883"},"mac":"18f543410e2869a6a843166f1c3fb6aae5a5ec0dc6fdd41d1d76d8e8b19c5983"},"id":"98123f84-3855-4f12-b844-8c0d8ac02c09","version":3}
	`
	ioutil.WriteFile("./benchmark/ethInvoke/eth/keystore/key1", []byte(key1), 0644)
	client, err = New(b)
	assert.Nil(t, client)
	assert.Error(t, err)

}

func TestDeployContract(t *testing.T) {
	t.Skip()
	viper.Set("rpc.port", "8545")
	viper.Set("rpc.node", "localhost")
	op := make(map[string]interface{})
	op["wkIdx"] = int64(0)
	op["vmIdx"] = int64(0)
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "eth",
		ConfigPath:   "./../../../benchmark/ethInvoke/eth",
		ContractPath: "./../../../benchmark/ethInvoke",
		Args:         nil,
		Options:      op,
	})
	client, err := New(b)
	assert.NotNil(t, client)
	assert.NoError(t, err)

	err = client.DeployContract()
	assert.Error(t, err)

	b.ContractPath = "./../../../benchmark/ethInvoke/contract"
	client, _ = New(b)
	err = client.DeployContract()
	assert.NoError(t, err)

	defer os.RemoveAll("./benchmark")

	b.ContractPath = ""
	client, _ = New(b)
	err = client.DeployContract()
	assert.NoError(t, err)

	os.MkdirAll("./benchmark/ethInvoke", 0755)
	b.ContractPath = "./benchmark/ethInvoke/contract"
	client, _ = New(b)
	err = client.DeployContract()
	assert.Error(t, err)

	os.MkdirAll("./benchmark/ethInvoke/contract/a.abi", 0755)
	b.ContractPath = "./benchmark/ethInvoke/contract"
	client, _ = New(b)
	err = client.DeployContract()
	assert.Error(t, err)

	os.Remove("./benchmark/ethInvoke/contract/a.abi")
	os.MkdirAll("./benchmark/ethInvoke/contract/b.bin", 0755)
	err = client.DeployContract()
	assert.Error(t, err)

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	os.Remove("./benchmark/ethInvoke/contract/b.bin")
	abi := `
	[{"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"key","type":"string"},{"indexed":false,"internalType":"string","name":"value","type":"string"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"string","name":"","type":"string"}],"name":"items","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"key","type":"string"},{"internalType":"string","name":"value","type":"string"}],"name":"test","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]
	`
	ioutil.WriteFile("./benchmark/ethInvoke/contract/a.abi", []byte(abi), 0644)
	ioutil.WriteFile("./benchmark/ethInvoke/contract/b.bin", []byte("abi"), 0644)
	err = client.DeployContract()
	assert.Error(t, err)

}

func TestTransaction(t *testing.T) {
	t.Skip()
	//invoke
	viper.Set("rpc.port", "8545")
	viper.Set("rpc.node", "localhost")
	op := make(map[string]interface{})
	op["wkIdx"] = int64(0)
	op["vmIdx"] = int64(0)
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "eth",
		ConfigPath:   "./../../../benchmark/ethInvoke/eth",
		ContractPath: "./../../../benchmark/ethInvoke/contract",
		Args:         nil,
		Options:      op,
	})
	client, _ := New(b)
	client.DeployContract()
	res := client.Invoke(bcom.Invoke{Func: "test", Args: []interface{}{"foo", "bar"}})
	assert.NotNil(t, res)

	res = client.Invoke(bcom.Invoke{Func: "111", Args: []interface{}{"foo", "bar"}})
	assert.NotNil(t, res)
	//getcontext
	msg, err := client.GetContext()
	assert.NoError(t, err)

	//setcontext
	err = client.SetContext(msg)
	assert.NoError(t, err)

	client.contract.ABI = "111"
	msg, err = client.GetContext()
	assert.NoError(t, err)
	err = client.SetContext(msg)
	assert.Error(t, err)

	//transfer
	client, _ = New(b)
	res = client.Transfer(bcom.Transfer{From: "74d366e0649a91395bb122c005917644382b9452", To: "74d366e0649a91395bb122c005917644382b9452", Amount: int64(1)})
	assert.Equal(t, res.Status, fcom.Status("success"))
	// confirm
	res = client.Confirm(res)
	assert.Equal(t, res.Status, fcom.Status("confirm"))

	res.UID = ""
	res.Status = "success"
	res = client.Confirm(res)
	assert.Equal(t, res.Status, fcom.Status("success"))

	res.UID = "111"
	res = client.Confirm(res)
	assert.Equal(t, res.Status, fcom.Status("unknown"))

	client.nonce -= 1
	res = client.Transfer(bcom.Transfer{From: "74d366e0649a91395bb122c005917644382b9452", To: "74d366e0649a91395bb122c005917644382b9452", Amount: int64(1)})
	assert.Equal(t, res.Status, fcom.Status("failure"))

	defer os.RemoveAll("./benchmark")

	b = base.NewBlockchainBase(base.ClientConfig{
		ClientType: "eth",
		ConfigPath: "./benchmark/ethTransfer/eth",
		Args:       nil,
		Options:    op,
	})
	os.MkdirAll("./benchmark/ethTransfer/eth/keystore/-3", 0755)
	ioutil.WriteFile("./benchmark/ethTransfer/eth/keystore/-1", []byte(`{"address":"74d366e0649a91395bb122c005917644382b9452","crypto":{"cipher":"aes-128-ctr","ciphertext":"fc4e8e2c753a98762828fad76697322da6a0143d6bfe223ce8a590637b433b75","cipherparams":{"iv":"9eab2eb01311d078ac7e3325150eecb2"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"0a8bda7b2e61a563a277601e65f6f30a92ab58e6e18f806105bb7218dff4c883"},"mac":"18f543410e2869a6a843166f1c3fb6aae5a5ec0dc6fdd41d1d76d8e8b19c5983"},"id":"98123f84-3855-4f12-b844-8c0d8ac02c09","version":3}`), 0644)
	ioutil.WriteFile("./benchmark/ethTransfer/eth/keystore/-2", []byte(`{"address":"3b2b643246666bfa1332257c13d0d1283736838d","crypto":{"cipher":"aes-128-ctr","ciphertext":"50b10e30295ff3a5b729b3bc62e89145ebf6b5839cd3b8c13dcbbf099584cec6","cipherparams":{"iv":"fe3dd61296891e6654fd1b39ff2401a2"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1244ce8b522ff776be571c8814a7dfd6c8607aecd1f3b145385a27aa0d1443c7"},"mac":"78ff225aa55470d242b1464b9b42476313024bbfab343a5bbd27e748c67b44d8"},"id":"c066e226-9b72-4d63-a556-70034d0a135b","version":3}`), 0644)
	ioutil.WriteFile("./benchmark/ethTransfer/eth/eth.toml", []byte(``), 0644)
	client, _ = New(b)
	assert.Nil(t, client)

	os.Remove("./benchmark/ethTransfer/eth/keystore/-3")
	ioutil.WriteFile("./benchmark/ethTransfer/eth/keystore/-4", []byte(""), 0644)
	client, _ = New(b)
	assert.Nil(t, client)

	os.Remove("./benchmark/ethTransfer/eth/keystore/-4")
	client, _ = New(b)

	//getcontext
	msg, err = client.GetContext()
	assert.NoError(t, err)

	//setcontext
	err = client.SetContext(msg)
	assert.NoError(t, err)

	err = client.SetContext("")
	assert.NoError(t, err)

	err = client.SetContext("111")
	assert.Error(t, err)

	//statistic
	client.startBlock -= 1
	result, err := client.Statistic(bcom.Statistic{From: 0, To: 1})
	assert.NoError(t, err)
	assert.NotNil(t, result)

	result, err = client.Statistic(bcom.Statistic{From: 2, To: 1})
	assert.Error(t, err)
	assert.Nil(t, result)

	err = client.ResetContext()
	assert.NoError(t, err)

	//option
	err = client.Option(bcom.Option{})
	assert.NoError(t, err)

}
