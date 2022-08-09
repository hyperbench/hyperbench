package glua

import (
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_blockchain(t *testing.T) {

	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	client, _ := NewMock()
	cLua := newBlockchain(L, client)
	L.SetField(mt, "blockchain", cLua)
	t.Run("DeployContract", func(t *testing.T) {
		scripts := []string{`
	    function run()
		    local i = 0
            case.blockchain.DeployContract()
            return i
		end
		`,
			`
	    function run()
		    local i = 1
            case.blockchain:DeployContract()
            return i
		end
		`,
		}
		for idx, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			assert.Nil(t, err)
			assert.Equal(t, ret, lua.LNumber(idx))
			assert.Equal(t, client.tempData[DeployContract], DeployContract)
		}

	})

	t.Run("Invoke", func(t *testing.T) {
		scripts := []string{`
		function run()
            ret = case.blockchain:Invoke({
		   func="123",
		   args={"123", "123"}
		},{aa="aa"},{bb="bb"})
            return ret
		end
			`,
			`
		function run()
            ret = case.blockchain.Invoke({
		   func="123",
		   args={"123", "123"}
		},{aa="aa"},{bb="bb"})
            return ret
		end
			`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			result := &fcom.Result{}
			err = TableLua2GoStruct(ret.(*lua.LTable), result)
			assert.Nil(t, err)
			assert.Equal(t, result, &fcom.Result{Label: "label", UID: "UUID", BuildTime: 0, SendTime: 0, ConfirmTime: 0, WriteTime: 0, Status: "success", Ret: []interface{}{"demo", "demo"}})
			assert.Equal(t, client.tempData[Invoke], fcom.Invoke{"123", []interface{}{"123", "123"}})
			assert.Equal(t, client.tempData[Option], []fcom.Option{{"aa": "aa"}, {"bb": "bb"}})
		}
	})

	t.Run("Transfer", func(t *testing.T) {
		scripts := []string{`
		function run()
            ret = case.blockchain:Transfer({
		   from="123",to="123",amount="1",extra="extra"
		},{aa="aa"},{bb="bb"})
            return ret
		end
			`,
			`
		function run()
            ret = case.blockchain.Transfer({
		   func="123",
		   args={"123", "123"}
		},{aa="aa"},{bb="bb"})
            return ret
		end
			`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			result := &fcom.Result{}
			err = TableLua2GoStruct(ret.(*lua.LTable), result)
			assert.Nil(t, err)
			assert.Equal(t, result, &fcom.Result{Label: "label", UID: "UUID", BuildTime: 0, SendTime: 0, ConfirmTime: 0, WriteTime: 0, Status: "success", Ret: []interface{}{"demo", "demo"}})
			assert.Equal(t, client.tempData[Transfer], fcom.Transfer{From: "123", To: "123", Amount: 1, Extra: "extra"})
		}
	})

	t.Run("Confirm", func(t *testing.T) {
		scripts := []string{`
		function run()
            ret = case.blockchain:Transfer({
		   from="123",to="123",amount="1",extra="extra"
		},{aa="aaa"},{bb="bbb"})
			ret = case.blockchain.Confirm(ret)
            return ret
		end
			`,
			`
		function run()
            ret = case.blockchain.Transfer({
		   from="123",to="123",amount="1",extra="extra"
		},{aa="aaa"},{bb="bbb"})
			ret = case.blockchain:Confirm(ret)
            return ret
		end	
		`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			result := &fcom.Result{}
			err = TableLua2GoStruct(ret.(*lua.LTable), result)
			assert.Nil(t, err)
			assert.Equal(t, result, &fcom.Result{Label: "Confirm", UID: "UUID", BuildTime: 0, SendTime: 0, ConfirmTime: 0, WriteTime: 0, Status: "confirm", Ret: []interface{}{"Confirm", "Confirm"}})
			assert.Equal(t, client.tempData[Transfer], fcom.Transfer{From: "123", To: "123", Amount: 1, Extra: "extra"})
			assert.Equal(t, client.tempData[Option], []fcom.Option(nil))
		}
	})

	t.Run("Query", func(t *testing.T) {
		scripts := []string{`
		function run()
            ret = case.blockchain:Query({
		    func="123",
			args={"banana","orange","apple"}}
			,{aa="aaa"},{bb="bbb"})
			return ret
		end
			`,
			`
			function run()
			ret = case.blockchain.Query({
			func="123",
			args={"banana","orange","apple"}}
			,{aa="aaa"},{bb="bbb"})
			return ret
			end
			`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			goRet, err := Lua2Go(ret)

			assert.Nil(t, err)
			assert.Equal(t, goRet, "nil")
			assert.Equal(t, client.tempData[Query], fcom.Query{Func: "123", Args: []interface{}{"banana", "orange", "apple"}})
			assert.Equal(t, client.tempData[Option], []fcom.Option{{"aa": "aaa"}, {"bb": "bbb"}})
		}
	})

	t.Run("Option", func(t *testing.T) {
		scripts := []string{`
		function run()
            ret = case.blockchain:Option(
				{aa="aaa",bb="bbb"})
			return ret
		end
			`,
			`
		function run()
            ret = case.blockchain.Option(
				{aa="aaa",bb="bbb"})
			return ret
		end
			`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			goRet, err := Lua2Go(ret)
			assert.Nil(t, err)
			assert.Equal(t, goRet, "")
			assert.Equal(t, client.tempData[Option], fcom.Option{"aa": "aaa", "bb": "bbb"})
		}
	})

	t.Run("GetContext", func(t *testing.T) {

		scripts := []string{`
		function run()
            ret = case.blockchain:GetContext()
			return ret
		end
			`,
			`
		function run()
            ret = case.blockchain.GetContext()
			return ret
		end
			`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			goRet, err := Lua2Go(ret)
			assert.Nil(t, err)
			assert.Equal(t, goRet, "GetContext")
		}

	})

	t.Run("SetContext", func(t *testing.T) {
		scripts := []string{`
		function run()
            ret = case.blockchain:SetContext("SetContext")
			return ret
		end
			`,
			`
		function run()
            ret = case.blockchain.SetContext("SetContext")
			return ret
		end
			`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			goRet, err := Lua2Go(ret)
			assert.Nil(t, err)
			assert.Equal(t, goRet, "")
		}
	})

	t.Run("RestContext", func(t *testing.T) {
		scripts := []string{`
		function run()
            ret = case.blockchain:ResetContext()
			return ret
		end
			`,
			`
		function run()
            ret = case.blockchain.ResetContext()
			return ret
		end
			`,
		}
		for _, script := range scripts {
			ret, err := runLuaRunFunc(L, script)
			goRet, err := Lua2Go(ret)
			assert.Nil(t, err)
			assert.Equal(t, goRet, "")
		}

	})
}
