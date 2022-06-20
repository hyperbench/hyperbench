local index = 0
local case = testcase.new()
function case:Run()
    index = index + 1
    local result = self.blockchain:Invoke({
        Func = "setHash",
        Args = { tostring(self.index.Tx), tostring(index) },
    })
    self.blockchain:Confirm(result)
    print(result.Ret[1])
    return result
end
return case



