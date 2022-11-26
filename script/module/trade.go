package module

import (
	"CryptoQuant-v2/exchange"
	"CryptoQuant-v2/notify"
	"context"
	"log"

	"github.com/shopspring/decimal"
	lua "github.com/yuin/gopher-lua"
)

func GetTradeExports(exchangeManager *exchange.Manager) map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"entry":           getEntryLGFunc(exchangeManager),
		"exit":            getExitLGFunc(exchangeManager),
		"exitAll":         getExitAllLGFunc(exchangeManager),
		"order":           getOrderLGFunc(exchangeManager),
		"cancelAllOrders": getCancelAllOrderLGFunc(exchangeManager),
		"getAllOrders":    unsupport,
		"hasPosition":     getHasPositionLGFunc(exchangeManager),
		"getBalance":      unsupport,
		"lineNotif":       getLineNotifLGFunc(),
	}
}

/*
cryptoquant.entry(side, qty) --市價開倉
no return value
*/
func getEntryLGFunc(exchangeManager *exchange.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		paramCount := L.GetTop()
		if paramCount != 2 {
			log.Println("cryptoquant.entry() paramCount != 2")
			return 0
		}

		exchangeName := L.GetGlobal("ExchangeName").String()
		userID := L.GetGlobal("UserID").String()
		symbol := L.GetGlobal("Symbol").String()

		side := L.CheckBool(1)
		qty := L.CheckNumber(2)

		ctx := context.Background()
		ex, err := exchangeManager.GetExchange(ctx, exchangeName, userID)
		if err != nil {
			log.Println("cryptoquant.entry() fail, exchangeManager.GetExchange error")
			log.Println(err)
			return 0
		}
		err = ex.CreateMarketOrder(ctx, symbol, side, decimal.NewFromFloat(float64(qty)))
		if err != nil {
			log.Println("cryptoquant.entry() fail, exchangeManager.CreateMarketOrder error")
			log.Println(err)
			return 0
		}

		return 0
	}
	return fn
}

/*
cryptoquant.exit(qty) --市價關倉
no return value
*/
func getExitLGFunc(exchangeManager *exchange.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		paramCount := L.GetTop()
		if paramCount != 1 {
			log.Println("cryptoquant.exit() paramCount != 1")
			return 0
		}

		exchangeName := L.GetGlobal("ExchangeName").String()
		userID := L.GetGlobal("UserID").String()
		symbol := L.GetGlobal("Symbol").String()

		qty := L.CheckNumber(1)
		qtyD := decimal.NewFromFloat(float64(qty))
		if qtyD.IsNegative() {
			log.Println("cryptoquant.exit() fail, input qty is negative")
			return 0
		}

		ctx := context.Background()
		ex, err := exchangeManager.GetExchange(ctx, exchangeName, userID)
		if err != nil {
			log.Println("cryptoquant.exit() fail, exchangeManager.GetExchange error")
			log.Println(err)
			return 0
		}
		p, err := ex.GetPosition(ctx, symbol)
		if err != nil {
			log.Println("cryptoquant.exit() fail, ex.GetPosition error")
			log.Println(err)
			return 0
		}

		side := p.Quantity.IsNegative() // side和持倉方向相反

		if p.Quantity.Abs().LessThan(qtyD) {
			qtyD = p.Quantity.Abs()
		}

		err = ex.CreateMarketOrder(ctx, symbol, side, qtyD)
		if err != nil {
			log.Println("cryptoquant.exit() fail, ex.CreateMarketOrder error")
			log.Println(err)
		}

		return 0
	}
	return fn
}

/*
cryptoquant.exitAll() --市價全部平倉
no return value
*/
func getExitAllLGFunc(exchangeManager *exchange.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		exchangeName := L.GetGlobal("ExchangeName").String()
		userID := L.GetGlobal("UserID").String()
		symbol := L.GetGlobal("Symbol").String()

		ctx := context.Background()
		ex, err := exchangeManager.GetExchange(ctx, exchangeName, userID)
		if err != nil {
			log.Println("cryptoquant.exitAll() fail, exchangeManager.GetExchange error")
			log.Println(err)
			return 0
		}
		p, err := ex.GetPosition(ctx, symbol)
		if err != nil {
			log.Println("cryptoquant.exitAll() fail, ex.GetPosition error")
			log.Println(err)
			return 0
		}

		side := p.Quantity.IsNegative() // side和持倉方向相反

		err = ex.CreateMarketOrder(ctx, symbol, side, p.Quantity.Abs())
		if err != nil {
			log.Println("cryptoquant.exitAll() fail, ex.CreateMarketOrder error")
			log.Println(err)
		}

		return 0
	}
	return fn
}

