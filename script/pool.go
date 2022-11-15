package script

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

/*
	Get一個vm要記得用完Put回去
*/

type luaStatePool struct {
	mux    sync.Mutex
	length int           //vm 總數
	vmList []*lua.LState //vm list
}

func (p *luaStatePool) get() *lua.LState {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.vmList) == 0 {
		return nil
	}

	v := p.vmList[0]
	p.vmList = p.vmList[1:]

	return v
}

func (p *luaStatePool) put(L *lua.LState) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.vmList = append(p.vmList, L)
}

func (p *luaStatePool) shutdown() {
	for _, L := range p.vmList {
		L.Close()
	}
}
