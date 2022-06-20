local case = testcase.new()

function case:Run()
    self.blockchain:Option({ extraid = { 123, "abc", "efg" } })
    local result = case.blockchain:Transfer({
        From = "0", -- account0 alias
        To = "1", -- account1 alias
        Amount = 0,
        Extra = self.toolkit:RandStr(1024),
    })
    -- print(result:result()[1]) -- 0x0
    return result
end

return case

