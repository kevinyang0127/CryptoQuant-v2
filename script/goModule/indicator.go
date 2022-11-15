package gomodule

import (
	"log"
	"strconv"

	"github.com/cinar/indicator"
	lua "github.com/yuin/gopher-lua"
)

func GetIndicatorExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		// "ma":  rsi,
		// "ema": rsi,
		"rsi": rsi,
		// "atr": rsi,
	}
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
		p, err := strconv.ParseFloat(l2.String(), 64)
		if err != nil {
			closing = []float64{}
			log.Println("when parse closingTable, strconv.ParseFloat fail")
			return
		}
		closing = append(closing, p)
	})

	_, rsi := indicator.RsiPeriod(timeframe, closing)
	t := &lua.LTable{}
	for i, v := range rsi {
		t.Insert(i+1, lua.LNumber(v))
	}
	L.Push(t)

	return 1
}
