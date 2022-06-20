local case = testcase.new()
function case:Run()
    local result = self.blockchain:Invoke({
        Func = "test",
        Args = {"foo","bar"},
    })
    self.blockchain:Confirm(result)
    return result
end
return case