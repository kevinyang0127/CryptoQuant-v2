package script

import (
	"CryptoQuant-v2/indicator"
	"fmt"
	"log"
	"strconv"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type luaScriptHandler struct {
	precompileManager     *luaPrecompileManager
	statePool             *luaStatePool
	moduleManager         *moduleManager
	backtestStatePool     *luaStatePool
	backtestModuleManager *moduleManager
}

var handler *luaScriptHandler
var once sync.Once

func GetLuaScriptHandler() *luaScriptHandler {
	once.Do(func() {
		handler = &luaScriptHandler{
			precompileManager: newLuaPrecompileManager(),

			statePool: &luaStatePool{
				length: 10,
				vmList: []*lua.LState{},
			},
			moduleManager: newModuleManager(),

			backtestStatePool: &luaStatePool{
				length: 10,
				vmList: []*lua.LState{},
			},
			backtestModuleManager: newModuleManager(),
		}

		/*
			PreloadModule只是先註冊module被require時要呼叫的func，所以loader會在require("cryptoquant")時被執行
			所以manager.GetExports()會在每次require都被叫到一次
		*/
		loader := func(L *lua.LState) int {
			mod := L.SetFuncs(L.NewTable(), handler.moduleManager.getExports())
			L.Push(mod)
			return 1
		}

		// add lua vm into pool
		for i := 0; i < handler.statePool.length; i++ {
			L := lua.NewState()
			L.PreloadModule(moduleName, loader)
			handler.statePool.vmList = append(handler.statePool.vmList, L)
		}

		// backtest
		backtestLoader := func(L *lua.LState) int {
			mod := L.SetFuncs(L.NewTable(), handler.backtestModuleManager.getExports())
			L.Push(mod)
			return 1
		}

		// add lua vm into pool
		for i := 0; i < handler.backtestStatePool.length; i++ {
			L := lua.NewState()
			L.PreloadModule(moduleName, backtestLoader)
			handler.backtestStatePool.vmList = append(handler.backtestStatePool.vmList, L)
		}
	})
	return handler
}

func (h *luaScriptHandler) RunScriptHandleKline(script string) error {
	L := h.statePool.get()
	defer h.statePool.put(L)

	if err := h.precompileManager.doScript(L, script); err != nil {
		fmt.Println("L.DoString fail")
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
		log.Println(err)
		return err
	}

	return nil
}

func (h *luaScriptHandler) AddNewModule(funcName string, script string, nRet int) error {
	key, err := h.precompileManager.precompile(script)
	if err != nil {
		log.Println("luaPrecompileManager.Precompile fail")
		return err
	}
	h.moduleManager.addNewExport(funcName, h.genLuaScriptModuleFunc(funcName, key, nRet))
	return nil
}

/*
執行並呼叫lua寫的模組的callback func
*/
func (h *luaScriptHandler) genLuaScriptModuleFunc(funcName string, scriptHashKey string, nRet int) lua.LGFunction {
	return func(l *lua.LState) int {
		err := h.precompileManager.doScriptByKey(l, scriptHashKey)
		if err != nil {
			log.Println("luaPrecompileManager.DoScriptByKey fail")
			log.Println(err)
			return 0
		}

		args := []lua.LValue{}
		for i := 1; i <= l.GetTop(); i++ {
			args = append(args, l.CheckAny(i))
		}

		err = l.CallByParam(lua.P{
			Fn:      l.GetGlobal(funcName), // 呼叫funcName函數
			NRet:    nRet,                  // 指定返回值數量
			Protect: true,                  // 如果出現異常，是panic還是返回err
		}, args...) // 傳遞輸入參數
		if err != nil {
			log.Println("L.CallByParam fail")
			log.Println(err)
		}
		return 0
	}
}

func (h *luaScriptHandler) RunBacktestHandleKline(strategyID string, userID string, simulationID string, script string, kls []indicator.Kline, kl indicator.Kline) error {
	L := h.backtestStatePool.get()
	defer h.backtestStatePool.put(L)

	L.SetGlobal("SimulationID", lua.LString(simulationID))
	L.SetGlobal("UserID", lua.LString(userID))
	L.SetGlobal("StrategyID", lua.LString(strategyID))
	L.SetGlobal("NowPrice", lua.LString(kl.Close.StringFixed(6)))
	L.SetGlobal("KlineEndTime", lua.LString(strconv.FormatInt(kl.EndTime, 10)))

	if err := h.precompileManager.doScript(L, script); err != nil {
		log.Println("precompileManager.doScript fail")
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
