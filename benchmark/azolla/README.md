# Azolla

If you are going to test Azolla, please use `self.blockchain:option({nonce=nonce})` to set nonce for Azolla.
If you do not set it or set a negative number, nonce will be random for each transaction. You should remember
to increace the nonce each time you set it.

example

```lua
self.blockchain:Option({nonce=0})
```

If you want to test invokeContract for Azolla, please use the `self.blockchain:option({account="random0"})` 
to set a random account as the default account name which will be used to sign transaction

example:
```lua
self.blockchain:Option({account="random0", nonce=0})
self.blockchain:Invoke("sethash", "123", "123")
```

If you want to test transfer for Azolla, from account must be set to a random one

example:
```lua
self.blockchain:Option({nonce=0})
self.blockchain:Transfer("random0", "1", 0)
```



