package module

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

func GetSaveDataExports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"saveData": saveData,
		"getData":  getData,
	}
}

func saveData(L *lua.LState) int {
	strategyID := L.GetGlobal("StrategyID").String()
	GetSaveStrategyData().save(strategyID, L.CheckAny(1))
	return 0
}

func getData(L *lua.LState) int {
	strategyID := L.GetGlobal("StrategyID").String()
	val := GetSaveStrategyData().get(strategyID)
	L.Push(val)
	return 1
}

var saveStrategyData *saveMetaData
var once sync.Once

type saveMetaData struct {
	mux         sync.Mutex
	metaDataMap map[string]lua.LValue
}

func GetSaveStrategyData() *saveMetaData {
	once.Do(func() {
		saveStrategyData = &saveMetaData{
			metaDataMap: make(map[string]lua.LValue),
		}
	})

	return saveStrategyData
}

func (s *saveMetaData) save(key string, metadata lua.LValue) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.metaDataMap[key] = metadata
}

func (s *saveMetaData) get(key string) (metadata lua.LValue) {
	s.mux.Lock()
	defer s.mux.Unlock()
	v, ok := s.metaDataMap[key]
	if !ok {
		return lua.LNil
	}
	return v
}
