local case = testcase.new()
function case:Run()
    --print("run")

    local aBool = "true"
    local aChar = "c"
    local aByte = "20"
    local aShort = "100"
    local anInt = "1000"
    local aLong = "10000"
    local aFloat = "1.1"
    local aDouble = "1.11"
    local aString = "string"
    local person1 = { "tom", "21" }
    local person2 = { "jack", "18" }
    local bean1 = { "hvm-beam1", person1 }
    local bean2 = { "hvm-bean2", person2 }
    local strList = { "strList1", "strList2" }
    local personList = { person1, person2 }
    local personMap = { { "person1", person1 }, { "person2", person2 } }
    local beanMap = { { "bean1", bean1 }, { "bean2", bean2 } }
    local result = self.blockchain:Invoke({
        func="cn.hyperchain.contract.invoke.EasyInvoke",
        args={
            aBool,
            aChar,
            aByte,
            aShort,
            anInt,
            aLong,
            aFloat,
            aDouble,
            aString,
            person1,
            bean1,
            strList,
            personList,
            personMap,
            beanMap
        },
    })
    --self.blockchain:Confirm(result)
    --print(result.Ret[1])
    return result
end
return case