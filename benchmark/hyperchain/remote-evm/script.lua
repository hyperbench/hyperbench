local case = testcase.new()

function case:Run()
    local ret = case.blockchain:Invoke({
        func="setHash",
        args={tostring(case.index.Tx),
              tostring(case.index.Worker)}
    })
    case.blockchain:Confirm(ret)
    return ret
end

return case