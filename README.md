# HyperBench

hyperbench is a distributed stress testing tool, used to perform stress testing on blockchain platforms, written by go.

detail introduction : [white paper](https://upload.hyperchain.cn/HyperBench%E7%99%BD%E7%9A%AE%E4%B9%A6.pdf)

## Features

-  flexible:  provide programmable use case extension based on Lua script and user hook provided by virtual machine.

- high-efficiency: the virtual machine has a built-in Go blockchain client with a unified interface, which directly tests the blockchain system without additional application services

- distributed: support distributed test function, support the use of multiple presses to simultaneously test the blockchain system, simple and easy to use

## Quick Start

### Build from Source Code

```bash
# clone Hyperbench repository into the $GOPATH/src/github.com/meshplus/hyperbench directory:
git clone git@github.com:meshplus/hyperbench.git

# enable go modules: https://github.com/golang/go/wiki/Modules
export GO111MODULE=on

# build main program
make build

# copy build hyperbench program to $GOPATH/bin
cp hyperbench $GOPATH/bin
```

### Run HyperBench

Before start stess test, prepare a blockchain network(support Flato, Hyperchain, Fabric),  the gosdk and config file which are used to connect the network.

For example, run `local` benchmark in benchmark directory, send tx to hyperchain network for stress test. Actually, you can create new benchmark for you Application scenario, write Lua script and configure the config file.

1. init hyperbench work directory

```bash
# create empty directory test
mkdir test & cd test

# init hyperbench work directory
hyperbench init

# show example test case in benchmark
ls benchmark

```

2. set hpc.toml which is prepared to connnect hyperchain network with gosdk in local/hyperchain directory.

now, loca directory's structure is:

```bash
local
|_hyperchain   # set config file which is used to connect hyperchain
| |_hpc.toml   # hyperchain's gosdk will use the config file to connect hyperchain
|_script.lua   # lua script, execute transfer tx
|_config.toml  # the config file of hyperbench
```

3. start stress test.

```bash
# use benchmark/local as test case, star benchmark test
hyperbench start benchmark/local
```

## Contribution

Thank you for considering to help out with the source code! No matter it's a system bug report, a new feature purposal, or code contributing, we're all huge welcome.

Please check the contributing guide for the full details.

## License

Hyperbench is currently under Apache 2.0 license. See the LICENSE file for details.