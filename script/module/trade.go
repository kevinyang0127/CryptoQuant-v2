package script

import lua "github.com/yuin/gopher-lua"

func GetTradeExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry": Entry,
		"exit":  Exit,
		"order": Order,
	}
}

func Entry(L *lua.LState) int {
	return 0
}

func Exit(L *lua.LState) int {
	return 0
}

func Order(L *lua.LState) int {
	return 0
}
