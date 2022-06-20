local case = testcase.new()

function case:BeforeRun()
    case.blockchain:Option({
        --hvm="method",
        account=case.index.Tx,
    })
end
function case:Run()
    local result = case.blockchain:Invoke({
            Func = "setHash",
            Args = { "key",
                    "value",
            }
        })
    return result
end
return case