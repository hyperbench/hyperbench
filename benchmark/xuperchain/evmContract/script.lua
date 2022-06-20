local case = testcase.new()
function case:Run()
    local result = self.blockchain:Invoke({
        Func = "increase",
        Args = {{"key","test"}},
    })
    --self.blockchain:Confirm(result)
    return result
end
return case