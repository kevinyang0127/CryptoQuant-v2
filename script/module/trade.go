package module

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func GetTradeExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry":           entry,
		"exit":            exit,
		"exitAll":         exitAll,
		"order":           order,
		"cancelAllOrders": cancelAllOrder,
		"getAllOrders":    getAllOrders,
		"hasPosition":     hasPosition,
		"getBalance":      unsupport,
	}
}

/*
cryptoquant.entry(side, qty) --市價開倉
no return value
*/
func entry(L *lua.LState) int {
	paramCount := L.GetTop()
	if paramCount != 2 {
		log.Println("backtestEntry paramCount != 2")
		return 0
	}

	// side := L.CheckBool(1)
	// qty := L.CheckNumber(2)

	return 0
}

func exit(L *lua.LState) int {
	return 0
}

func exitAll(L *lua.LState) int {
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

/*
cryptoquant.hasPosition() --目前是否還有倉位
return bool
*/
func hasPosition(L *lua.LState) int {

	return 0
}

func unsupport(L *lua.LState) int {
	return 0
}
