//Package fabric provide operate for blockchain of fabric
package fabric

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"os"
)

var (
	goPath = os.Getenv("GOPATH")
)

// ExecuteCC invoke chaincode
func ExecuteCC(client *channel.Client, ccID, fcn string, args [][]byte, endpoints []string, invoke bool) (channel.Response, error) {
	ccConstruct := channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: args}
	if invoke {
		return client.Execute(ccConstruct, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(endpoints...))
	}

	return client.Query(ccConstruct, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(endpoints...))
}

// InstallCC install chaincode
func InstallCC(ccPath, ccID, ccVersion string, orgResMgmt *resmgmt.Client) ([]resmgmt.InstallCCResponse, error) {
	ccPkg, err := gopackager.NewCCPackage(ccPath, goPath)
	if err != nil {
		return nil, err
	}

	//install cc to org peers
	installCCRequest := resmgmt.InstallCCRequest{
		Name:    ccID,
		Path:    ccPath,
		Version: ccVersion,
		Package: ccPkg,
	}

	return orgResMgmt.InstallCC(installCCRequest, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
}

// InstantiateCC instantiate chaincode
func InstantiateCC(ccPath, ccID, ccVersion, channelID string, initArgs [][]byte, ccPolicy *common.SignaturePolicyEnvelope, orgResMgmt *resmgmt.Client) (resmgmt.InstantiateCCResponse, error) {

	instantiateCCRequest := resmgmt.InstantiateCCRequest{
		Name:    ccID,
		Path:    ccPath,
		Version: ccVersion,
		Args:    initArgs,
		Policy:  ccPolicy,
	}

	// Org resource manager will instantiate 'example_cc' on channel
	return orgResMgmt.InstantiateCC(channelID, instantiateCCRequest, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
}
