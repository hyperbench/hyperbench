本文是hyperbench的入门手册，包含了最基本的安装、配置和使用。

# 第一章 环境安装

## 编译安装-Go环境与Hyperbench

Hyperbench基于Go语言开发，因此需要在系统中安装golang开发环境。安装方式主要有3种：

### 源码安装Golang

从 [官网](https://golang.org/dl/)下载源码，找到相应的`goVersion.src.tar.gz`文件下载到本地`$HOME`目录，解压至`$HOME`目录，并运行`all.bash`脚本安装，出现`ALL TESTS PASSED`字样时安装成功，然后设置环境变量。假设下载的目标版本是go1.12.4，那么执行的脚本如下：

```javascript
# 下载并解压
cd $HOME
wget https://dl.google.com/go/go1.15.6.src.tar.gz
tar -xvf go1.15.6.src.tar.gz

# 编译安装
cd go/src
./all.bash

# 设置环境变量
echo 'export PATH=$PATH:$HOME/go/bin:' >> ~/.bashrc
source ~/.bashrc
```

### 标准包安装Golang

Go语言[官网](https://golang.org/dl/)提供了大部分平台打包好的一键安装包，将这些包默认安装到`/usr/local/go`

**Mac**

下载系统对应的安装包(pkg文件)进行安装，完成之后go会自动被添加进PATH中，此时在终端输入go即可使用。

**Linux**

下载安装包并解压到安装目录，假设是`$HOME`目标安装目录，安装的目标版本为go1.15.6(64位)，那么安装执行以下脚本即可:

```shell
# 下载并解压
cd $HOME
wget https://dl.google.com/go/go1.15.6.linux-amd64.tar.gz
tar -xvf go1.15.6.linux-amd64.tar.gz
# 添加到path
echo 'export PATH=$PATH:$HOME/go/bin:' >> ~/.bashrc
source ~/.bashrc
```

### **第三方工具安装**Golang（推荐）

**yum**

```text
sudo yum install go
```

**apt-get**

```text
sudo apt-get install python-software-properties
sudo add-apt-repository ppa:gophers/go
sudo apt-get update
sudo apt-get install golang-stable
```

**homebrew**

homebrew是Mac系统下使用最多的管理工具软件，推荐在Mac上使用homebrew安装和管理

```text
# 安装homebrew
/usr/bin/ruby -e "$(curl -fsSLhttps://raw.githubusercontent.com/Homebrew/install/master/install)"
# homebrew更新和go安装
brew update && brew upgrade
brew install go
```

### 编译安装最新的`Hyperbench`

安装Hyperbench也有多种方式，任选其一即可。

**注意编译需要要求go版本1.11+**

下载源码

```text
cd $GOPATH/src
git clone git@github.com:hyperbench/hyperbench.git
# 选择版本分支，例如v1.0.1
git checkout v1.0.1

export GO111MODULE=on 
```

使用makefile进行编译

```text
make build
# 此处出现packr缺失错误时，请确保packr的安装路径在您的环境目录下
```

编译后会生成二进制文件`hyperbench`，可以按照说明进行运行，推荐将`hyperbench`放到`/usr/local/bin`或者`$GOPATH/bin`目录下，方便使用。
### 编译区块链适配插件`Hyperbench-plugins`
下载源码
``` text
git clone https://github.com/hyperbench/hyperbench-plugins.git
# 选择对应版本分支,例如v0.0.1
git checkout v0.0.1
# 选择需要使用的区块链适配插件，例如hyperchain
cd hyperbench-plugins/hyperchain
# 编译插件,插件名称在配置阶段需要使用
make build
```

# 第二章 使用与配置

第二章中的命令默认您已经将hyperbench安装到了`/usr/local/bin`或者`$GOPATH/bin`目录下了，如果您没有安装hyperbench，也可以将其当做一个普通的二进制文件使用`./hyperbench [cmd]`进行各种操作，此时必须确保在目录文件夹下含有hyperbench二进制文件。

## 工作目录初始化

在使用hyperbench时，需要使用一个独立的工作目录。所有hyperbench的操作都需要在工作目录（benchmark）同级目录下执行，不可以在目录之外或者是子目录中，否则会产生异常。您可以通过以下命令初始化出一个工作目录。

```text
# 创建test空目录
mkdir test && cd test
# 初始化目录
hyperbench init
```

该命令会将hyperbench使用过程中需要使用的配置文件等进行初始化，得到的目录结构如下

```text

test
└── benchmark     # 预先创建的测试用例
```

其中benchmark目录下主要是各种预先创建的一些测试用例，目录结构为

```text
benchmark
├── hyperchain       # hyperchain测试用例
│	├── azolla          # azolla测试用例
│	├── evmType         # evm合约类型测试用例
│	├── extraID         # transfer extra参数测试用例
│	├── hvmContract     # hvm 合约用例
│	├── hvmSBank        # hvmSBank 合约用例
│	├── invokeExample   # solidity合约 setHash用例
│	├── local           # 单机压测转账例子
│	├── remote          # 分布式压测转账用例
│	└── remote-evm      # 分布式压测solidity setHash用例
├── fabric           # fabric测试用例
│	├── Sacc            # go合约存取用例（key，value）
│	└── example         # go合约转账、查询、删除用例
├── ethereum         # ethereum测试用例  
│	├── invoke          # evm合约存取用例（key，value）
│	└── transfer        # 账户转账用例
├── xuperchain       # xuperchain测试用例  
│	├── evmContract     # evm合约increase、get用例
│	└── goContract      # go合约increase、get用例
```

对于每一个例子，目录下会有对应平台的网络配置目录（区块链SDK配置)、测试脚本以及压力测试相关的文件（合约、配置文件等）。例如，benchmark/hyperchain/local目录下的文件结构如下所示：

```text
local
├── hyperchain    # hyperchain的网络配置目录
│   └── hpc.toml  # hyperchain gosdk对应的网络配置文件
├── script.lua    # 测试脚本
└── config.toml   # 压测相关配置
```

benchamrk/hyperchain/remote-evm目录下的文件结构如下所示：

```text
remote-evm
├── contract                         # 合约目录
│   ├── evm                          # 具体合约类型相关目录
│   │   ├── source_solc_SetHash.bin  # 合约bin文件
│   │   ├── source_solc_SetHash.solc # 合约源文件
│   │   └── source_solc_SetHash.abi  # 合约abi文件
├── hyperchain                       # hyperchain的网络配置目录
│   └── hpc.toml                     # hyperchain gosdk对应的网络配置文件
├── script.lua                       # 测试脚本
└── config.toml                      # 压测相关配置
```



## 测试配置config.toml

初始化完成后就可以对区块链进行压力测试了，使用hyperbench进行压力测试需要编写toml格式的测试配置。

运行压测时使用start子命令，指令格式如下：

```text
hyperbench start path/to（/configuration）
# 需要说明的是这里的path/to是指测试例子的文件夹的路径
# path/to/configuration指的是具体的配置文件路径，这两种方式都支持
# 若不指明配置文件名，则需要保证压测目录下有名为config的配置文件

# 例如有一个本机压测转账的例子，在当前hyperbench目录下的/benchmark/local
# 那么指令就是
hyperbench start benchmark/local
或
hyperbench start benchmark/local/config.toml
```

在start子命令运行完成后我们可以看到命令行提示如下：

```none
[ctrl][NOTIC] 14:17:58.707 controller.go:113 finish 
```

在运行压测过程中，可以在命令行中看到压力测试的发送情况，以及最终的一个远程统计结果。

压力测试配置主要包括以下几个部分的配置项：

- **engine**  主要包括压测数据相关的配置

- **client** 主要包括压测平台及需要用到的文件、参数配置

- **recorder** 主要包括压测数据统计结果输出相关的配置

### **engine**

engine是用来进行压力控制的，即发送的压力数。

| 参数名      | 概述              | 类型     | 示例           |
|----------|-----------------|--------|--------------|
| rate     | 每秒发送的交易数量（压力设置） | number | 100          |
| duration | 压测持续时间          | string | "3m" (表示3分钟) |
| cap      | 启动的用户数量（同时也是最大并发数） | number | 100          |
| instant  | 每个批次发送的交易数      | number | 100          |
| wait     | 压力输出虚拟机获取最大等待时间 | string | "5ms"        |
下面详细说明一下各个项的具体用途：

1. **rate** ：压力控制的手段之一，压力引擎会根据rate指定的数值N，尝试每秒发送N笔交易。

2. **duration** ：压力控制的手段之一，描述压力引擎发送交易的持续时间，可以用"ns", "us" ( "µs"), "ms", "s", "m", "h"这几个单位来组合描述时长，例如持续90分钟的压测，既可以用"1h30m"也可以"90m"来表示。

3. **cap**：压力控制的手段之一，限制系统的最大并发数（当达到最大并发数的时候，系统不会继续提高并发压力，会出现实际发送的TPS和设定的TPS不一致的情况），同时这个数值限制系统最多同时保存多少份脚本控制的上下文。更通俗的来说，这个指标实际上是同时最多模拟多少个用户产生压力，每一个模拟用户的行为是根据测试脚本和各自的保存的上下文进行压力参数的生成。

4. **instant**: 压力控制的手段之一，压力引擎发送交易并不是一笔一笔发送的，而是一批一批发送的，instant指定每一批的交易数量。决定压力发送策略，即决定共分多少批次发送交易，此项可能会直接影响压力发送效果，为必须设置项，且不可为0，建议小于等于cap数。

5. **wait**: 获取压力输出最小单位虚拟机时的最大等待时间，即可以自由设置为block或非block形式。不设置或设置为0时即为非block形式获取。

### **client**

client是用来配置所使用的区块链连接、测试脚本、合约及客户端选项的。

| 参数名       | 概述                                 | 类型     | 实例                                |
|-----------|------------------------------------| ------ |-----------------------------------|
| script    | 指定测试脚本的路径                          | string | "benchmark/remote-evm/script.lua" |
| type      | 区块链类型                              | string | "hyperchain"                      |
| config    | 区块链sdk配置路径                         | string | "benchmark/remote-evm/hyperchain" |
| contract  | 要测试的合约文件夹的路径（作为参数传给deployContract) | string | "benchmark/remote-evm/contract"   |
| args      | 合约参数路径                             | []string | ["A", "123"]                      |
| plugin    | 区块链适配插件路径                          | string | "hyperchain.so"                   |
**【注意】**这里的配置的路径是**相对工作目录的路径或者绝对路径**，不是相对于config.toml的路径。

