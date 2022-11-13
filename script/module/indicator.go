package script

import lua "github.com/yuin/gopher-lua"

func GetIndicatorExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"ma":  Ma,
		"ema": Ema,
		"rsi": Rsi,
		"atr": Atr,
	}
}

func Ma(L *lua.LState) int {
	return 0
}

func Ema(L *lua.LState) int {
	return 0
}

func Rsi(L *lua.LState) int {
	return 0
}

func Atr(L *lua.LState) int {
	return 0
}
