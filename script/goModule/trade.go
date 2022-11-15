package gomodule

import lua "github.com/yuin/gopher-lua"

func GetTradeExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry": entry,
		"exit":  exit,
		"order": order,
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