下面详细说明一下各项的用途：

1. **script** ：系统会根据script指定的脚本的内容将其拼接成一个完整的lua测试脚本，测试脚本的编写方法具体参见“测试脚本编写”部分的说明。

2. **type** ：标识所使用的区块链网络的类型，系统会根据你使用的type来进行适配层的选择，目前系统支持hyperchain、fabric、eth（ethereum）以及xuperchain。。

3. **config** ：连接区块链网络配置文件目录，config指向的是被测网络的配置文件以及keystore目录，hyperchain详细配置方案请参见hyperchain的go sdk文档，ethereum和xuperchain的网络配置只需要配置节点ip与端口。但需要注意的是ethereum的ip地址需要指明rpc类型，例如http://localhost或者ws://localhost。fabric的网络配置文件较为复杂，详见fabric官方文档。类似xuperchain和ethereum必须要有代币账户类型的区块链，必须在config/keystore目录下配置有代币的账户文件，具体配置细节可参考benchmark内readme。

4. **contract**  ：系统会根据contract项指向的目录下的文件结构，初始化合约的初始化（当contract指向一个无效路径时，例如配置为空字符串、不配置或者是指向一个无效路径时，不会初始化合约，只能正常执行转账），具体请参见”合约初始化“部分的说明进行contract目录的组织。（__当测试的合约为fabric合约时，合约代码需要放在gopath目录下，配置的contract合约目录路径需要是相对于gopath/src的路径__）

