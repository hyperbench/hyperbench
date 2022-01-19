local case = testcase.new()
local accountNum=10
local idx=math.random( 0,accountNum )

function case:Run()
    local result = self.blockchain:Transfer({
        From = tostring(idx),
        To = tostring(math.random( 0,accountNum )),
        Amount = 1,
    })
    idx=(idx+1)%accountNum
    --self.blockchain:Confirm(result)
    return result
end
return case