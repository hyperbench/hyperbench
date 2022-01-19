local case = testcase.new()

function case:Run()
    local ret = case.blockchain:Transfer({
        from = "0",
        to = "1",
        amount = 0,
        extra = tostring(case.index.tx),
    })
    print(case.index.worker, case.index.tx, ret.status)
    return ret
end

return case