

local nonce = 0
local case = testcase.new()
function case:BeforeDeploy()
    self.blockchain:option({nonce=0})
end


function case:Run()

    -- transfer
    self.blockchain:option({nonce=nonce})
    local result = self.blockchain:transfer(
            "random0",                   -- random account `random0`
            "0",                         -- account in keystore
            0
    )

    -- invokeContract
    --if nonce == 0 then
    --    self.blockchain:option({account="random0"})
    --end
    --self.blockchain:option({nonce=nonce})
    --local result = self.blockchain:invokeContract(
    --        "setHash",                   -- random account `random0`
    --        tostring(idx),
    --        tostring(idx)
    --)

    -- increase nonce
    nonce = nonce + 1
    -- print(result:result()[1]) -- 0x0
    return result
end

return case

