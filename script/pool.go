package script

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

/*
	Get一個vm要記得用完Put回去
*/

var realTradePool *luaStatePool //處理真實交易的虛擬機池
var backtestPool *luaStatePool  //處理回測的虛擬機池

func init() {
	realTradePool = &luaStatePool{
		length: 10,
		vmList: []*lua.LState{},
	}

	// add lua vm into pool
	for i := 0; i < realTradePool.length; i++ {
		L := lua.NewState()
		L.PreloadModule(moduleName, loadTradeModule) // loadTradeModule
		realTradePool.vmList = append(realTradePool.vmList, L)
	}

	backtestPool = &luaStatePool{
		length: 10,
		vmList: []*lua.LState{},
	}

	// add lua vm into pool
	for i := 0; i < realTradePool.length; i++ {
		L := lua.NewState()
		L.PreloadModule(moduleName, loadBacktestModule) // loadBacktestModule
		backtestPool.vmList = append(backtestPool.vmList, L)
	}
}

type luaStatePool struct {
	mux    sync.Mutex
	length int           //vm 總數
	vmList []*lua.LState //vm list
}

func (p *luaStatePool) Get() *lua.LState {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.vmList) == 0 {
		return nil
	}

	v := p.vmList[0]
	p.vmList = p.vmList[1:]

	return v
}

func (p *luaStatePool) Put(L *lua.LState) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.vmList = append(p.vmList, L)
}

func (p *luaStatePool) Shutdown() {
	for _, L := range p.vmList {
		L.Close()
	}
}
