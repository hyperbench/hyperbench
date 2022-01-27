# Xuperchain Test

If you want to customize your Xuperchain test, please read this first.
For example, in `benchmark/xuperchain/evmContract` test case.
The keystore is in the `evmContract/xuperchain/keystore` directory, and the xuperchain network connection configuration is set in `evmContract/xuperchain/xuperchain.toml`.
The contract files are placed in the evmContract/contrat directory.

## Account

At least one account file with a token must be configured in the keystore, because the operation of xuperchain requires payment of tokens.
The account file must be an xuperchain-generated file. The accounts files shoulde be seperated in directories under the keystore path, there must be a `main` directory which contains an accounts with enough token like `evmContract/xuperchain/keystore/main`.

## Network connection

Please set ip and port of xuperchain network in the eth.toml file.
example
```
[rpc]
node = "127.0.0.1"
port = "37101"
```

## Contract

The current adaptation supports evm and go contracts, please indicate the type of contract used by creating a directory named as the type of contract like the `evm` directory in `evmContract/contract/evm`.

### EVM

If you want to deploy a evm contract, please place both `*.abi` and `*.bin` files compiled from the source solidity contract like those in the directory of `evmContract/contract/evm` (solidity file is optional).

### GO

If you want to deploy a go contract, please place the compiled contract file from the source go contract like that in the directory of `goContract/contract/go` (go file is optional).
Details for compiling contract [here](https://www.bookstack.cn/read/XuperChain-5.1-zh/ee1cca974bbc0699.md).

## Attention

Here are some usage details.

### UTXO

If the same account is used to Transfer test, it is very likely that the transaction will fail because of UTXO model. Therefore, in the configuration of xuperchain, you can choose to initialize a certain amount of test accounts for transfer before the stress test.
You can set the "instant" configuration under client.options in the config.toml referring to `benchmark/xuperchain/transfer/config.toml`.

### Transfer Script

The FROM and TO parameters of the Transfer function can be set empty. If FROM is empty, the main account will be used as the FROM account. When TO is empty, if instant is set, an account will be randomly obtained as the TO account. Otherwise, a new account will be created as the TO account.
If FROM and TO are set but the account is invalid, it is the same as not set.
The Amount parameter must be be greater than 0, because xuperchain does not allow the transfer amount to be 0.

### Contract Invoke Parameter

The args parameter of the Invoke function in the script must be set in pairs, as shown in the example:

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



