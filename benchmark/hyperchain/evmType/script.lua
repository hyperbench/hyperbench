local switch = {
    [0] = function(self)
        -- uint8
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeUint8",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(ret.Status)
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3

        return ret
    end,
    [1] = function(self)
        -- uint16
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeUint16",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(ret.Status)
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [2] = function(self)
        -- uint32
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeUint32",
            Args = {p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [3] = function(self)
        -- uint64
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeUint64",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [4] = function(self)
        -- uint128
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeUint128",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [5] = function(self)
        -- uint256
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeUint256",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [6] = function(self)
        -- int8
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeInt8",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [7] = function(self)
        -- int16
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeInt16",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [8] = function(self)
        -- int32
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeInt32",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [9] = function(self)
        -- int64
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeInt64",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [10] = function(self)
        -- int128
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeInt128",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [11] = function(self)
        -- int256
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeInt256",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 1
        print(result[2][1], result[2][2]) -- 2, 2
        print(result[3][1], result[3][2], result[3][3]) -- 3, 3, 3
        return ret
    end,
    [12] = function(self)
        -- bytes1
        local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        local ret = self.blockchain:Invoke({
            Func = "typeBytes1",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1][1]) -- 49    is byte converted from "1"
        print(self.toolkit:String(result[1]))   -- 1
        print(self.toolkit:String(result[2][1]), self.toolkit:String(result[2][2])) -- 2, 2
        print(self.toolkit:String(result[3][1]), self.toolkit:String(result[3][2]), self.toolkit:String(result[3][3])) -- 3, 3, 3
        return ret
    end,
    [13] = function(self)
        -- bytes2
        local p1 = "11"
        local p2 = { "22", "22" }
        local p3 = { "33", "33", "33" }
        local ret = self.blockchain:Invoke({
            Func = "typeBytes2",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(self.toolkit:String(result[1]))   -- 11
        print(self.toolkit:String(result[2][1]), self.toolkit:String(result[2][2])) -- 22, 22
        print(self.toolkit:String(result[3][1]), self.toolkit:String(result[3][2]), self.toolkit:String(result[3][3])) -- 33, 33, 33
        return ret
    end,
    [14] = function(self)
        -- bytes7
        local p1 = "1111111"
        local p2 = { "2222222", "2222222" }
        local p3 = { "3333333", "3333333", "3333333" }
        local ret = self.blockchain:Invoke({
            Func = "typeBytes7",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(self.toolkit:String(result[1]))   -- 1111111
        print(self.toolkit:String(result[2][1]), self.toolkit:String(result[2][2])) -- 2222222, 2222222
        print(self.toolkit:String(result[3][1]), self.toolkit:String(result[3][2]), self.toolkit:String(result[3][3])) -- 3333333, 3333333, 3333333
        return ret
    end,
    [15] = function(self)
        -- bytes24
        local p1 = "1111111"
        local p2 = { "2222222", "2222222" }
        local p3 = { "3333333", "3333333", "3333333" }
        local ret = self.blockchain:Invoke({
            Func = "typeBytes24",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(self.toolkit:String(result[1]))   -- 1111111
        print(self.toolkit:String(result[2][1]), self.toolkit:String(result[2][2])) -- 2222222, 2222222
        print(self.toolkit:String(result[3][1]), self.toolkit:String(result[3][2]), self.toolkit:String(result[3][3])) -- 3333333, 3333333, 3333333
        return ret
    end,
    [16] = function(self)
        -- bytes32
        local p1 = "1111111"
        local p2 = { "2222222", "2222222" }
        local p3 = { "3333333", "3333333", "3333333" }
        local ret = self.blockchain:Invoke({
            Func = "typeBytes32",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(self.toolkit:String(result[1]))   -- 1111111
        print(self.toolkit:String(result[2][1]), self.toolkit:String(result[2][2])) -- 2222222,2222222
        print(self.toolkit:String(result[3][1]), self.toolkit:String(result[3][2]), self.toolkit:String(result[3][3])) -- 3333333, 3333333, 3333333
        return ret
    end,
    [17] = function(self)
        -- bool
        local p1 = "true"
        local p2 = { "true", "false" }
        local p3 = { "false", "true", "false" }
        local ret = self.blockchain:Invoke({
            Func = "typeBool",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- true
        print(result[2][1], result[2][2]) -- true, false
        print(result[3][1], result[3][2], result[3][3]) -- false, true, false
        return ret
    end,
    [18] = function(self)
        -- address
        local p1 = "123123"
        local p2 = { "123123", "123123" }
        local p3 = { "123123", "123123", "123123" }
        local ret = self.blockchain:Invoke({
            Func = "typeAddress",
            Args = { p1, p2, p3 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- 0x0000000000000000000000000000000000123123
        print(result[2][1], result[2][2]) -- 0x0000000000000000000000000000000000123123, 0x0000000000000000000000000000000000123123
        print(result[3][1], result[3][2], result[3][3]) -- 0x0000000000000000000000000000000000123123, 0x0000000000000000000000000000000000123123, 0x0000000000000000000000000000000000123123
        return ret
    end,
    [19] = function(self)
        -- string
        local p1 = "test string"
        local ret = self.blockchain:Invoke({
            Func = "typeString",
            Args = { p1 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(result[1])   -- test string
        return ret
    end,
    [20] = function(self)
        -- bytes
        local p1 = self.toolkit:Hex("test bytes") -- hex of "test bytes"
        local ret = self.blockchain:Invoke({
            Func = "typeBytes",
            Args = { p1 },
        })
        ret = self.blockchain:Confirm(ret)
        local result = ret.Ret
        print(self.toolkit:String(result[1]))   -- test bytes
        return ret
    end,

}

local case = testcase.new()

function case:Run()
    local num = (case.index.Tx) % 21
    print("current index:", case.index.Tx, "current switch:", num)
    return switch[num](case)
end

return case



