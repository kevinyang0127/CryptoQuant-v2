package script

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"sync"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

type luaPrecompileManager struct {
	mux      sync.Mutex
	protoMap map[string]*lua.FunctionProto // key: sha256(sourceCode)
}

func newLuaPrecompileManager() *luaPrecompileManager {
	return &luaPrecompileManager{
		protoMap: make(map[string]*lua.FunctionProto),
	}
}

func (m *luaPrecompileManager) precompile(sourceCode string) (hashKey string, err error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	proto, err := m.compileLuaString(sourceCode)
	if err != nil {
		log.Println("compileLuaString fail")
		return "", err
	}
	key := m.genKey(sourceCode)
	m.protoMap[key] = proto
	return key, nil
}

func (m *luaPrecompileManager) doScriptByKey(L *lua.LState, hashKey string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	p, ok := m.protoMap[hashKey]
	if !ok {
		return fmt.Errorf("DoScriptByKey fail, can't find this key")
	}
	return m.doCompiledFile(L, p)
}

func (m *luaPrecompileManager) doScript(L *lua.LState, sourceCode string) error {
	key := m.genKey(sourceCode)
	p, ok := m.protoMap[key]
	if !ok {
		k, err := m.precompile(sourceCode)
		if err != nil {
			log.Println("Precompile fail")
			return err
		}
		return m.doScriptByKey(L, k)
	}
	return m.doCompiledFile(L, p)
}

func (m *luaPrecompileManager) genKey(sourceCode string) string {
	bkey := sha256.New().Sum([]byte(sourceCode))
	return fmt.Sprintf("%x", bkey)
}

// CompileLua reads the passed lua file from disk and compiles it.
func (m *luaPrecompileManager) compileLuaString(source string) (*lua.FunctionProto, error) {
	reader := strings.NewReader(source)
	chunk, err := parse.Parse(reader, source)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, source)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func (m *luaPrecompileManager) doCompiledFile(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}
