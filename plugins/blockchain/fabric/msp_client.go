package fabric

import (
	"errors"
	"strconv"
	"time"

	clientMSP "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/op/go-logging"
)

//SECRET the default secret for account
const SECRET = "123456"

//ClientManager the manager of client
type ClientManager struct {
	Clients   map[string]*Client
	MspClient *clientMSP.Client
	OrgName   string
	EnableMSP bool
	Logger    *logging.Logger
}

//Client contains a client account fields
type Client struct {
	Key       string
	Name      string
	Secret    string
	EnrSecret string
	OrgName   string
	IsEnroll  bool
}

//NewClientManager create a ClientManager with sdk
//return nil error if success
func NewClientManager(sdk *SDK, enableMSP bool, logger *logging.Logger) (*ClientManager, error) {
	var (
		e         error
		mspClient *clientMSP.Client
	)

	// generate the msp client if msp is enabled
	if enableMSP {
		mspClient, e = sdk.GetMspClient()
		if e != nil {
			return nil, e
		}
	} else {
		mspClient = nil
	}

	cm := &ClientManager{
		Clients:   map[string]*Client{},
		OrgName:   sdk.OrgName,
		MspClient: mspClient,
		EnableMSP: enableMSP,
		Logger:    logger,
	}
	return cm, nil
}

//GetAccount get account with name
//if success, return nil error
func (cm *ClientManager) GetAccount(name string) (*Client, error) {
	client := cm.Clients[name]
	if client != nil {
		return client, nil
	}
	client = &Client{
		Key:     name,
		Name:    name + strconv.Itoa(int(time.Now().UnixNano())),
		Secret:  SECRET,
		OrgName: cm.OrgName,
	}
	e := cm.register(client)
	if e != nil {
		return nil, e
	}
	e = cm.enroll(client)
	if e != nil {
		return nil, e
	}
	return client, nil
}

//InitAccount init the number of account
//if success, return nil error
func (cm *ClientManager) InitAccount(count int) error {
	var e error
	if cm.EnableMSP {
		// use random Account
		for i := 0; i < count; i++ {
			_, e = cm.GetAccount(strconv.Itoa(i))
			if e != nil {
				cm.Logger.Error(e)
			}
		}
	} else {
		// use only Admin
		admin := &Client{
			Key:       "0",
			Name:      "Admin",
			Secret:    "",
			EnrSecret: "",
			OrgName:   cm.OrgName,
			IsEnroll:  true,
		}
		for i := 0; i < count; i++ {
			cm.Clients[strconv.Itoa(i)] = admin
		}
	}

	return e
}

func (cm *ClientManager) register(c *Client) error {
	registReq := &clientMSP.RegistrationRequest{
		Name:        c.Name,
		Secret:      c.Secret,
		Affiliation: c.OrgName,
	}
	secret, e := cm.MspClient.Register(registReq)
	if e != nil {
		return e
	}
	c.EnrSecret = secret
	cm.Clients[c.Key] = c
	return nil
}

func (cm *ClientManager) enroll(c *Client) error {
	c = cm.Clients[c.Key]
	if c == nil {
		return errors.New("user not register")
	}
	if c.IsEnroll {
		return nil
	}

	err := cm.MspClient.Enroll(c.Name, clientMSP.WithSecret(c.EnrSecret))
	if err != nil {
		return err
	}
	c.IsEnroll = true
	return nil
}
