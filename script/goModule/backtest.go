package script

import (
	"CryptoQuant-v2/simulation"
	"context"
	"log"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

func GetBacktestExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry": backtestEntry,
		"exit":  backtestExit,
		"order": backtestOrder,
	}
}

/*
cryptoquant.entry(side, qty) --市價開倉
no return value
*/
func backtestEntry(L *lua.LState) int {
	paramCount := L.GetTop()
	if paramCount != 2 {
		log.Println("BacktestEntry paramCount != 2")
		return 0
	}

	side := L.CheckBool(1)
	qty := L.CheckNumber(2)

	simulationID := L.GetGlobal("SimulationID").String()
	nowPrice := L.GetGlobal("NowPrice").String()
	klineEndTimeS := L.GetGlobal("KlineEndTime").String()
	klineEndTime, _ := strconv.ParseInt(klineEndTimeS, 10, 64)
	simulation.SimulationManager.Entry(context.Background(), simulationID, side, nowPrice, qty.String(), false, klineEndTime)

	return 0
}

/*
cryptoquant.exit(qty, closeAll) --市價關倉
no return value
*/
func backtestExit(L *lua.LState) int {
	return 0
}

/*
cryptoquant.order(side, price, qty) --限價下單
no return value
*/
func backtestOrder(L *lua.LState) int {
	return 0
}
