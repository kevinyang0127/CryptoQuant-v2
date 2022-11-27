package module

import (
	"CryptoQuant-v2/simulation"
	"context"
	"log"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

func GetBacktestExports(simulationManager *simulation.Manager) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry":           getBacktestEntryLGFunc(simulationManager),
		"exit":            getBacktestExitLGFunc(simulationManager),
		"exitAll":         getBacktestExitAllLGFunc(simulationManager),
		"order":           getBacktestOrderLGFunc(simulationManager),
		"cancelAllOrders": getBacktestCancelAllOrderLGFunc(simulationManager),
		"getAllOrders":    getBacktestGetAllOrdersLGFunc(simulationManager),
		"hasPosition":     getBacktestHasPositionLGFunc(simulationManager),
		"lineNotif":       unsupport,
		"stopLossOrder":   getBacktestStopLossOrderLGFunc(simulationManager),
		"takeProfitOrder": getBacktestTakeProfitOrderLGFunc(simulationManager),
	}
}

/*
cryptoquant.entry(side, qty) --市價開倉
no return value
*/
func getBacktestEntryLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		paramCount := L.GetTop()
		if paramCount != 2 {
			log.Println("backtestEntry paramCount != 2")
			return 0
		}

		side := L.CheckBool(1)
		qty := L.CheckNumber(2)

		simulationID := L.GetGlobal("SimulationID").String()
		nowPrice := L.GetGlobal("NowPrice").String()
		klineEndTimeS := L.GetGlobal("KlineEndTime").String()
		klineEndTime, _ := strconv.ParseInt(klineEndTimeS, 10, 64)

		log.Printf("market entry, nowPrice = %s, side = %v qty = %f", nowPrice, side, qty)

		simulationManager.Entry(context.Background(), simulationID, side, nowPrice, qty.String(), false, klineEndTime)

		return 0
	}
	return fn
}

/*
cryptoquant.exit(qty) --市價關倉
no return value
*/
func getBacktestExitLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
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
		simulationManager.Exit(context.Background(), simulationID, nowPrice, qty.String(), false, klineEndTime)

		return 0
	}
	return fn
}

/*
cryptoquant.exitAll() --市價全部平倉
no return value
*/
func getBacktestExitAllLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {

		log.Printf("exit all position")

		simulationID := L.GetGlobal("SimulationID").String()
		nowPrice := L.GetGlobal("NowPrice").String()
		klineEndTimeS := L.GetGlobal("KlineEndTime").String()
		klineEndTime, _ := strconv.ParseInt(klineEndTimeS, 10, 64)
		simulationManager.ExitAll(context.Background(), simulationID, nowPrice, false, klineEndTime)

		return 0
	}
	return fn
}

/*
cryptoquant.order(side, price, qty) --限價下單
no return value
*/
func getBacktestOrderLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
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
		simulationManager.Order(context.Background(), simulationID, side, price.String(), qty.String())

		return 0
	}
	return fn
}

/*
cryptoquant.cancelAllOrder() --取消所有掛單
no return value
*/
func getBacktestCancelAllOrderLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {

		simulationID := L.GetGlobal("SimulationID").String()
		simulationManager.CloseAllOrder(context.Background(), simulationID)

		log.Printf("cancel all order")

		return 0
	}
	return fn
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
func getBacktestGetAllOrdersLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {

		simulationID := L.GetGlobal("SimulationID").String()
		orders, _ := simulationManager.GetAllOrder(context.Background(), simulationID)

		// TODO: 改成回傳整個list
		L.Push(lua.LNumber(len(orders)))

		return 1
	}
	return fn
}

/*
cryptoquant.hasPosition() --目前是否還有倉位
return bool
*/
func getBacktestHasPositionLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {

		simulationID := L.GetGlobal("SimulationID").String()
		position, _ := simulationManager.GetPosition(context.Background(), simulationID)
		hasPosition := position != nil

		L.Push(lua.LBool(hasPosition))

		return 1
	}
	return fn
}

/*
cryptoquant.stopLossOrder(side, price, qty, stopPrice) --限價停損單
no return value
當side為true時，當前價格小於stopPrice時觸發
當side為false時，當前價格大於stopPrice時觸發
*/
func getBacktestStopLossOrderLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {

		simulationID := L.GetGlobal("SimulationID").String()

		side := L.CheckBool(1)
		price := L.CheckNumber(2)
		qty := L.CheckNumber(3)
		stopPrice := L.CheckNumber(4)

		log.Printf("get new stop loss order, side = %v, price = %f, qty = %f, stopPirce = %f", side, price, qty, stopPrice)

		simulationManager.StopLossOrder(context.Background(), simulationID, side, price.String(), qty.String(), stopPrice.String())

		return 0
	}
	return fn
}

/*
cryptoquant.takeProfitOrder(side, price, qty, stopPrice) --限價停利單
no return value
當side為true時，當前價格大於stopPrice時觸發
當side為false時，當前價格小於stopPrice時觸發
*/
func getBacktestTakeProfitOrderLGFunc(simulationManager *simulation.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {

		simulationID := L.GetGlobal("SimulationID").String()

		side := L.CheckBool(1)
		price := L.CheckNumber(2)
		qty := L.CheckNumber(3)
		stopPrice := L.CheckNumber(4)

		log.Printf("get new take profit order, side = %v, price = %f, qty = %f, stopPirce = %f", side, price, qty, stopPrice)

		simulationManager.TakeProfitOrder(context.Background(), simulationID, side, price.String(), qty.String(), stopPrice.String())

		return 0
	}
	return fn
}
