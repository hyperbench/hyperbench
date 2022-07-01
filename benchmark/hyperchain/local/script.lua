local case = testcase.new()

function case:BeforeRun()
    case.blockchain:Option({
        account=case.index.Tx
    })
end

function case:Run()
    local ret = case.blockchain:Transfer({
        from = "0",
        to = "1",
        amount = 0,
        extra = tostring(case.index.Tx),
    })
    return ret
end

return case