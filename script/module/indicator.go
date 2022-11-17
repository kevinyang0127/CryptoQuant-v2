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
		"closedArr":  closedArr,
		"openedArr":  openedArr,
		"highArr":    highArr,
		"lowArr":     lowArr,
		"sma":        sma,
		"ema":        ema,
		"rsi":        rsi,
		"atr":        atr,
		"macd1226":   macd1226,
		"supertrend": supertrend,
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
cryptoquant.sma(values, timeframe) --values is an array
return sma values
*/
func sma(L *lua.LState) int {
	if L.GetTop() != 2 {
		log.Println("rsi input param != 2")
		return 0
	}

	valuesTable := L.CheckTable(1)
	timeframe := L.CheckInt(2)
	values := []float64{}
	valuesTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := decimal.NewFromString(l2.String())
		if err != nil {
			log.Println("when parse valuesTable, decimal.NewFromString fail")
		}
		values = append(values, p.InexactFloat64())
	})

	sma := indicator.Sma(timeframe, values)
	t := &lua.LTable{}
	for i, v := range sma {
		t.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t)

	return 1
}

/*
cryptoquant.ema(values, timeframe) --values is an array
return ema value array
*/
func ema(L *lua.LState) int {
	if L.GetTop() != 2 {
		log.Println("rsi input param != 2")
		return 0
	}

	valuesTable := L.CheckTable(1)
	timeframe := L.CheckInt(2)
	values := []float64{}
	valuesTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := decimal.NewFromString(l2.String())
		if err != nil {
			log.Println("when parse valuesTable, decimal.NewFromString fail")
		}
		values = append(values, p.InexactFloat64())
	})

	ema := indicator.Ema(timeframe, values)
	t := &lua.LTable{}
	for i, v := range ema {
		t.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t)

	return 1
}

/*
cryptoquant.rsi(values, timeframe) --values is an array
return rsi value array
*/
func rsi(L *lua.LState) int {
	if L.GetTop() != 2 {
		log.Println("rsi input param != 2")
		return 0
	}

	valuesTable := L.CheckTable(1)
	timeframe := L.CheckInt(2)
	values := []float64{}
	valuesTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := decimal.NewFromString(l2.String())
		if err != nil {
			log.Println("when parse valuesTable, decimal.NewFromString fail")
		}
		values = append(values, p.InexactFloat64())
	})

	_, rsi := indicator.RsiPeriod(timeframe, values)
	t := &lua.LTable{}
	for i, v := range rsi {
		t.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t)

	return 1
}

/*
cryptoquant.atr(high, low, closing, timeframe) --closing is an array
return atr value array
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

	tr, _ := indicator.Atr(timeframe, high, low, closing)
	atr := indicator.Rma(timeframe, tr)
	t := &lua.LTable{}
	for i, v := range atr {
		t.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t)

	return 1
}

/*
cryptoquant.macd1226(values) --values is an array
return macd1226 value array and signal value array
*/
func macd1226(L *lua.LState) int {
	if L.GetTop() != 1 {
		log.Println("macd1226 input param != 1")
		return 0
	}

	valuesTable := L.CheckTable(1)
	values := []float64{}
	valuesTable.ForEach(func(l1, l2 lua.LValue) {
		p, err := decimal.NewFromString(l2.String())
		if err != nil {
			log.Println("when parse valuesTable, decimal.NewFromString fail")
		}
		values = append(values, p.InexactFloat64())
	})

	macd, signal := indicator.Macd(values)
	t1 := &lua.LTable{}
	for i, v := range macd {
		t1.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t1)

	t2 := &lua.LTable{}
	for i, v := range signal {
		t2.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t2)

	return 2
}

/*
cryptoquant.supertrend(high, low, closing, timeframe, factor)
return supertrend value array and direction array
*/
func supertrend(L *lua.LState) int {
	if L.GetTop() != 5 {
		log.Println("atr input param != 5")
		return 0
	}

	highTable := L.CheckTable(1)
	lowTable := L.CheckTable(2)
	closingTable := L.CheckTable(3)
	timeframe := L.CheckInt(4)
	factor := L.CheckInt(5)

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

	hl2 := []float64{}
	if len(high) != len(low) {
		log.Println("len(high) != len(low)")
	} else {
		for i, h := range high {
			hl2 = append(hl2, (h+low[i])/float64(2))
		}
	}

	tr, _ := indicator.Atr(timeframe, high, low, closing)
	atr := indicator.Rma(timeframe, tr)

	superTrend := []float64{}
	direction := []bool{}
	if len(hl2) != len(atr) || len(atr) != len(closing) {
		log.Println("len(hl2) != len(atr) || len(atr) != len(closing)")
	} else {
		for i, h := range hl2 {
			upperBand := h + float64(factor)*atr[i]
			lowerBand := h - float64(factor)*atr[i]

			if i == 0 {
				superTrend = append(superTrend, lowerBand)
				direction = append(direction, true)
				continue
			}

			// 前一根方向是多方
			if direction[i-1] {
				// 檢查目前收盤價是否低於前一個superTrend值
				if closing[i] < superTrend[i-1] {
					superTrend = append(superTrend, upperBand)
					direction = append(direction, false)
				} else {
					if lowerBand > superTrend[i-1] {
						superTrend = append(superTrend, lowerBand)
					} else {
						superTrend = append(superTrend, superTrend[i-1])
					}
					direction = append(direction, true)
				}
			} else { // 前一根方向是空方
				// 檢查目前收盤價是否高於前一個superTrend值
				if closing[i] > superTrend[i-1] {
					superTrend = append(superTrend, lowerBand)
					direction = append(direction, true)
				} else {
					if upperBand < superTrend[i-1] {
						superTrend = append(superTrend, upperBand)
					} else {
						superTrend = append(superTrend, superTrend[i-1])
					}
					direction = append(direction, false)
				}
			}
		}
	}

	t1 := &lua.LTable{}
	for i, v := range superTrend {
		t1.RawSetInt(i+1, lua.LNumber(v))
	}
	L.Push(t1)

	t2 := &lua.LTable{}
	for i, v := range direction {
		t2.RawSetInt(i+1, lua.LBool(v))
	}
	L.Push(t2)

	return 2
}
