# Ethereum Test

If you want to customize your Ethereum test, please read this first.
For example, in `benchmark/ethereum/invoke` test case.
The keystore is in the `invoke/eth/keystore` directory, and the Ethereum network connection configuration is set in `invoke/eth/eth.toml`.
The contract files are placed in the `invoke/contrat` directory.

## Account

At least one account file with a token must be configured in the keystore, because the operation of Ethereum requires payment of tokens.
The account file must be an Ethereum-generated file with filename unchanged.
If you want to test Transfer, please set up at least two accounts.

## Network connection

Please set the way of rpc connection, ip and port of ethereum network in the eth.toml file.
example
```
[rpc]
node = "http://localhost"
port = "8545"
```

## Contract

If you want to deploy a contract, please place both `*.abi` and `*.bin` files compiled from the source solidity contract like those in the directory of `invoke/contract` (solidity file is optional).

## Attention

Here are some usage details.

### NONCE

Because the nonce value of Ethereum must be continuous and increasing, it should be noted that the FROM account of Transfer test must be the same. If the alternate FROM account is used otherwise, the nonce value of the alternate account will be invalid and the transaction will fail.

### AccountName
You need to set the FROM and TO account addresses in the Transfer test. The account address should be set to the address in the name of account file.
For example, if the account filename is `UTC--2021-11-08T06-39-32.219546000Z--74d366e0649a91395bb122c005917644382b9452`, its address is `74d366e0649a91395bb122c005917644382b9452`.