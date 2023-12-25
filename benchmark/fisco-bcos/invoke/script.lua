local case = testcase.new()
function case:Run()
    local result = self.blockchain:Invoke({
        Func = "setItem",
        Args = {"foo","bar"},
    })
    self.blockchain:Confirm(result)
    return result
end
return case