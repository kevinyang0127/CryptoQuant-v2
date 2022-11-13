package script

import (
	module "CryptoQuant-v2/script/module"

	lua "github.com/yuin/gopher-lua"
)

const (
	moduleName = "cryptoquant"
)

func loadTradeModule(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{}

	for k, v := range module.GetTradeExports() {
		exports[k] = v
	}

	for k, v := range module.GetIndicatorExports() {
		exports[k] = v
	}

	for k, v := range module.GetDataExports() {
		exports[k] = v
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

func loadBacktestModule(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{}

	for k, v := range module.GetBacktestExports() {
		exports[k] = v
	}

	for k, v := range module.GetIndicatorExports() {
		exports[k] = v
	}

	for k, v := range module.GetDataExports() {
		exports[k] = v
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}
