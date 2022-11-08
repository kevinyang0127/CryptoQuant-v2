package script

import (
	lua "github.com/yuin/gopher-lua"
)

const (
	moduleName = "cryptoquant"
)

func loadmodule(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"saveData": saveData,
	"getData":  getData,
	"entry":    entry,
	"exit":     exit,
	"order":    order,
	"ma":       ma,
	"ema":      ema,
	"rsi":      rsi,
	"atr":      atr,
}

func saveData(L *lua.LState) int {
	return 0
}

func getData(L *lua.LState) int {
	return 0
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

func ma(L *lua.LState) int {
	return 0
}

func ema(L *lua.LState) int {
	return 0
}

func rsi(L *lua.LState) int {
	return 0
}

func atr(L *lua.LState) int {
	return 0
}
