local case = testcase.new()

function case:Run()
    local ret = case.blockchain:Invoke({
        func="setHash",
        args={tostring(case.index.tx),
              tostring(case.index.worker)}
    })
    case.blockchain:Confirm(ret)
    return ret
end

return case