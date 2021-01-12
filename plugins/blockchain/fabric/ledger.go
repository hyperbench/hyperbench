package fabric

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/meshplus/hyperbench/common"
)

// GetTPS get remote tps
func GetTPS(client *ledger.Client, startNum uint64, beginTime, endTime int64) (*common.RemoteStatistic, error) {

	blockInfo, err := client.QueryInfo()
	if err != nil {
		return nil, err
	}

	var (
		blockCounter int
		txCounter    int
	)

	height := blockInfo.BCI.Height
	for i := startNum; i < height; i++ {
		block, err := client.QueryBlock(i)
		if err != nil {
			return nil, err
		}
		txCounter += len(block.GetData().Data)
		blockCounter++
	}

	statistic := &common.RemoteStatistic{
		Start:    beginTime,
		End:      endTime,
		BlockNum: blockCounter,
		TxNum:    txCounter,
	}
	return statistic, nil
}
