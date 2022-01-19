local index = 0
local case = testcase.new()
function case:Run()
    index = index + 1
    local result = self.blockchain:Invoke({
        func = "setHash",
        args = { tostring(self.index.tx), tostring(index) },
    })
    self.blockchain:Confirm(result)
    print(result.Ret[1])
    return result
end
return case



