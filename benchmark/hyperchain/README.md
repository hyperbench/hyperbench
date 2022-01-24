# Hyperchain Test

If you want to customize your Hyperchain test, please read this first.
For example, in `benchmark/hyperchain/hvmSBank` test case.
The Hyperchain network connection configuration is set in `hvmSBank/hyperchain/hyperchain.toml`.
The contract files are placed in the `hvmSBank/contrat` directory.

## Network connection

The configuration of Hyperchain network is official, please refer to its official documentation for hyperchain network configuration of go sdk or you can refer to these test case.

## Contract

The system will deploy the contract according to the directory of the contract, and the initialization priority is as follows:
1. EVM solidity contract
2. HVM java contract
3. Unrecognized contract will not be deployed

### EVM
If you want to deploy solidity contract, please create a `evm` directory under `path/to/contract` like `evmType/contrat/evm`
The initialization priorities are as follows:
1. If you want to use the deployed contract, please set the ABI file as `*.abi`, and set the contract address to `*.addr`.
2. Compiled binary file `*.bin` and ABI file `*.abi` will do.
3. It can be the `*.solc` source code of the contract but it's not recommended because it may cause failure in compling contract.

### HVM
If you want to deploy hvm contract, please create a `hvm` directory under `path/to/contract` like `hvmSBank/contrat/hvm`
The initialization priorities are as follows:
1. If you want to use the deployed contract, please set the contract address to `*.addr` and set the ABI file of contract to `*.addr`.
2. If you want to deploy java contract, please set the JAR file compiled by the contract as `*.jar`, and set the ABI file of contract to `*.addr`.