5. **args**：部署合约需要用到的参数路径，系统会根据指定的参数进行合约部署，一般在部署合约时不用指定参数，但是当部署一些特殊合约时可能需要用到。例如，部署fabric合约时，需要调用合约的init方法对合约进行初始化，并传入相应的init参数，这些都可通过args进行配置。

6. **plugin**: 编译后区块链适配插件的路径，支持绝对路径和相对路径。若此路径无效或者未配置，则系统不能正常启动。

#### client下有options配置项，用于配置客户端选项。每个区块链平台有不同的options配置项，具体请阅读第四章。


**recorder**

recorder用来配置压力统计结果输出相关的配置。目前压力统计结果输出有两种方式：

- csv 格式，以csv格式输出时需要指定csv格式文件输出的目录

- log 格式，以log格式输出时，需要配置log文件相关的配置

对于以csv格式输出，在recoder下有csv配置项，用于配置csv文件相关的配置。

| 参数名 | 概述              | 类型     | 实例      |
| --- | --------------- | ------ | ------- |
| dir | csv格式输出的csv文件目录 | string | "./csv" |

下面详细说明一下各项的用途：

1. **dir** ：用于指定csv文件的存放目录，当配置了`recorder.csv` ，但是没有配置`recoder.csv.dir` 的值时，会将`./csv` 设为存放csv文件的默认路径。

对于以log格式输出，在recoder下有log配置项，用于配置log格式输出相关的配置。

| 参数名   | 概述              | 类型     | 实例       |
| ----- | --------------- | ------ | -------- |
| level | 日志输出级别          | string | "NOTICE" |
| dir   | log格式输出的log文件目录 | string | "./log"  |
| dump  | 是否导出到文件         | bool   | "true"   |

下面详细说明一下各项的用途：

1. **level**：用于指定log文件输出的日志级别，当没有指定值，或者指定的是非法日志级别时，会将`NOTICE` 设为默认的日志级别。从高到低，日志级别依次为：CRITICAL、ERROR、WARNING、NOTICE、INTO、DEBUG。

1. **dir** ：用于指定log文件的存放目录，当没有指定值，或指定的空值时，会将`./log` 设为存放log文件的默认路径。

1. **dump** ：用于指定是否将log日志导出到文件。当没有指定值，或者指定的值为`false` 时，日志只在控制台输出，不导出到文件；当值为`true`时会将日志导出到文件。

对于需求使用influxdb进行数据可视化输出，在recorder下有influxdb配置项，用于配置influxdb服务的url、数据库名、用户名及密码。
需要配置`recorder.influxdb`配置项

| 参数名      | 概述                     | 类型     | 实例                      |
|----------|------------------------|--------|-------------------------|
| url      | influxdb服务的url，包括ip与端口 | string | "http://127.0.0.1:8086" |
| database | 指定influxdb内数据库名称       | string | "hyperbench"            |
| username | 访问influxdb数据库的用户名      | string | "admin"                 |
| password | 访问influxdb数据库的密码       | string | "admin"                 |

## 合约初始化

系统会根据测试计划`config.toml`中`client.contract`项所指定的目录下的目录结构进行规则匹配，从而初始化测试所使用的合约，对于hyperchain的测试来说，初始化的优先级依次如下：

1. EVM solidity合约 

2. HVM java合约初始化

3. 无法识别，不初始化合约

### hyperchain EVM合约初始化

如果你希望测试solidity编写的合约，请在`contract` 指定的目录下创建一个名为`evm`的目录，按照你所希望的初始化方式组织`evm`目录下的文件结构，初始化优先级依次如下（注意*表示任意字符串都可以）：

1. 如果希望使用已经部署好的合约，那么请将合约的ABI存放到扩展名为`abi`的文件，存放合约地址到扩展名为`addr`文件。

1. 如果希望系统帮你部署合约，那么请将编译好的二进制存放到扩展名为`bin`的文件，并将ABI文件存放到扩展名为`abi`的文件。

1. 如果希望系统帮你部署合约，但是本地又没有编译好合约，那么可以将合约源码存放到扩展名为`solc`的文件中， 注意这是一种不推荐的初始化方式，因为系统有可能无法正确编译你的合约从而导致初始化失败。

目前已经支持solidity的所有基本类型和他们的数组、切片类型，具体使用方式请参考`benchmark/evmType/script.lua`中的脚本编写方式。

**示例**

```text
# 如果目录结构是这样子的，那么会使用addr文件中的地址，和abi文件中的ABI，不会重新部署合约
contract
└── evm                      # 使用evm合约
    ├── contract.abi         # 合约abi文件
    └── contract.addr        # 合约addr文件

# 如果目录结构是这样子的，那么会使用bin文件中的二进制进行部署，和abi文件中的ABI
contract
└── evm                      # 使用evm合约
    ├── contract.abi         # 合约abi文件
    └── contract.bin         # 合约bin文件

# 如果目录结构是这样子的，那么会使用bin文件中的二进制进行部署，和abi文件中的ABI
contract
└── evm                      # 使用evm合约
    └── contract.solc        # solidity合约 
```

### hyperchain HVM合约初始化

如果你希望测试hyperchain的HVM合约，请在`contract`指定的目录下创建一个名为`hvm`的目录，按照你所希望的初始化方式，组织`hvm`目录下的文件，初始化优先级依次如下：

1. 如果你希望使用已经部署了的合约，你可以将合约地址存到扩展名为`addr`的文件中，并且将合约的abi放置到扩展名为`abi`的 文件中。

