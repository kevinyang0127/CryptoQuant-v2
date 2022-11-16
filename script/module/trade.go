package module

import lua "github.com/yuin/gopher-lua"

func GetTradeExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry":           entry,
		"exit":            exit,
		"order":           order,
		"cancelAllOrders": cancelAllOrder,
		"getAllOrders":    getAllOrders,
	}
}

func entry(L *lua.LState) int {
	return 0
}

func exit(L *lua.LState) int {
	return 0
}

func order(L *lua.LState) int {
	return 0
}

/*
cryptoquant.cancelAllOrder(side, price, qty) --取消所有掛單
no return value
*/
func cancelAllOrder(L *lua.LState) int {
	return 0
}

func getAllOrders(L *lua.LState) int {
	L.Push(lua.LNumber(0))
	return 1
}
