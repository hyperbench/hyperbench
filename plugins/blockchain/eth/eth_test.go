package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestClient(t *testing.T) {
	t.Skip()
	client, err := ethclient.Dial("/Users/aiyoa/desktop/eth-test/data/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections
}
func TestBalance(t *testing.T) {
	t.Skip()
	client, err := ethclient.Dial("/Users/aiyoa/desktop/eth-test/data/geth.ipc")
	account := common.HexToAddress("01eec173917c429901b41b98ac3dd300e060e698")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance)
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue)
}
func TestTransfer(t *testing.T) {
	t.Skip()
	client, err := ethclient.Dial("/Users/aiyoa/desktop/eth-test/data/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}

	privKey, address, err := KeystoreToPrivateKey("/Users/aiyoa/desktop/eth-test/data/keystore/UTC--2021-10-21T03-41-34.288690000Z--01eec173917c429901b41b98ac3dd300e060e698", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("privKey:%s\naddress:%s\n", privKey, address)

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000) // in wei (1 eth)
	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("a6df6489927a9d0172185efe68de1f9aace82639")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	fmt.Println(signedTx.ChainId())
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

}

func TestDeployContract(t *testing.T) {
	t.Skip()
	client, err := ethclient.Dial("/Users/aiyoa/desktop/eth-test/data/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}

	privKey, address, err := KeystoreToPrivateKey("/Users/aiyoa/desktop/eth-test/data/keystore/UTC--2021-10-21T03-41-34.288690000Z--01eec173917c429901b41b98ac3dd300e060e698", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("privKey:%s\naddress:%s\n", privKey, address)
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	input := "1.0"
	contract, _ := newContract()
	parsed, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		log.Fatal(err)
	}
	contractAddress, tx, instance, err := bind.DeployContract(auth, parsed, common.FromHex(contract.BIN), client, input)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(contractAddress.Hex()) // 0xA103dA779fCB208c02759BB6fBD3eD1d354B9E16
	fmt.Println(tx.Hash().Hex())       // 0x1690003dca86eba1491e0aaa5a1cfde3fa39cafd90058537d8a0c8c4b6863d25

	_ = instance
}

func TestInvoke(t *testing.T) {
	t.Skip()
	client, err := ethclient.Dial("/Users/aiyoa/desktop/eth-test/data/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}

	privKey, address, err := KeystoreToPrivateKey("/Users/aiyoa/desktop/eth-test/data/keystore/UTC--2021-10-21T03-41-34.288690000Z--01eec173917c429901b41b98ac3dd300e060e698", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("privKey:%s\naddress:%s\n", privKey, address)
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	contractAddress := common.HexToAddress("0x7c376C8ED768018aa53d3C37Eed637912fEAA782")
	contract, _ := newContract()
	parsed, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		log.Fatal(err)
	}
	instance := bind.NewBoundContract(contractAddress, parsed, client, client, client)
	if err != nil {
		log.Fatal(err)
	}
	key := "foo"
	value := "bar"

	tx, err := instance.Transact(auth, "setItem", key, value)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", tx.Hash().Hex()) // tx sent: 0x5012ba7c07e46da3e1fbec454ed0e4079936b605d7ef0c0b0d0572972bb32dc6
}
