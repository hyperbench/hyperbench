package hyperchain

import (
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//AccountType the type of sign account
type AccountType uint8

const (
	//ECDSA account type of ecdsa
	ECDSA AccountType = iota
	//SM2 account type of SM2
	SM2
)

//PASSWORD the default password of account
const PASSWORD = ""

//AccountManager the manager of account
type AccountManager struct {
	AccountType  AccountType
	Accounts     map[string]Account
	AccountsJSON map[string]string
	logger       *logging.Logger
}

//Account define the operate of account
type Account interface {
	GetAddress() common.Address
}

//NewAccountManager create a new AccountManager with keystore and accountType and return
func NewAccountManager(keystore, accountType string, logger *logging.Logger) *AccountManager {
	var (
		acType AccountType
	)

	switch strings.ToLower(accountType) {
	case "sm2":
		acType = SM2
	default:
		acType = ECDSA
	}

	am := &AccountManager{
		AccountType:  acType,
		Accounts:     map[string]Account{},
		AccountsJSON: map[string]string{},
		logger:       logger,
	}

	if keystore != "" {
		am.InitFromKeyStore(keystore, PASSWORD)
	}

	return am
}

//InitFromKeyStore init account with keystore and password
func (am *AccountManager) InitFromKeyStore(keystore, password string) {
	// get all file from keystore dir, try to parse it into account and store it in accounts map
	// notice that ioutil.ReadDir do not ensure file is sorted so that alias of account maybe different.
	var (
		acJSON    []byte
		acJSONStr string
		counter   int
		err       error
		rd        []os.FileInfo
	)
	if rd, err = ioutil.ReadDir(keystore); err == nil {
		for _, fi := range rd {
			if fi.IsDir() {
				continue
			}
			acJSON, _ = ioutil.ReadFile(filepath.Join(keystore, fi.Name()))
			acJSONStr = string(acJSON)
			_, _ = am.SetAccount(strconv.Itoa(counter), acJSONStr, password)
			counter++
		}

	}
}

//GetAccount get account with accountName and return
func (am *AccountManager) GetAccount(accountName string) (Account, error) {
	ac, ok := am.Accounts[accountName]
	if ok {
		return ac, nil
	}
	acJSON := am.genAccountJSON(PASSWORD)
	newAc, err := am.SetAccount(accountName, acJSON, PASSWORD)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot account %v from account json %v", accountName, acJSON))
	}
	return newAc, nil
}

//GetAccountJSON ger accountJson with accountName and return
func (am *AccountManager) GetAccountJSON(accountName string) (string, error) {
	acJSON, ok := am.AccountsJSON[accountName]
	if ok {
		return acJSON, nil
	}
	acJSON = am.genAccountJSON(PASSWORD)
	_, err := am.SetAccount(accountName, acJSON, PASSWORD)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("cannot account %v from account json %v", accountName, acJSON))
	}
	return acJSON, nil
}

//SetAccount set account with accountName, accountJson and password and return
func (am *AccountManager) SetAccount(accountName string, accountJSON string, password string) (Account, error) {
	var (
		ac  Account
		err error
	)
	switch am.AccountType {
	case ECDSA:
		ac, err = account.NewAccountFromAccountJSON(accountJSON, password)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("parse ecdsa account error"))

		}
	case SM2:
		ac, err = account.NewAccountSm2FromAccountJSON(accountJSON, password)
		if err != nil {
			return nil, errors.Wrap(err, "parse sm2 account error")
		}

	default:

		return nil, errors.New(fmt.Sprintf("unknow sign type %v", am.AccountType))

	}

	// Map account's name and address to account
	// then accountManager can get account through it's name or address
	am.Accounts[accountName] = ac
	am.Accounts[ac.GetAddress().Hex()] = ac

	// Map account's name to account but not the address
	// Account should only be used to generate and sync context of accounts
	am.AccountsJSON[accountName] = accountJSON
	return ac, nil
}

func (am *AccountManager) genAccountJSON(password string) string {
	// generate account json according to type require
	switch am.AccountType {
	case ECDSA:
		accountJSON, err := account.NewAccount(password)
		am.logger.Debugf("new ecdsa account %v", accountJSON)
		if err != nil {
			am.logger.Errorf("account gen: %v", err)
			panic(err)
			//return ""
		}
		return accountJSON
	case SM2:
		accountJSON, err := account.NewAccountSm2(password)
		am.logger.Debugf("new sm account %v", accountJSON)
		if err != nil {
			am.logger.Errorf("account gen: %v", err)
			panic(err)
			//return ""
		}
		return accountJSON
	default:
		panic(fmt.Sprintf("can not recognize sign type: %v", am.AccountType))
		//return ""
	}

}