1. 如果希望系统帮你部署合约，那么你可以将合约编译成的jar包存放到目录下，扩展名为`jar`，并且将合约的abi放置到扩展名为`abi`的文件中。

**示例**

```text
# 如果目录是这样子的，那么会使用addr文件中的地址
contract
└── hvm                          # 使用jvm合约
    ├── contract.abi             # 合约abi文件
    └── contract.addr            # 合约addr文件

# 如果目录是这样子的，那么会部署合约
contract
└── hvm                          # 使用jvm合约  
    ├── contract.abi             # 合约abi文件
    └── contract.jar             # 合约jar包
```

## 测试脚本编写

在测试配置中的`script`项指定了测试脚本的路径，编写测试脚本是使用定制测试逻辑的手段。系统使用__**lua5.1语法**__进行脚本编写，需要编写者大概了解lua的语法。

### **简单使用**

要简单地编写一个测试脚本，只需要创建一个lua文件，然后像下面这样实现一个`Run`函数即可，其中testcase为提供的测试示例，通过这个测试示例可以调用提供的插件，并实现相应的钩子函数，函数会在运行过程中被调用执行。

```lua
local case = testcase.new()

function case:Run()
    local ret = case.blockchain:Transfer({
        from = "0",  -- account中别名为0的账户为转账的from，以string形式入参
        to = "1", -- account中别名为1的账户为转账的to，以string形式入参
        amount = 0, -- 转账金额
        extra = tostring(case.index.tx),-- 设置转账交易的extra字段，以string形式入参
    })
    return ret
end

return case
```



```lua
local case = testcase.new()

function case:Run()
    local ret = case.blockchain:Invoke({
        func="setHash", -- 调用的合约方法
        args={tostring(case.index.tx),
              tostring(case.index.worker)} -- 合约方法需要的参数列表。以string的形式入参
    })
    case.blockchain:Confirm(ret)
    return ret
end

return case

```

要注意的地方有：

1. testcase为hyperbench提供的测试case，需要创建一个新的实例，实现其中提供的钩子函数，最后将这个实例返回。

1. `Run`返回值要是invokeContract或者transfer的返回值。

来看一个evm合约的例子，下面是一个sethash的简单合约：

```text
contract SetHash {
	mapping(string => string) hashMap;
    function setHash(string key, string value) returns(string){
        hashMap[key] = value;
        return (value);
    }
    function getHash(string key) returns(string){
        return (hashMap[key])
    }
}
```

对于这么一个合约，如果希望测试合约中的`setHash`函数，那么我们可以编写这样一个测试脚本`script.lua`：

```lua
local case = testcase.new() -- 创建测试实例

function case:Run() -- 实现钩子函数Run函数
    local ret = case.blockchain:Invoke({
        func="setHash", -- 调用的合约方法
        args={tostring(case.index.Tx),
              tostring(case.index.Worker)} -- 合约方法需要的参数列表。以string的形式入参
    }) -- 调用case中提供的插件blockchain中的Invoke方法用于执行合约
    case.blockchain:Confirm(ret) -- 调用case中提供的查看blockchain中的Confirm函数用于查询交易回执
    return ret
end

return case

```

`case.index.Tx` 是用来标记这个是当前压力机发送的第几个交易，从0开始。`case.index.Worker` 是用来标记当前交易是使用的第几个压力机发送的，从0开始。这两个组合起来可以唯一标识一笔交易。

将上面这一段代码用于测试，就可以以这样的方式进行合约`setHash`函数的测试，每一次调用时，key是一个全局的编号，index是一个本地的编号，假设在测试配置中是单机压测，那么最后插入的数据可能是这样的：

```text
key    value
0       1
1       1
2       1
3       1
4       1
5       1
6       1
7       1
8       1
9       1
10      2
11      2
12      2
13      2
14      2
15      2
16      2
17      2
18      3       
19      2
```

**【注意】**系统无法保证模拟用户之间的执行次序，有的模拟用户的交易耗时长，有的模拟用户的交易耗时短，这就很有可能出现在测试期间内，每个模拟用户执行交易的次数不一致且差距比较大的情况，这都是无法预计的。

### 压测前置操作

上一小结中我们主要涉及的是`Run`这个hook，这个hook是每一次引擎尝试产生压力所会使用的钩子，那么如果系统需要在压测之前做一些前置的预备操作要怎么做，应该怎么做呢？这里提供了两种前置操作的hook，强化测试的灵活性，这两个钩子分别是：

- **GetContext()** 仅会被调用一次，需要返回context

- **SetContext(context)** 会被所有的”虚拟用户“调用一次，需要setContext

这两个钩子函数，前者会被特定的某一个部署合约的“虚拟用户”所调用，后者则会被每一个“虚拟用户”所调用。

下面举个例子，我们预先调用setHash，然后在压测中使用getHash，注意这里testplan配置中的`confirm`项要置为true才可以打印出内容。

```text
// 假设我们使用的是下面这个solidity合约
contract SetHash {
	mapping(string => string) hashMap;
    function setHash(string key, string value) returns(string){
        hashMap[key] = value;
        return (value);
    }
    function getHash(string key) returns(string){
        return (hashMap[key])
    }
}
```

```lua
local case = testcase.new()

function case:BeforeGet()
	-- 在这之前合约的信息会被初始化好
	for i = 1, 100 do
		case.blockchain:Invoke({
			func="setHash",
	        args={tostring(i),
	              tostring(i)}
		})
	end
end 

function case:Run()
	local num = (idx % 100) + 1
	local ret = case.blockchain:Invoke({
        func="getHash", 
        args={tostring(num)} 
    }) 
    ret = case.blockchain:Confirm(ret) 
	print(ret.ret[1])
    return ret
end

return case

```

### **钩子函数**

