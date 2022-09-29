local case = testcase.new()
function case:Run()
    local ret = case.blockchain:Transfer({
        From = "74d366e0649a91395bb122c005917644382b9452",
        To = "3b2b643246666bfa1332257c13d0d1283736838d",
        Amount = 100,
        Extra = "11",
    })
    self.blockchain:Confirm(ret)
    return ret
end

return case
