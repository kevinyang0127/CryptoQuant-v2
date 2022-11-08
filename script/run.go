package script

import (
	"CryptoQuant-v2/indicator"
	"log"

	lua "github.com/yuin/gopher-lua"
)

func RunScriptHandleKline(script string, value indicator.Kline) error {
	L := pool.Get()
	defer pool.Put(L)

	if err := L.DoString(script); err != nil {
		log.Println("L.DoString fail")
		return err
	}

	klines := &lua.LTable{}
	kline := &lua.LTable{}

	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("handleKline"), // 呼叫handleKline函數
		NRet:    0,                          // 指定返回值數量
		Protect: true,                       // 如果出現異常，是panic還是返回err
	}, klines, kline) // 傳遞輸入參數
	if err != nil {
		log.Println("L.CallByParam fail")
		return err
	}

	return nil
}