在脚本中，所有需要与虚拟机进行的交互被统一封装成一致的钩子函数，在对应的虚拟机中有各自的实现和调用方式，选择lua虚拟机的实现方式是使用`testcase.new()` 创建lua虚拟机的实现示例，然后根据需要实现里面的钩子函数。hyperbench提供了一些钩子函数，绝大部分情况下，使用者只需要使用到`Run` 函数。

| 函数名          | 说明                      |
|--------------| ----------------------- |
| BeforeDeploy | 部署合约前调用（master调用一次）     |
| BeforeGet    | 生成上下文前调用（master调用一次）    |
| BeforeSet    | 设置上下文前调用（每个cap调用一次）     |
| BeforeRun    | 运行前调用（每个worker运行前都会调用一次） |
| Run          | 压力测试运行的函数（每次发交易时调用）     |
| AfterRun     | 运行后调用（每个worker运行后都会调用一次） |

### Lua API

**blockchain API**

在脚本中，所有需要与区块链系统进行的交互被统一封装成一致的API，在适配层各自实现，具体选择哪种blockchain实现的方式是根据`type` 所指定的区块链系统来选择的。blockchain实现了一些函数，绝大部分情况下，使用者只需要使用到`Invoke`和`Transfer`两个函数：

| 函数名            | 参数                                                                      | 返回值                 | 说明                                                    |
| -------------- |-------------------------------------------------------------------------|---------------------|-------------------------------------------------------|
| DeployContract | -                                                                       | -                   | 部署合约                                                  |
| Invoke         | {func: string 函数名args: []string 参数列表}                                   | userdata: result    | 调用合约函数，返回的userdata如何使用请参见下一小节                         |
| Transfer       | {from: string 账户别名to: string 账户别名amount: number数额extra: string extra字段} | userdata: result    | 转账，返回值的userdata如何使用请参见下一小节, 默认会使用from别名对应的account进行签名 |
| Confirm        | userdata: result                                                        | userdata: result    | 查询交易回执，返回值的userdata如何使用请参见下一小节                        |
| Query          | {func: string 函数名args: []string 参数列表}                                   | interface{}         | 进行一些查询操作，属于预留接口                                       |
| Option         | map[string]interface{}                                                  | -                   | 设置区块链客户端相关参数                                          |
| GetContext     | -                                                                       | string：json格式的上下文   | 用来生成一个客户端实例的上下文，上下文内容主要是合约的一些信息，比如合约地址等等              |
| SetContext     | context: string getContext的返回值                                          | -                   | 同步上下文                                                 |
| ResetContext   | -                                                                       | -                   | 重置上下文                                                 |
| Statistic      | {from: 链起始状态，to: 链结束状态}                                                 | userdata: statistic | 用来计算两个链状态之间的tps，bps                                   |
| LogStatus      | -                                                                       | ChainInfo，error     | 链上的状态，包括区块高度、交易数和时间戳                                  |
**result 结构体**

result是一种特殊的userdata，脚本的`Run`函数需要以这种类型作为返回值，在系统中有四种方式可以产生result类型的userdata

- **blockchain:Invoke**返回值

- **blockchain:Transfer** 返回值

- **blockchain:Confirm** 返回值

- result.new()的返回值（正常来说测试脚本编写时不会使用到这个，只有在脚本中进行区块链平台适配时才使用这个）

对于result类型的userdata，实现了以下函数：

| 字段名         | 类型            | 说明                                                                                                                                |
| ----------- | ------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| uid         | uid: string   | 获取交易的uid，使用hyperchain时将获取txHash                                                                                                   |
| confirmTime | time: number  |                                                                                                                                   |
| label       | label: string | 本地统计数据将根据label进行分类统计，一般使用的所调用合约函数的名称                                                                                              |
| ret         | ret: table    | 将结果通过json.Marshal然后在json.Unmarshal进行类型裁剪之后产生的对应的table，对于evm各种result形式请参见benchmark/evmType/script.lua中的使用例子，该函数所转换出来的结果有一些特殊类型会有问题 |

还有一些API一般只会在使用脚本进行适配时使用，这里也罗列出来：

| 函数名         | 类型           | 说明         |
| ----------- | ------------ | ---------- |
| uid         | uid: string  | 变更uid      |
| sendTime    | time: number | 变更发送时间(纳秒) |
| buildTime   | time: number | 变更构造时间(纳秒) |
| confirmTime | time: number | 变更确认时间(纳秒) |
| confirm     | -            | 变更状态为确认    |
| failure     | -            | 变更状态为失败    |
| success     | -            | 变更状态为成功    |
| unknown     | -            | 变更状态为未知    |

**statistic 结构体**

`statistic`是用来组织链上统计TPS时所使用的操作，一般来说除了使用lua做底层适配之外，不需要使用这边的API，statistic实例的获取方法是:

- **blockchain:statistic({from, to})**   获取适配层实现统计所返回的实例，可能是在go里面生成的也可能是在lua里面生成

- **statistic.new({from, to})**  除非你要实现一个适配层，不然的话不建议这么做

| 字段名      | 类型     | 说明    |
|----------| ------ |-------|
| start    | Number | 开始时间  |
| end      | Number | 结束时间  |
| blockNum | Number | 总区块数  |
| txNum    | Number | 总交易数  |
| TPS      | Number | 压测tps |
| BPS      | Number | 压测bps |


**【注意】**

- 合约调用使用别名为"0"的账户来进行签名，转账使用`from`参数作为别名指定的账户进行签名。

- 在使用测试性能时尽量避免使用与账号强相关的case，或者把合约中一些账号控制的部分给注释掉，这个更多的应该是和功能测试相关的工作。

- 建议account配合testplan中的`keystore`配置项使用。



## 账户仓库keystore

账号仓库`keystore`在config.toml中是一个可选配置。但对于某些链必须要有可以使用代币但账户，例如Ethereum和Xuperchain，如果用户希望使用某几个特定账户，那么可以指定`keystore`。使用`keystore`的时候需要注意的事情有以下几点：

