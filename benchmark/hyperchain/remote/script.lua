local case = testcase.new()

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