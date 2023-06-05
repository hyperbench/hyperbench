# hyperbench

hyperbench是一个由go编写的用于在区块链平台上进行压力测试的分布式压力测试工具。

详细介绍: [white paper](https://upload.hyperchain.cn/HyperBench%E7%99%BD%E7%9A%AE%E4%B9%A6.pdf)

## 我们的特色

- 灵活：基于虚拟机提供的Lua脚本和用户钩子来提供可编程的用例扩展

- 高效：虚拟机内置Go区块链客户端，具有统一的接口，无需额外的应用服务即可直接测试区块链系统

- 分布式：支持分布式测试功能，支持使用多个压力机同时测试区块链系统，简单易用

## 版本

| hyperbench 版本 | hyperbench 插件版本 |
|--------------------|----------------------------|
| v1.1.0             | v0.0.6                     |
| v1.0.5+            | v0.0.5                     |
| v1.0.4             | v0.0.4                     |
| v1.0.3             | v0.0.3                     |
| v1.0.2             | v0.0.2                     |
| v1.0.1             | v0.0.1                     |
| v1.0.0             | v0.0.1-alpha               |

## 注意
需要注意的是，当前系统使用的是go组件，go仅支持macOS和Linux系统，Windows系统尚未支持。

## 快速开始
###环境准备
平台支持多种操作系统如下
- Centos6 - Centos7
- Ubuntu14 – Ubuntu20
- Suse11，Suse12
- macOS
推荐操作系统
- centos7
- ubuntu18
如果您想在windows环境体验hyperbench平台，您可选择安装vmware虚拟机，然后安装合适的操作系统
VMware官网：https://www.vmware.com/cn/products/workstation-pro.html
推荐镜像（ubuntu18桌面版）:https://developer.aliyun.com/mirror/?spm=a2c6h.25603864.0.0.23583decq5r7CE  

### 安装go环境
```bash
#以go1.17版本为例
#从https://golang.google.cn/dl/下载go安装包
wget https://golang.google.cn/dl/go1.17.13.linux-amd64.tar.gz
tar -zxvf go1.17.13.linux-amd64.tar.gz
#将go可执行文件添加到环境变量
mv go/ ~/
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
echo 'export GOPATH=~/go' >> ~/.bashrc

source ~/.bashrc

#添加国内代理
go env -w GOPROXY=https://goproxy.cn,direct
 
```

### 安装docker及docker-compose
```bash
#docker版本建议v 19.03及以上，docker-compose建议v1.27.4及以上

sudo yum install -y yum-utils
sudo yum-config-manager  --add-repo  https://download.docker.com/linux/centos/docker-ce.repo
yum install docker-ce docker-ce-cli containerd.io
# 设置为开机启动 
systemctl enable docker 
systemctl daemon-reload
# 启动 docker 
systemctl start docker

#安装docker-compose,安装命令如下： 
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
 
```

### 安装 packr
```bash
# 基于  Go 1.16 及以上版本
go install github.com/gobuffalo/packr/v2/packr2@latest
# 基于  Go 1.15 已以下版本
go install github.com/gobuffalo/packr/packr2@latest
```
更多安装方式请查看  https://github.com/gobuffalo/packr

### 从源代码构建hyperbench

```bash
# 克隆hyperbench源代码至 hyperbench 工作目录:
mkdir -p $GOPATH/src/github.com/hyperbench && cd $GOPATH/src/github.com/hyperbench
git clone https://github.com/hyperbench/hyperbench.git
cd hyperbench

# 编译
make build

# 将编译完成的hyperbench可执行文件拷贝至环境变量
cp hyperbench $GOPATH/bin

```

### 从源代码构建联盟链插件

```bash
# 将hyperbench-plugins源代码克隆至hyperbench-plugins 工作目录:
cd $GOPATH/src/github.com/hyperbench
git clone https://github.com/hyperbench/hyperbench-plugins.git

# 以编译hyperchain插件为例
cd hyperbench-plugins/hyperchain
make build
cp hyperchain.so ../../hyperbench/hyperchain.so
```

### hyperbench  对接  hyperchain

在启动压力测试前，您需要提前安装 docker和docker-compose 来完成hyperchain网络的搭建.

1. 启动hyperchain
```bash
# start hyperchain network
cd $GOPATH/src/github.com/hyperbench/hyperbench/benchmark/hyperchain
bash control.sh start
# 如果这是你第一次拉取docker镜像，请耐心等待
# 如果出现超时错误，请重试
```
  想了解更多hyperchain相关信息，请访问 https://hyperchain.cn/
2. 开始压力测试.
```bash
# 使用 benchmark/hyperchain/local 目录下的测试用例，开始压力测试
cd $GOPATH/src/github.com/hyperbench/hyperbench
hyperbench start benchmark/hyperchain/local
```

### hyperbench  对接  Fabric

在启动压力测试前，您需要提前安装 docker和docker-compose 来完成fabric网络的搭建.
1. 启动一个fabric网络

```bash
# 启动一个fabric网络
cd $GOPATH/src/github.com/hyperbench/hyperbench/benchmark/fabric/example/fabric
bash deamon.sh
# 如果这是你第一次拉取docker镜像，请耐心等待
# 如果出现超时错误，请重试
```
  想了解更多fabric相关信息，请访问  https://github.com/hyperledger/fabric
2. 开始压力测试

```bash
# 使用 benchmark/fabric/example 目录下的测试用例 case, 开始压力测试
cd $GOPATH/src/github.com/hyperbench/hyperbench
hyperbench start benchmark/fabric/example
```

### hyperbench  对接  ethereum

1. 创建Geth可执行文件
```bash
# 将go-ethereum源代码克隆至本地并编译
git clone https://github.com/ethereum/go-ethereum.git
cd go-ethereum
make geth
```
  以上命令将在go-ethereum/build/bin文件夹中创建一个Geth可执行文件，如果需要，可将该文件放入环境变量。二进制文件是独立的，不需要任何额外的文件。
2. 启动 ethereum 网络
```bash
#启动之前，您需要一个创世账户 来初始化geth节点，以确保正确设置所有区块链参数，完成初始化后，启动
#初始化本地eth
geth --datadir ./test/ init path/genesis.json

#启动本地eth，更多参数说明请参考 ./geth --help
nohup ./geth --datadir "./test/" --http --http.addr=0.0.0.0 --http.port 8545 --http.corsdomain "*" --http.api "eth,net,web3,personal,admin,txpool,debug,miner" --nodiscover --maxpeers 30 --networkid 2081 --port 30403 --allow-insecure-unlock  2>> ./test/geth.log &
```

3. 准备进行转账压力测试
在压力测试之前，您需要准备两个账户，且其中一个账户拥有balance (也可使用拥有余额的创世账户)
```bash
#将eth数据目录下的账户私钥文件拷贝至benchmark工作目录下
cd $GOPATH/src/github.com/hyperbench/hyperbench
rm  benchmark/ethereum/transfer/eth/keystore/*
cp PATH/keystore/UTC-* benchmark/ethereum/transfer/eth/keystore/

#配置转账交易的发起方和接收方
cd $GOPATH/src/github.com/hyperbench/hyperbench
vim benchmark/ethereum/transfer/script.lua

#配置用于http访问的eth端口
cat benchmark/ethereum/transfer/eth/eth.toml
```
  想了解更多ethereum相关信息，请访问   https://github.com/ethereum/go-ethereum

4. 开启压力测试
```bash
# 使用 benchmark/ethereum/transfer/example 中的测试用例, 开始压力测试
cd $GOPATH/src/github.com/hyperbench/hyperbench
hyperbench start benchmark/ethereum/transfer/
```



### Possible problems
1. 您可能会遇到类似问题，这是由于软件包版本的自动更新，主存储库和插件使用了不同的软件包版本。唯一的解决办法是，在编译之前，将讲两者的go-plugin版本对齐
```bash
[blockchain][ERROR] 10:37:45.853 blockchain.go:39 plugin failed: plugin.Open("./hyperchain"): plugin was built with a different version of package golang.org/x/sys/internal/unsafeheader
```
    针对以上问题，您需要在$GOPATH/src/github.com/hyperbench/hyperbench/go.mod中找到提示的插件版本，同步至$GOPATH/src/github.com/hyperbench/hyperbench-plugins/hyperchain/go.mod中，并重新编译hyperbench-plugins

2. 当您在测试eth时，若最终得到的统计数据为0，并不是代表交易失败，而是因为eth必须要有矿工工作才能落块，以下是其中一种解决办法
```bash
#进入ipc控制台界面
./geth attach test/geth.ipc

#账户绑定矿工（挖矿时，该账户会有余额增长）
miner.setEtherbase("0x7a3f4c28a85ec58f3c8bd862ab201b8a03dfff5f")

#矿工开始工作
miner.start()

#矿工停止工作
miner.stop()
```

3.缺少依赖
```bash
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH

#安装依赖
sudo yum install -y gcc 或 sufo apt install gcc
 ```

### 其他
hyperbench是一个区块链测试工具，基准测试中提供的相应示例仅供参考。在测试相应的区块链平台时，使用相应的区块链平台是指相应平台的用户手册

## 捐赠

感谢您考虑帮助编写源代码！无论是系统错误报告、新功能还是代码贡献，我们都非常欢迎。

有关详细信息，请查看投稿指南。

## 许可

hyperbench目前使用Apache 2.0许可证。有关详细信息，请参阅许可证文件。