- 账户文件必须能够被适配层所读取，对于hyperchain，账户文件必须是由gosdk产生的非加密账号，类型需要和config.toml `client.option.type` 字段配置的一致

- 注意即使文件目录下文件结构一致也可能出现特定账户别名发生不同的情况，因此，对于hyperchain，keystore中的账户也可以用地址映射到该账户

**使用范例**

Transfer

```text
// keystore指定的目录下的账号文件之一
{
	"address": "0xcf8dc52bab9775e3df68d7e2f82f52a382bf7706",
	"algo": "0x03",
	"encrypted": "5ebd455ff5f7db8d59f1f8712fa48e28b958c03f265415225086189f3d74a489",
	"version": "2.0",
	"privateKeyEncrypted": false
}
```

```text
// keystore指定的目录下的账号文件之一
{    
	"address": "0xb0249132126707c2b07aa165cd32927c10396fba",    
	"algo": "0x03",    
	"encrypted": "2682b42c285640eefb7ad2e6131e7a5fd901349516dae37e8b82497056f98776",    
	"version": "2.0",    
	"privateKeyEncrypted": false
}
```

```lua
-- 转账逻辑
local case = testcase.new()

function case:Run()
    local ret = case.blockchain:Transfer({
        from = "0",
        to = "1",
        amount = 0,
        extra = tostring(case.index.tx),
    })
    return ret
end
-- 这样就可以完成一次从keystore指定的某一个特定账号的到另一个账号的转账
return case

```

# 第三章 其他使用说明

## 指令概览

先通过`hyperbench --help`来看一下hyperbench有哪些指令：

```none
Usage:
  hyperbench [command]

Examples:
hyperbench --doc ./doc (generate document to specify path)

Available Commands:
  help        Help about any command
  init        init a stress test dir
  new         initialize a test plan
  start       start a benchmark
  version     get code version
  worker      start as a worker server 

Flags:
      --debug        enable debug mode
      --doc string   use to create doc and specify the doc path
  -h, --help         help for hyperbench

Use "hyperbench [command] --help" for more information about a command.
```

所有hyperbench的子命令都可以使用hyperbench的flags。

## 执行压测

使用hyperbench执行压测的命令使用说明如下：

```text
start a benchmark

Usage:
  hyperbench start [flags]

Examples:
hyperbench start benchmark/transfer

Flags:
  -h, --help   help for start

Global Flags:
      --debug        enable debug mode
      --doc string   use to create doc and specify the doc path
```

hyperbench同样支持对fabric进行压力测试，这需要自行编写测试文档，在benchmark目录下同样有fabric的测试用例可以使用，例如：

```text
hyperbench start benchmark/fabricSacc
```

## 启动worker

使用hyperbench启用worker服务的命令说明如下：

```text
start as a worker server

Usage:
  hyperbench worker [flags]

Examples:
hyperbench worker

Flags:
  -h, --help       help for worker
  -p, --port int   port of worker (default 8080)

Global Flags:
      --debug        enable debug mode
      --doc string   use to create doc and specify the doc path
```

例如，使用hyperbench在8081端口启动一个worker服务：

```text
hyperbench worker -p 8081
```

worker启动后可以看到如下命令行提示：

```text
[GIN-debug] POST   /set-nonce                --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func1 (3 handlers)
[GIN-debug] POST   /upload                   --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func2 (3 handlers)
[GIN-debug] POST   /init                     --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func3 (3 handlers)
[GIN-debug] POST   /set-context              --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func4 (3 handlers)
[GIN-debug] POST   /do                       --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func5 (3 handlers)
[GIN-debug] POST   /checkout-collector       --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func6 (3 handlers)
[GIN-debug] POST   /teardown                 --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func7 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8081
```

启动了worker服务后，在使用分布式压力测试时，在config.toml中`engine.urls` 配置项中配置远程worker的ip和端口即可。

# 第四章 适配平台配置说明

## Hyperchain

如果您想进行hyperchain的测试，请先阅读此手册。
例如，在`benchmark/hyperchain/hvmSBank`测试用例中。

### 网络连接

hyperchain详细配置方案请参见hyperchain的go sdk文档。Hyperchain网络连接配置文件在`hvmSBank/hyperchain/hyperchain.toml`路径。

### 合约
系统会根据测试计划`config.toml`中`client.contract`项所指定的目录下的目录结构进行规则匹配，从而初始化测试所使用的合约，对于hyperchain的测试来说，初始化的优先级依次如下：

1. EVM solidity合约

2. HVM java合约初始化

3. 无法识别，不初始化合约

#### EVM合约初始化

如果你希望测试solidity编写的合约，请在`contract` 指定的目录下创建一个名为`evm`的目录，按照你所希望的初始化方式组织`evm`目录下的文件结构，初始化优先级依次如下（注意*表示任意字符串都可以）：

1. 如果希望使用已经部署好的合约，那么请将合约的ABI存放到扩展名为`abi`的文件，存放合约地址到扩展名为`addr`文件。

1. 如果希望系统帮你部署合约，那么请将编译好的二进制存放到扩展名为`bin`的文件，并将ABI文件存放到扩展名为`abi`的文件。

1. 如果希望系统帮你部署合约，但是本地又没有编译好合约，那么可以将合约源码存放到扩展名为`solc`的文件中， 注意这是一种不推荐的初始化方式，因为系统有可能无法正确编译你的合约从而导致初始化失败。

#### HVM合约初始化
如果你希望测试hyperchain的HVM合约，请在`contract`指定的目录下创建一个名为`hvm`的目录，按照你所希望的初始化方式，组织`hvm`目录下的文件，初始化优先级依次如下：

1. 如果你希望使用已经部署了的合约，你可以将合约地址存到扩展名为`addr`的文件中，并且将合约的abi放置到扩展名为`abi`的 文件中。

