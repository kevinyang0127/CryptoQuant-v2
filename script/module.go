package script

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

const (
	moduleName = "cryptoquant"
)

type moduleManager struct {
	mux     sync.Mutex
	exports map[string]lua.LGFunction
}

func newModuleManager() *moduleManager {
	return &moduleManager{
		exports: make(map[string]lua.LGFunction),
	}
}

func (m *moduleManager) getExports() map[string]lua.LGFunction {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.exports
}

func (m *moduleManager) addNewExport(funcName string, lgFunc lua.LGFunction) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.exports[funcName] = lgFunc
}

func (m *moduleManager) deleteExport(funcName string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.exports, funcName)
}
