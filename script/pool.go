package script

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

/*
	Get一個vm要記得用完Put回去
*/

// Global LState pool
var pool *luaStatePool

func init() {
	pool = &luaStatePool{
		length: 10,
		vmList: []*lua.LState{},
	}

	// add lua vm into pool
	for i := 0; i < pool.length; i++ {
		L := lua.NewState()
		L.PreloadModule(moduleName, loadmodule)
		pool.vmList = append(pool.vmList, L)
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