2. 如果希望系统帮你部署合约，那么你可以将合约编译成的jar包存放到目录下，扩展名为`jar`，并且将合约的abi放置到扩展名为`abi`的文件中。

### client.options

| 参数名             | 概述                        | 类型     | 实例                                    |
|-----------------|---------------------------| ------ | ------------------------------------- |
| keystore        | 账户仓库路径                    | string | "benchmark/remote-evm/keystore/ecdsa" |
| type            | 账户签名类型                    | string | "ECDSA"、"SM2"（默认ECDSA）                |
| request         | 交易连接类型，rpc或grpc，必需配置项     | string | "rpc"            |
| crosschain      | 合约调用时是否使用跨链交易             | string | "true"                           |
| simulate        | 是否使用simulate模式发送交易        | bool   | "false"                               |
| vmtype          | 合约类型，使用bvm时需要指明，其他无需指明    | string   | "bvm"            |
| fvmadvancedtype | 使用fvm合约时，是否使用advancedtype | bool   | "false"                               |

下面详细说明一下各项的用途：

1. **keystore** ：如果需要使用指定的账号，可以配置keystore，系统会读取keystore指向的目录下所有文件（不递归，只读取第一级文件），对于hyperchain，每个文件表示一个账号，文件名无所谓，但是文件内容必须是由hyperchain的go SDK生成的sign指定的类型的account json文件，否则无法正常识别。

2. **type** ：系统会根据这个标识来判断使用哪种类型的账户进行交易的发送，对于hyperchain，目前支持sm2和ecdsa两种账户，对大小写不敏感。

3. **request** ：用于配置交易连接类型，rpc或grpc，必需配置项。有rpc和grpc的区别。

4. **crosschain** ：用于配置使用合约时是否使用跨链交易。

5. **simulate** ：用于配置simulate模式交易。

6. **vmtype** ：用于配置合约虚拟机类型，仅需在使用bvm时需要配置。

7. **fvmadvancedtype** ：用于配置使用fvm合约时，是否使用advanced类型。

## Fabric
如果您想进行Fabric的测试，请先阅读此手册。
例如，在`benchmark/fabric/Sacc`测试用例中。

### 合约
再部署合约时，请将合约放置在正确未知，例如`benchmark/fabric/Sacc/contract`。在部署合约时，fabric go-sdk会自动将合约路径加上$GOPATH/src，所以在配置文件中配置合约路径时，请参考`benchmark/fabric/Sacc/config.toml`。
### client.options
| 参数名         | 概述                   | 类型     | 实例                                    |
|-------------|----------------------| ------ | ------------------------------------- |
| instant     | fabric、xuperchain初始化账户数量               | number | 10                                    |
| channel     | fabric网络对应的channelID | string | "mychannel"                           |
| option.MSP  | fabric网络是否启用MSP      | bool   | "false"                               |

下面详细说明一下各项的用途：

1. **instant** ：初始化账户数量。

2. **channel** ：用于指定fabric网络中对应的channelID。

3. **option.MSP** ：用于配置fabric网络是否启用MSP服务。

## Ethereum
如果您想进行Ethereum的测试，请先阅读此手册。
例如，在`benchmark/eth/invoke`测试用例中。

### 账户
至少需要配置一个有代币的账户在正确路径下，需要使用ethereum生成的账户文件。若您需要测试transfer，则至少需要配置两个账户文件。


### 网络连接
请配置ethereum网络连接方式、ip以及端口。连接方式包括：http或者websocket。
例如：
```
[rpc]
node = "http://localhost"
port = "8545"
```
或
```
[rpc]
node = "ws://localhost"
port = "8546"
```
### 合约
如果您需要部署合约，请将合约编译后的abi和bin文件放置于正确路径下，例如`benchmark/eth/invoke/contract`，solidity文件非必须。

### 注意项

此处为使用ethereum压测时的一些注意项。

#### NONCE值
由于ethereum的nonce设置必须连续且递增才能使交易生效，在测试transfer过程中时，FROM账户需保持不变，否则将导致交易失败。

#### 账户地址
在测试transfer测试用例时，需要设置FROM和TO账户地址。地址名为账户文件最后一段编号，例如，文件名为`UTC--2021-11-08T06-39-32.219546000Z--74d366e0649a91395bb122c005917644382b9452`,它的地址为`74d366e0649a91395bb122c005917644382b9452`。
### client.options
| 参数名         | 概述           | 类型     | 实例    |
|-------------|--------------|--------|-------|
| keypassword     | ethereum账户密码 | string | "111" |


下面详细说明一下各项的用途：

1. **keypassword** ：ethereum账户密码，若未设置密码也可以不配置。

## Xuperchain
如果您想进行Fabric的测试，请先阅读此手册。
例如，在`benchmark/xuperchain/evmContract`测试用例中。
### 账户
至少需要配置一个有可使用代币的账户在main路径下，例如`benchmark/xuperchain/evmContract/xuperchain/keystore/main`。
账户文件需为xuperchain生成的账户文件，且在配置时需要分开路径，main路径下放置主账户文件，若需要配置其余账户，则需类似main路径，配置为不同路径。

### 网络连接
请配置部署的xuperchain网络的ip以及端口在`benchmark/xuperchain/evmContract/xuperchain/xuperchain.toml`
例如：
```
[rpc]
node = "127.0.0.1"
port = "37101"
```

### 合约
当前适配合约支持go以及evm合约，请通过文件路径指明合约类型，例如`benchmark/xuperchain/evmContract/contract/evm`中最后的evm路径。

#### EVM
如果您需要部署的合约为evm，请将编译后的abi以及bin文件在`benchmark/xuperchain/evmContract/contract/evm`路径下，solidity文件为可选。

#### GO
如果您需要部署的合约为evm，请将go合约编译后的文件在例如`benchmark/xuperchain/goContract/contract/go`路径下，go文件为可选。

