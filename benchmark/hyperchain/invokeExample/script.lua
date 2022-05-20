local index = 0
local case = testcase.new()
function case:Run()
    index = index + 1
    local result = self.blockchain:Invoke({
        func = "setHash",
        args = { tostring(self.index.Tx), tostring(index) },
    })
    result=self.blockchain:Confirm(result)
    --print(result.Status)
    return result
end
return case



