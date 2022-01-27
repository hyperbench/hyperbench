# Fabric Test

If you want to customize your Fabric test, please read this first.
For example, in `benchmark/fabric/example` test case.
The Fabric network connection configuration is set in example/fabric.
The chaincode files are placed in the `example/contrat` directory.

## Network connection

The configuration of Fabric network is official, please refer to its official documentation for fabric network configuration or you can refer to these test case.

## Contract

If you want to deploy a contract, please chaincode file in the directory of `example/contract`. In addition, gosdk of Fabric will add `GOPATH/src` to the path of chaincode file autometically, so you should refer to example/config.toml when you set the contract path.

## Attention

### Instant

The instant parameter in the client.options must be set.