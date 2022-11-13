package script

import (
	"CryptoQuant-v2/indicator"
	"log"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

func RunScriptHandleKline(script string, value indicator.Kline) error {
	L := realTradePool.Get()
	defer realTradePool.Put(L)

	if err := L.DoString(script); err != nil {
		log.Println("L.DoString fail")
		return err
	}

	klines := &lua.LTable{}
	kline := &lua.LTable{}

	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("HandleKline"), // 呼叫HandleKline函數
		NRet:    0,                          // 指定返回值數量
		Protect: true,                       // 如果出現異常，是panic還是返回err
	}, klines, kline) // 傳遞輸入參數
	if err != nil {
		log.Println("L.CallByParam fail")
		return err
	}

	return nil
}

func RunBacktestHandleKline(script string, userID string, simulationID string, kls []indicator.Kline, kl indicator.Kline) error {
	L := backtestPool.Get()
	defer backtestPool.Put(L)

	L.SetGlobal("SimulationID", lua.LString(simulationID))
	L.SetGlobal("UserID", lua.LString(userID))
	L.SetGlobal("NowPrice", lua.LString(kl.Close.StringFixed(6)))
	L.SetGlobal("KlineEndTime", lua.LString(strconv.FormatInt(kl.EndTime, 10)))

	if err := L.DoString(script); err != nil {
		log.Println("L.DoString fail")
		return err
	}

	klines := &lua.LTable{}
	kline := &lua.LTable{}

	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("HandleKline"), // 呼叫HandleKline函數
		NRet:    0,                          // 指定返回值數量
		Protect: true,                       // 如果出現異常，是panic還是返回err
	}, klines, kline) // 傳遞輸入參數
	if err != nil {
		log.Println("L.CallByParam fail")
		return err
	}

	return nil
}
