

local index = 0
local case = testcase.new()
function case:Run()
    index = index + 1
    local result
    if index % 2 == 1 then
        result = self.blockchain:Invoke({
            func="issue",
            args={"acc", tostring(index)}
        })
        print(result:result()[1])
    else
        result = self.blockchain:Invoke({
            func="getAccountBalance",
            args={"acc"},
        })
        print(result:result()[1])
    end
    return result
end


