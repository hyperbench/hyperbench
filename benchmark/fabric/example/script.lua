local index = 0
local invoke = { mode = "invoke" }
local query = { mode = "query" }
local case = testcase.new()
function case:Run()
    --print("run")
    index = index + 1
    local option
    if self.index.tx % 2 == 0 then
        option = invoke
    else
        option = query
    end
    self.blockchain:Option(option)
    local result = self.blockchain:Invoke({
        func = "query",
        args = { "A" },
    })
    --print(result)
    return result
end
return case

