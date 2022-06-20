local nonce = 0
local case = testcase.new()
function case:BeforeDeploy()
    self.blockchain:Option({nonce=0})
end

function case:Run()
    -- transfer
    self.blockchain:Option({nonce=nonce})
    local result = self.blockchain:Transfer({
        From="random0",                   -- random account `random0`
        To="0",                         -- account in keystore
        Amount=0,
    })

    -- increase nonce
    nonce = nonce + 1
    -- print(result:result()[1]) -- 0x0
    return result
end

return case

