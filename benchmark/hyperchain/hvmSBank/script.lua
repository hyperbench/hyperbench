-- prepare account
local accountStrings = { "0", "1", "2", "3", "4", "5", "6", "7", "8", "9" }
local idx=1
local case = testcase.new()
function case:BeforeGet()
    -- prepare account balance
    for key, value in ipairs(accountStrings) do
        self.blockchain:Invoke({
            func = "com.hpc.sbank.invoke.IssueInvoke",
            args = { value, "1000000000" }
        })
    end
end

function case:Run()
    local result = self.blockchain:Invoke({
        func = "com.hpc.sbank.invoke.TransferInvoke",
        args = { accountStrings[idx % 10],
                 accountStrings[(idx + 1) % 10],
                 "100",
        }
    })
    --print(result)
    return result
end
return case