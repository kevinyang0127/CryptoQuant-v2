package module

import (
	"CryptoQuant-v2/simulation"
	"context"
	"log"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

func GetBacktestExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry":           backtestEntry,
		"exit":            backtestExit,
		"exitAll":         backtestExitAll,
		"order":           backtestOrder,
		"cancelAllOrders": backtestCancelAllOrder,
		"getAllOrders":    backtestGetAllOrders,
		"hasPosition":     backtestHasPosition,
	}
}

/*
cryptoquant.entry(side, qty) --市價開倉
no return value
*/
func backtestEntry(L *lua.LState) int {
	paramCount := L.GetTop()
	if paramCount != 2 {
		log.Println("backtestEntry paramCount != 2")
		return 0
	}

	side := L.CheckBool(1)
	qty := L.CheckNumber(2)

	log.Printf("open, side = %v, qty = %f", side, qty)

	simulationID := L.GetGlobal("SimulationID").String()
	nowPrice := L.GetGlobal("NowPrice").String()
	klineEndTimeS := L.GetGlobal("KlineEndTime").String()
	klineEndTime, _ := strconv.ParseInt(klineEndTimeS, 10, 64)
	simulation.SimulationManager.Entry(context.Background(), simulationID, side, nowPrice, qty.String(), false, klineEndTime)

	return 0
}

/*
cryptoquant.exit(qty) --市價關倉
no return value
*/
func backtestExit(L *lua.LState) int {
	paramCount := L.GetTop()
	if paramCount < 1 {
		log.Println("backtestExit paramCount < 1")
		return 0
	}

	qty := L.CheckNumber(1)

	log.Printf("exit, qty = %f", qty)

	simulationID := L.GetGlobal("SimulationID").String()
	nowPrice := L.GetGlobal("NowPrice").String()
	klineEndTimeS := L.GetGlobal("KlineEndTime").String()
	klineEndTime, _ := strconv.ParseInt(klineEndTimeS, 10, 64)
	simulation.SimulationManager.Exit(context.Background(), simulationID, nowPrice, qty.String(), false, klineEndTime)

	return 0
}

/*
cryptoquant.exitAll() --市價全部平倉
no return value
*/
func backtestExitAll(L *lua.LState) int {

	log.Printf("exit all position")

	simulationID := L.GetGlobal("SimulationID").String()
	nowPrice := L.GetGlobal("NowPrice").String()
	klineEndTimeS := L.GetGlobal("KlineEndTime").String()
	klineEndTime, _ := strconv.ParseInt(klineEndTimeS, 10, 64)
	simulation.SimulationManager.ExitAll(context.Background(), simulationID, nowPrice, false, klineEndTime)

	return 0
}

/*
cryptoquant.order(side, price, qty) --限價下單
no return value
*/
func backtestOrder(L *lua.LState) int {
	paramCount := L.GetTop()
	if paramCount != 3 {
		log.Println("backtestOrder paramCount != 3")
		return 0
	}

	side := L.CheckBool(1)
	price := L.CheckNumber(2)
	qty := L.CheckNumber(3)

	log.Printf("get new order, side = %v, prict = %f, qty = %f", side, price, qty)

	simulationID := L.GetGlobal("SimulationID").String()
	simulation.SimulationManager.Order(context.Background(), simulationID, side, price.String(), qty.String())

	return 0
}

/*
cryptoquant.cancelAllOrder() --取消所有掛單
no return value
*/
func backtestCancelAllOrder(L *lua.LState) int {

	simulationID := L.GetGlobal("SimulationID").String()
	simulation.SimulationManager.CloseAllOrder(context.Background(), simulationID)

	log.Printf("cancel all order")

	return 0
}

/*
cryptoquant.getOrders() --取得目前所有掛單
return orders table

	orders{
		order{
			["side"] = true,
			["price"] = 1300.5,
			["qty"] = 0.5,
		},
		order{
			["side"] = false,
			["price"] = 1350.5,
			["qty"] = 0.5,
		}
	}
*/
func backtestGetAllOrders(L *lua.LState) int {

	simulationID := L.GetGlobal("SimulationID").String()
	orders, _ := simulation.SimulationManager.GetAllOrder(context.Background(), simulationID)

	// TODO: 改成回傳整個list
	L.Push(lua.LNumber(len(orders)))

	return 1
}

/*
cryptoquant.hasPosition() --目前是否還有倉位
return bool
*/
func backtestHasPosition(L *lua.LState) int {

	simulationID := L.GetGlobal("SimulationID").String()
	position, _ := simulation.SimulationManager.GetPosition(context.Background(), simulationID)
	hasPosition := position != nil

	L.Push(lua.LBool(hasPosition))

	return 1
}