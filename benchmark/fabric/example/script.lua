local invoke = { mode = "invoke" }
local query = { mode = "query" }
local case = testcase.new()
function case:Run()
    --print("run")
    local option
    if self.index.Tx % 2 == 0 then
        option = invoke
    else
        option = query
    end
    self.blockchain:Option(option)
    local result = self.blockchain:Invoke({
        Func = "query",
        Args = { "A" },
    })
    --print(result)
    return result
end
return case

