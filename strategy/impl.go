package strategy

import (
	"CryptoQuant-v2/indicator"
	"CryptoQuant-v2/script"
)

type LuaScriptStrategy struct {
	script string
}

func NewLuaScriptStrategy(script string) *LuaScriptStrategy {
	return &LuaScriptStrategy{
		script: script,
	}
}

func (s *LuaScriptStrategy) HandleKline(kline indicator.Kline) {
	script.RunScriptHandleKline(s.script, kline)
}
