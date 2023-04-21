# hyperbench

hyperbench is a distributed stress testing tool, used to perform stress testing on blockchain platforms, written by go.

detail introduction : [white paper](https://upload.hyperchain.cn/HyperBench%E7%99%BD%E7%9A%AE%E4%B9%A6.pdf)

## Features

-  flexible:  provide programmable use case extension based on Lua script and user hook provided by virtual machine.

- high-efficiency: the virtual machine has a built-in Go blockchain client with a unified interface, which directly tests the blockchain system without additional application services

- distributed: support distributed test function, support the use of multiple presses to simultaneously test the blockchain system, simple and easy to use

## Version

| hyperbench version | hyperbench-plugins version |
|--------------------|----------------------------|
| v1.0.5+            | v0.0.5                     |
| v1.0.4             | v0.0.4                     |
| v1.0.3             | v0.0.3                     |
| v1.0.2             | v0.0.2                     |
| v1.0.1             | v0.0.1                     |
| v1.0.0             | v0.0.1-alpha               |
## Attention
It should be noted that the current system uses go-plugin, go-plugin only supports macOS and Linux systems, and Windows systems are not yet supported.
## Quick Start
### install packr
```bash
# Go 1.16 and above
go install github.com/gobuffalo/packr/v2/packr2@latest
# Go 1.15 and below
go install github.com/gobuffalo/packr/packr2@latest
```
### Build hyperbench from Source Code

```bash
# clone hyperbench repository into the $GOPATH/src/github.com/hyperbench/hyperbench directory:
mkdir $GOPATH/src/github.com/hyperbench && cd $GOPATH/src/github.com/hyperbench
git clone https://github.com/hyperbench/hyperbench.git
cd hyperbench

# build main program
make build

# copy build hyperbench program to $GOPATH/bin
cp hyperbench $GOPATH/bin
```

### Build plugins from Source Code

```bash
# clone hyperbench-plugins into the $GOPATH/src/github.com/hyperbench/hyperbench-plugins directory:
cd $GOPATH/src/github.com/hyperbench
git clone https://github.com/hyperbench/hyperbench-plugins.git

# build hyperchain for example
cd hyperbench-plugins/hyperchain
make build
cp hyperchain.so ../../hyperbench/hyperchain.so
```

### Run hyperbench for Hyperchain

Before start stess test, docker and docker-compose are needed to be installed before preparing a hyperchian network.
1. start hyperchain network

```bash
# start hyperchain network
cd $GOPATH/src/github.com/hyperbench/hyperbench/benchmark/hyperchain
bash control.sh start
```

2. start stress test.

```bash
# use benchmark/hyperchain/local as test case, star benchmark test
cd $GOPATH/src/github.com/hyperbench/hyperbench
hyperbench start benchmark/hyperchain/local
```

### Run hyperbench for Fabric

Before start stess test, docker and docker-compose are needed to be installed before preparing a fabric network.
1. start fabric network

```bash
# start fabric network
cd $GOPATH/src/github.com/hyperbench/hyperbench/benchmark/fabric/example/fabric
bash deamon.sh
# if it's the first time for you to pull the docker images, please be patient
# and if a timeout bug occurs, please try again
```

2. start stress test.

```bash
# use benchmark/fabric/example as test case, star benchmark test
cd $GOPATH/src/github.com/hyperbench/hyperbench
hyperbench start benchmark/fabric/example
```
### Possible problems
You may encounter a similar issue where the main repository and plugins share different package versions due to automatic package version updates. This is required by go-plugin, there is no other way than to unify the version.
```text
[blockchain][ERROR] 10:37:45.853 blockchain.go:39 plugin failed: plugin.Open("./hyperchain"): plugin was built with a different version of package golang.org/x/sys/internal/unsafeheader
```

### Others
hyperbench is a blockchain testing tool, the corresponding examples provided in the benchmark are for reference only. When testing the corresponding blockchain platform, the use of the corresponding blockchain platform refers to the user manual of the corresponding platform.

## Contribution

Thank you for considering to help out with the source code! No matter it's a system bug report, a new feature purposal, or code contributing, we're all huge welcome.

Please check the contributing guide for the full details.

## License

hyperbench is currently under Apache 2.0 license. See the LICENSE file for details.