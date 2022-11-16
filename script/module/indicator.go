package module

import (
	"log"
	"strconv"

	"github.com/cinar/indicator"
	"github.com/shopspring/decimal"
	lua "github.com/yuin/gopher-lua"
)

func GetIndicatorExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"closedArr": closedArr,
		"openedArr": openedArr,
		"highArr":   highArr,
		"lowArr":    lowArr,
		// "ma":  rsi,
		// "ema": rsi,
		"rsi": rsi,
		"atr": atr,
	}
}

/*
cryptoquant.closedArr(klines)
return closed price array
*/
func closedArr(L *lua.LState) int {
	if L.GetTop() != 1 {
		log.Println("closedArr input param != 1")
		return 0
	}

	klinesTable := L.CheckTable(1)
	closing := &lua.LTable{}
	klinesTable.ForEach(func(l1, l2 lua.LValue) {
		kline, ok := l2.(*lua.LTable)
		if ok {
			v := kline.RawGet(lua.LString("close"))
			closing.RawSet(l1, v)
		}
	})

	L.Push(closing)
	return 1
}

/*
cryptoquant.closedArr(klines)
return closed price array
*/
func openedArr(L *lua.LState) int {
	if L.GetTop() != 1 {
		log.Println("openedArr input param != 1")
		return 0
	}

	klinesTable := L.CheckTable(1)
	open := &lua.LTable{}
	klinesTable.ForEach(func(l1, l2 lua.LValue) {
		kline, ok := l2.(*lua.LTable)
		if ok {
			v := kline.RawGet(lua.LString("open"))
			open.RawSet(l1, v)
		}
	})

	L.Push(open)
	return 1
}

/*
cryptoquant.closedArr(klines)
return closed price array
*/
func highArr(L *lua.LState) int {
	if L.GetTop() != 1 {
		log.Println("highArr input param != 1")
		return 0
	}

	klinesTable := L.CheckTable(1)
	high := &lua.LTable{}
	klinesTable.ForEach(func(l1, l2 lua.LValue) {
		kline, ok := l2.(*lua.LTable)
		if ok {
			v := kline.RawGet(lua.LString("high"))
			high.RawSet(l1, v)
		}
	})

	L.Push(high)
	return 1
}

/*
cryptoquant.closedArr(klines)
return closed price array
*/
func lowArr(L *lua.LState) int {
	if L.GetTop() != 1 {
		log.Println("lowArr input param != 1")
		return 0
	}

	klinesTable := L.CheckTable(1)
	low := &lua.LTable{}
	klinesTable.ForEach(func(l1, l2 lua.LValue) {
		kline, ok := l2.(*lua.LTable)
		if ok {
			v := kline.RawGet(lua.LString("low"))
			low.RawSet(l1, v)
		}
	})

	L.Push(low)
	return 1
}

/*
cryptoquant.rsi(closing, timeframe) --closing is an array
return lastest rsi value array
*/
func rsi(L *lua.LState) int {
	if L.GetTop() != 2 {
		log.Println("rsi input param != 2")
		return 0
	}

	closingTable := L.CheckTable(1)
	timeframe := L.CheckInt(2)
	closing := []float64{}
	closingTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := decimal.NewFromString(l2.String())
		if err != nil {
			log.Println("when parse closingTable, decimal.NewFromString fail")
		}
		closing = append(closing, p.InexactFloat64())
	})

	_, rsi := indicator.RsiPeriod(timeframe, closing)
	t := &lua.LTable{}
	for i, v := range rsi {
		t.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t)

	return 1
}

/*
cryptoquant.atr(high, low, closing, timeframe) --closing is an array
return lastest atr value array
*/
func atr(L *lua.LState) int {
	if L.GetTop() != 4 {
		log.Println("atr input param != 4")
		return 0
	}

	highTable := L.CheckTable(1)
	lowTable := L.CheckTable(2)
	closingTable := L.CheckTable(3)
	timeframe := L.CheckInt(4)

	high := []float64{}
	highTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := strconv.ParseFloat(l2.String(), 64)
		if err != nil {
			log.Println("when parse highTable, strconv.ParseFloat fail")
		}
		high = append(high, p)
	})

	low := []float64{}
	lowTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := strconv.ParseFloat(l2.String(), 64)
		if err != nil {
			log.Println("when parse lowTable, strconv.ParseFloat fail")
		}
		low = append(low, p)
	})

	closing := []float64{}
	closingTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := strconv.ParseFloat(l2.String(), 64)
		if err != nil {
			log.Println("when parse closingTable, strconv.ParseFloat fail")
		}
		closing = append(closing, p)
	})

	_, atr := indicator.Atr(timeframe, high, low, closing)
	t := &lua.LTable{}
	for i, v := range atr {
		t.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t)

	return 1
}