编译合约请参考 [here](https://www.bookstack.cn/read/XuperChain-5.1-zh/ee1cca974bbc0699.md).


### 注意项

此处为使用xuperchain测试的一些注意项。

#### UTXO模型
如果在transfer测试时使用同一账户作为FROM以及TO账户，极有可能会导致交易失败。因此，您可以在client.options下配置instant项，设置在压测前初始化的测试账户数量，初始化时会预先为测试账户转入一定的代币。具体的instant使用说明参考下文client.options栏。

#### 转账测试脚本
转账测试中的FROM和TO账户都可以配置为空。若FROM为空，main账户则作为FROM账户；若TO为空且配置instant项，则会从其中随机选择一个账户作为TO账户，若TO为空且为配置instant，则会生成一个账户作为TO账户。
若设置FROM与TO账户，但为无效账户，则与未配置相同效果。
由于xuperchain本身但设置，转账数量Amount必须大于0。
#### 合约调用参数

合约调用时参数必须成对设置，参考xuperchain官方合约使用文档，例如：
```lua
local case = testcase.new()

function case:Run()
    local result = self.blockchain:Invoke({
        func = "Increase",
        args = {{"creator","test"},{"key","test"}},
        -- Since the args parameter of the xuperchain contract call is in map[string]string format
    return result
end
return case
```
### client.options
| 参数名     | 概述        | 类型     | 实例 |
|---------|-----------|--------|----|
| instant | 初始化预设账户数量 | number | 10 |


下面详细说明一下各项的用途：

1. **instant** ：设置在压测前初始化的测试账户数量，初始化时会预先为测试账户转入一定的代币，供转账测试使用，可提高交易成功率。若设置数量较大，在初始化时可能需要时间相对较长，请耐心等待。


# 第五章 使用示例

## 单台压力机测试

使用benchmark中的local测试用例，单机对flato进行压力测试，其config.toml的配置如下：

```text
[engine]
rate = 20                            # 速率
duration = "20s"                     # 持续时间
cap = 10                             # 客户端虚拟机数量
instant = 5                          # 每个批次发的交易数
wait = "5ms"                         # 获取cap最大等待时间

[client]
script = "benchmark/local/script.lua"  # 脚本
type = "hyperchain"                    # 区块链类型
config = "benchmark/local/hyperchain"  # 区块链SDK配置路径
contract = "benchmark/local/contract"  # 合约目录路径
plugin = "hyperchain.so"               # 插件路径

[client.options] # 客户端选项
request = "rpc"
```

local的文件目录结构如下：

```text
local
├── hyperchain  # hyperchain的网络配置目录
│   └── hpc.toml  # hyperchain gosdk对应的网络配置文件
├── script.lua  # 测试脚本
└── config.toml # 压测相关配置
```

使用start子命令开始压力测试：

```bash
hyperbench start benchmark/local
```

在start子命令运行完成后我们可以看到命令行提示如下：

```text
[ctrl][NOTIC] 14:17:58.707 controller.go:113 finish 
```

## 分布式压力机测试

使用benchmark中的remote-evm测试用例，使用多台压力机对flato进行压力测试。

例如，使用172.0.1.10、172.0.1.11、172.0.1.12三台服务器进行压力测试，其中172.0.1.10、172.0.1.11作为worker，172.0.1.12作为master控制整个压力测试。

首先在172.0.1.10、172.0.1.11服务器的8081端口分别启动了一个worker。将hyperbench安装到172.0.1.10、172.0.1.11服务器上后，分别运行以下命令：

```text
hyperbench worker -p 8081
```

看到以下命令行，则表示worker启动完成。

```text
[GIN-debug] POST   /set-nonce                --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func1 (3 handlers)
[GIN-debug] POST   /upload                   --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func2 (3 handlers)
[GIN-debug] POST   /init                     --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func3 (3 handlers)
[GIN-debug] POST   /set-context              --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func4 (3 handlers)
[GIN-debug] POST   /do                       --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func5 (3 handlers)
[GIN-debug] POST   /checkout-collector       --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func6 (3 handlers)
[GIN-debug] POST   /teardown                 --> github.com/hyperbench/hyperbench/core/network/server.(*Server).Start.func7 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8081
```

然后在172.0.1.12服务器上，配置benchmark/remote-evm目录的config.toml，其配置如下：

```text
[engine]
rate = 20                            # 速率
duration = "20s"                     # 持续时间
cap = 10                             # 客户端虚拟机数量
instant = 5                          # 每个批次发的交易数
wait = "5ms"                         # 获取cap最大等待时间
urls = ["172.0.1.10:8081", "172.0.1.11:8081"]                 # 若不设置或者长度为0则在本地启动worker

[client]
script = "benchmark/remote-evm/script.lua"  # 脚本
type = "flato"                              # 区块链类型
config = "benchmark/remote-evm/hyperchain"  # 区块链SDK配置路径
contract = "benchmark/remote-evm/contract"  # 合约目录路径
plugin = "hyperchain.so"                    # 插件路径

[client.options] # 客户端选项
request = "rpc"

```

remote-evm的文件目录结构如下：

```text
remote-evm
├── contract                         # 合约目录
│   ├── evm                          # 具体合约类型相关目录
│   │   ├── source_solc_SetHash.bin  # 合约bin文件
│   │   ├── source_solc_SetHash.solc # 合约源文件
│   │   └── source_solc_SetHash.abi  # 合约abi文件
├── hyperchain                       # hyperchain的网络配置目录
│   └── hpc.toml                     # hyperchain gosdk对应的网络配置文件
├── script.lua                       # 测试脚本
└── config.toml                      # 压测相关配置
```

在master上使用start子命令开始压力测试：

```bash
hyperbench start benchmark/remote-evm
```

在start子命令运行完成后我们可以看到命令行提示如下：

```text
[ctrl][NOTIC] 14:17:58.707 controller.go:113 finish 
```