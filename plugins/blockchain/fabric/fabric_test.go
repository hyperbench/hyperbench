/*
@Time : 2019-04-12 14:13
@Author : lmm
*/
package fabric

//import (
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//var (
//	configs = `
//   [benchmark]
//      contract = "github.com/meshplus/hyperbench/benchmark/fabricExample/contract"    # 要测试的合约文件夹路径（作为参数传给deployContract)
//      script   = "benchmark/fabricExample/script.lua"  # 指定测试脚本路径
//      tps      = 1                                   # 指定tps(压力设置）
//      instant  = 1                                   # 按批发送 [optional][default:tps/10]
//      duration = "30s"                                 # 指定持续时间
//      skip     = false                                 # 跳过 [optional][default:false]
//      user     = 1                                     # 模拟用户数量(同时也是最大并发数) [optional][default:tps]
//      sign     = "ECDSA"                               # 签名类型 [optional]["ECDSA"|"SM2"][defualt:"ECDSA"]
//      confirm  = false                                 # 是否轮询结果 [optional][default:false]
//
//   [option]
//      channel     = "mychannel"
//      initArgs    = ["init","A","123","B","234"]
//      orgAdmin    = "Admin"
//      OrgName     = "Org1"
//      orgMspId    = "Org1MSP"
//      MSP         = false
//
//   # 网络配置
//   [network]
//      Name    = "fabric"
//      config  = "../../../config/fabric"                    # 仅仅指向配置文件夹的路径`
//)
//
//func TestNew(t *testing.T) {
//	client := &Fabric{}
//	e := client.New(configs)
//	assert.NoError(t, e)
//}
//
//func TestFabric(t *testing.T) {
//	var err error
//
//	//1.new
//	client := &Fabric{}
//	err = client.New(configs)
//	assert.NoError(t, err)
//
//	//print
//	//sdk := setupSDK("../../../config/fabric/config.yaml")
//	//adminContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
//	//resmgmtClient, _ := resmgmt.New(adminContext)
//	//fmt.Println(client.CCPath,"ccpath/n",client.SDK.GetResmgmtClient(),"ccsdk/n", resmgmtClient, "resclient/n")
//
//	//2.deploy
//	err = client.DeployContract()
//	assert.NoError(t, err)
//
//	//3.getContext
//	var context string
//	context, err = client.GetContext()
//	assert.NoError(t, err)
//
//	//4.setContext
//	err = client.SetContext(context)
//	assert.NoError(t, err)
//
//	////print
//	//intn := rand.Intn(len(client.AccountManager.Clients))
//	//fmt.Println("invoke account : ",intn)
//	//account, e := client.AccountManager.GetAccount(strconv.Itoa(intn))
//	//fmt.Println("account is:", account)
//	//fmt.Println("account name is: ", account.Name, "account orgname is: ", account.OrgName)
//	//assert.NoError(t, e)
//	//fmt.Println("channelid is: ", client.ChannelID, "admin", client.SDK.OrgAdmin, "orgname", client.SDK.OrgName)
//
//	//5.invokeContract
//	txResult := client.InvokeContract("query", "A")
//	fmt.Println(txResult)
//
//	//6.transfer
//	txResult = client.Transfer("0", "1", 1, "")
//	fmt.Println(txResult)
//
//	//7.reset
//	err = client.ResetContext()
//	assert.NoError(t, err)
//
//	//8.statistic
//	sta := client.Statistic(0, 1)
//	assert.NotNil(t, sta)
//
//	//string
//	s := client.String()
//	assert.NotNil(t, s)
//
//}
