package fabric

import (
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	clientMSP "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	"log"
	"os"
)

// SDK struct
type SDK struct {
	*base.BlockchainBase
	ConfigFileName string
	SDK            *fabsdk.FabricSDK
	OrgName        string
	OrgAdmin       string
	MspIds         []string
	EndPoints      []string
	MSPClient      *clientMSP.Client
}

var fabSDK *SDK

// NewSDK init fabric SDK
func NewSDK(blockchainBase *base.BlockchainBase, configFile string) *SDK {
	if fabSDK != nil {
		return fabSDK
	}
	sdk := setupSDK(configFile)
	fabricsdk := &SDK{
		SDK:            sdk,
		OrgAdmin:       "Admin",
		BlockchainBase: blockchainBase,
	}
	fabricsdk.cleanupUserData()
	ctxProvider := sdk.Context()
	ctx, e := ctxProvider()
	if e != nil {
		blockchainBase.Logger.Errorf("failed to get context: %v\n", e)
		return fabricsdk
	}
	//fabricsdk.MspId = ctx.Identifier().MSPID
	client := ctx.IdentityConfig().Client()
	fabricsdk.OrgName = client.Organization
	peers := ctx.EndpointConfig().NetworkPeers()
	for _, peer := range peers {
		fabricsdk.EndPoints = append(fabricsdk.EndPoints, peer.URL)
		fabricsdk.MspIds = append(fabricsdk.MspIds, peer.MSPID)
	}
	//ctx.EndpointConfig().
	//caConfig, b := ctx.IdentityConfig().CAConfig(client.Organization)
	//if !b {
	//	logger.Errorf("failed to get ca config: %v\n", e)
	//	return fabricsdk
	//}
	//fabricsdk.OrgAdmin = caConfig.Registrar.EnrollID
	//fabricsdk.adminSecret = caConfig.Registrar.EnrollSecret

	fabSDK = fabricsdk
	return fabSDK
}

// GetChannelClient get channel.client from channelName, userName and OrgName
func (s *SDK) GetChannelClient(channelName string, userName string, orgName string) *channel.Client {
	clientChannelContext := s.SDK.ChannelContext(channelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Printf("failed to create new channel client: %s\n", err)
	}
	return client
}

// GetLedgerClient get client.ledger.client from channelName, userName and OrgName
func (s *SDK) GetLedgerClient(channelName string, userName string, orgName string) *ledger.Client {
	clientChannelContext := s.SDK.ChannelContext(channelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	client, err := ledger.New(clientChannelContext)
	if err != nil {
		log.Printf("failed to create new channel client: %s\n", err)
	}
	return client
}

// GetResmgmtClient get channel.client from channelName, userName and OrgName
func (s *SDK) GetResmgmtClient() *resmgmt.Client {
	//prepare context
	adminContext := s.SDK.Context(fabsdk.WithUser(s.OrgAdmin), fabsdk.WithOrg(s.OrgName))

	//org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Printf("failed to create new regmgmt client: %s\n", err)
	}
	return orgResMgmt
}

// GetTPS get tps
func (s *SDK) GetTPS() *resmgmt.Client {
	//prepare context
	adminContext := s.SDK.Context(fabsdk.WithUser(s.OrgAdmin), fabsdk.WithOrg(s.OrgName))

	//org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Printf("failed to create new regmgmt client: %s\n", err)
	}
	return orgResMgmt
}

//GetMspClient get clientMSP.Client
//if success return nil error
func (s *SDK) GetMspClient() (*clientMSP.Client, error) {
	if s.MSPClient != nil {
		return s.MSPClient, nil
	}
	//prepare context
	adminContext := s.SDK.Context()

	//org resource management client
	mspClient, err := clientMSP.New(adminContext)
	if err != nil {
		log.Printf("failed to create new regmgmt client: %s\n", err)
		return nil, err
	}
	registrarEnrollID, registrarEnrollSecret := getRegistrarEnrollmentCredentials(adminContext)
	if err := mspClient.Enroll(registrarEnrollID, clientMSP.WithSecret(registrarEnrollSecret)); err != nil {
		log.Fatalf("enroll registrar failed: %v", err)
		return nil, err
	}
	s.MSPClient = mspClient
	return mspClient, nil
}

func setupSDK(configFileName string) *fabsdk.FabricSDK {
	var config = config.FromFile(configFileName)
	sdk, err := fabsdk.New(config)
	if err != nil {
		log.Printf("failed to create new SDK: %s\n", err)
	}
	return sdk
}

func getRegistrarEnrollmentCredentials(ctxProvider context.ClientProvider) (string, string) {
	ctx, err := ctxProvider()
	if err != nil {
		fabSDK.Logger.Errorf("failed to get context: %v\n", err)
	}
	clientConfig := ctx.IdentityConfig().Client()
	caConfig, ok := ctx.IdentityConfig().CAConfig(clientConfig.Organization)
	if !ok {
		log.Printf("CAConfig failed: %v\n", err)
	}
	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
}

func (s *SDK) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

// CleanupUserData clean up CA dir
func (s *SDK) cleanupUserData() {
	configBackend, err := s.SDK.Config()
	if err != nil {
		fabSDK.Logger.Fatal(err)
	}
	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
	identityConfig, err := msp.ConfigFromBackend(configBackend)
	if err != nil {
		fabSDK.Logger.Fatal(err)
	}
	keyStorePath := cryptoSuiteConfig.KeyStorePath()
	credentialStorePath := identityConfig.CredentialStorePath()
	cleanupPath(keyStorePath)
	cleanupPath(credentialStorePath)
}

func cleanupPath(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		fabSDK.Logger.Fatalf("Cleaning up directory '%s' failed: %v", storePath, err)
	}
}
