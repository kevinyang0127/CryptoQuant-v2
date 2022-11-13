package script

import (
	lua "github.com/yuin/gopher-lua"
)

func GetDataExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"saveData": SaveData,
		"getData":  GetData,
	}
}

func SaveData(L *lua.LState) int {
	return 0
}

func GetData(L *lua.LState) int {
	return 0
}
