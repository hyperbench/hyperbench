local case = testcase.new()
function case:Run()
    local result = self.blockchain:Invoke({
        func = "increase",
        args = {{"key","test"}},
    })
    --self.blockchain:Confirm(result)
    return result
end
return case