local index = 0
local case = testcase.new()
function case:Run()
    --print("run")
    index = index + 1
    local result = self.blockchain:invokeContract({
        func = "query",
        args = { "A" }
    })
    --print(result)
    return result
end

return case

