local nonce = 0
local case = testcase.new()
function case:BeforeDeploy()
    self.blockchain:option({nonce=0})
end

function case:Run()
    -- transfer
    self.blockchain:option({nonce=nonce})
    local result = self.blockchain:Transfer({
        from="random0",                   -- random account `random0`
        to="0",                         -- account in keystore
        amount=0,
    })

    -- increase nonce
    nonce = nonce + 1
    -- print(result:result()[1]) -- 0x0
    return result
end

return case

