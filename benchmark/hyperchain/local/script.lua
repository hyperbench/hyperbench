local case = testcase.new()

function case:BeforeRun()
    case.blockchain:Option({
        account=case.index.Tx
    })
end

function case:Run()
    local ret = case.blockchain:Transfer({
        From = "0",
        To = "1",
        Amount = 0,
        Extra = tostring(case.index.Tx),
    })
    return ret
end

return case