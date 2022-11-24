package script

import "testing"

func TestXxx(t *testing.T) {
	// script := `
	// function HandleKline(klines, kline)
	// 	local cryptoquant = require("cryptoquant")
	// 	local timeframe = 4
	// 	local indicator = {1232.92,1220.52,1215.60,1213.15,1216.27,1216.73,1215.31,1218.56,1222.40,1225.86,1227.03,1225.11,1222.76,1216.07,1219.16,1232.27}
	// 	local rsi = cryptoquant.rsi(indicator, timeframe)
	// 	print(rsi[#rsi-1])
	// end
	// `
	// GetLuaScriptHandler().RunScriptHandleKline(script)
	t.Error("errrrr")
}

func TestXxx2(t *testing.T) {
	// script := `
	// function HandleKline(klines, kline)
	// 	local cryptoquant = require("cryptoquant")
	// 	local closedArr = cryptoquant.closedArr(klines)
	// 	local high = cryptoquant.highArr(klines)
	// 	local low = cryptoquant.lowArr(klines)
	// 	local supertrend, dic = cryptoquant.supertrend(high, low, closedArr, 14, 3)
	// 	print(supertrend[#supertrend])
	// 	print(dic[#dic])
	// end
	// `
	// GetLuaScriptHandler().RunScriptHandleKline(script)
	t.Error("errrrr")
}