/*
cryptoquant.order(side, price, qty) --限價下單
no return value
*/
func getOrderLGFunc(exchangeManager *exchange.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		paramCount := L.GetTop()
		if paramCount != 3 {
			log.Println("cryptoquant.order() paramCount != 3")
			return 0
		}

		exchangeName := L.GetGlobal("ExchangeName").String()
		userID := L.GetGlobal("UserID").String()
		symbol := L.GetGlobal("Symbol").String()

		side := L.CheckBool(1)
		price := L.CheckNumber(2)
		priceD := decimal.NewFromFloat(float64(price))
		if priceD.IsNegative() {
			log.Println("cryptoquant.order() fail, input price is negative")
			return 0
		}
		qty := L.CheckNumber(3)
		qtyD := decimal.NewFromFloat(float64(qty))
		if qtyD.IsNegative() {
			log.Println("cryptoquant.order() fail, input qty is negative")
			return 0
		}

		ctx := context.Background()
		ex, err := exchangeManager.GetExchange(ctx, exchangeName, userID)
		if err != nil {
			log.Println("cryptoquant.order() fail, exchangeManager.GetExchange error")
			log.Println(err)
			return 0
		}

		err = ex.CreateLimitOrder(ctx, symbol, side, priceD, qtyD)
		if err != nil {
			log.Println("cryptoquant.order() fail, ex.CreateLimitOrder error")
			log.Println(err)
		}
		return 0
	}
	return fn
}

/*
cryptoquant.cancelAllOrder() --取消所有掛單
no return value
*/
func getCancelAllOrderLGFunc(exchangeManager *exchange.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		exchangeName := L.GetGlobal("ExchangeName").String()
		userID := L.GetGlobal("UserID").String()
		symbol := L.GetGlobal("Symbol").String()

		ctx := context.Background()
		ex, err := exchangeManager.GetExchange(ctx, exchangeName, userID)
		if err != nil {
			log.Println("cryptoquant.cancelAllOrder() fail, exchangeManager.GetExchange error")
			log.Println(err)
			return 0
		}

		err = ex.CancelAllOpenOrders(ctx, symbol)
		if err != nil {
			log.Println("cryptoquant.cancelAllOrder() fail, ex.CancelAllOpenOrders error")
			log.Println(err)
		}
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
func getAllOrdersLGFunc(exchangeManager *exchange.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		// exchangeName := L.GetGlobal("ExchangeName").String()
		// userID := L.GetGlobal("UserID").String()
		// symbol := L.GetGlobal("Symbol").String()

		// ctx := context.Background()
		// ex, err := exchangeManager.GetExchange(ctx, exchangeName, userID)
		// if err != nil {
		// 	log.Println("cryptoquant.cancelAllOrder() fail, exchangeManager.GetExchange error")
		// 	log.Println(err)
		// 	return 0
		// }

		// err = ex.(ctx, symbol)
		// if err != nil {
		// 	log.Println("cryptoquant.cancelAllOrder() fail, ex.CancelAllOpenOrders error")
		// 	log.Println(err)
		// }
		return 0
	}
	return fn
}

/*
cryptoquant.hasPosition() --目前是否還有倉位
return bool
*/
func getHasPositionLGFunc(exchangeManager *exchange.Manager) lua.LGFunction {
	fn := func(L *lua.LState) int {
		exchangeName := L.GetGlobal("ExchangeName").String()
		userID := L.GetGlobal("UserID").String()
		symbol := L.GetGlobal("Symbol").String()

		ctx := context.Background()
		ex, err := exchangeManager.GetExchange(ctx, exchangeName, userID)
		if err != nil {
			log.Println("cryptoquant.hasPosition() fail, exchangeManager.GetExchange error")
			log.Println(err)
			return 0
		}

		p, err := ex.GetPosition(ctx, symbol)
		if err != nil {
			log.Println("cryptoquant.hasPosition() fail, ex.GetPosition error")
			log.Println(err)
			return 0
		}
		hasPosition := p != nil

		L.Push(lua.LBool(hasPosition))
		return 1
	}
	return fn
}

/*
cryptoquant.lineNotif(msg) --發送line通知
no return value
*/
func getLineNotifLGFunc() lua.LGFunction {
	fn := func(L *lua.LState) int {
		paramCount := L.GetTop()
		if paramCount != 1 {
			log.Println("cryptoquant.lineNotif() paramCount != 1")
			return 0
		}

		msg := L.CheckString(1)
		err := notify.SendMsg(msg)
		if err != nil {
			log.Println("cryptoquant.lineNotif() fail,notify.SendMsg error")
		}

		return 0
	}
	return fn
}

func unsupport(L *lua.LState) int {
	return 0
}
