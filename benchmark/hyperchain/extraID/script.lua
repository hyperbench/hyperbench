local case = testcase.new()

function case:Run()
    self.blockchain:Option({ extraid = { 123, "abc", "efg" } })
    local result = case.blockchain:Transfer({
        from = "0", -- account0 alias
        to = "1", -- account1 alias
        amount = 0,
        extra = self.toolkit.RandStr(self.toolkit,1024),
    })
    -- print(result:result()[1]) -- 0x0
    return result
end

return case

