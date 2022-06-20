local case = testcase.new()

function case:Run()
    local ret = case.blockchain:Transfer({
        From = "0",
        To = "1",
        Amount = 0,
        Extra = tostring(case.index.Tx),
    })
    print(case.index.Worker, case.index.Tx, ret.Status)
    return ret
end

return case