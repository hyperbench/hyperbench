package fabric

import (
	clientMSP "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/meshplus/hyperbench/common"
	"github.com/stretchr/testify/assert"
	"testing"
	//fabsdk "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func TestClientManager(t *testing.T) {
	s := &SDK{
		ConfigFileName: "config.yaml",
		SDK:            &fabsdk.FabricSDK{},
		OrgName:        "org1",
		OrgAdmin:       "Admin",
		MspIds:         []string{"Org1MSP"},
		EndPoints:      nil,
		MSPClient:      &clientMSP.Client{},
	}
	cm, err := NewClientManager(s, false, common.GetLogger("client"))
	assert.NotNil(t, cm)
	assert.Nil(t, err)

	//cli, e2 := cm.GetAccount("test")
	//assert.NotNil(t, cli)
	//assert.Nil(t, e2)

	//for i:= 0; i<2; i++ {
	//	cli, e := client.GetAccount(strconv.Itoa(i))
	//	assert.NotNil(t, cli)
	//	assert.Nil(t, e)
	//}

}
