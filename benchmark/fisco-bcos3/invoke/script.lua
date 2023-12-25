local case = testcase.new()
function case:Run()
    local result = self.blockchain:Invoke({
        Func = "set",
        Args = {"blockchain"},
    })
    self.blockchain:Confirm(result)
    return result
end
return case