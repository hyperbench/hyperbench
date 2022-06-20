local case = testcase.new()
function case:Run()
    --print("run")
    local result = self.blockchain:Invoke({
        Func = "query",
        Args = { "A" }
    })
    --print(result)
    return result
end

return case

