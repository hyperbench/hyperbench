local case = testcase.new()

function case:BeforeRun()
    case.blockchain:Option({
        account=case.index.tx
    })
end

function case:Run()
    local ret = case.blockchain:Transfer({
        from = "0",
        to = "1",
        amount = 0,
        extra = tostring(case.index.tx),
    })
    return ret
end

return case