package strategy

import (
	"CryptoQuant-v2/indicator"
	"CryptoQuant-v2/script"
)

type LuaScriptStrategy struct {
	script string
	userID string
}

func NewLuaScriptStrategy(script string, userID string) *LuaScriptStrategy {
	return &LuaScriptStrategy{
		script: script,
		userID: userID,
	}
}

func (s *LuaScriptStrategy) HandleKline(klines []indicator.Kline, kline indicator.Kline) {
	script.RunScriptHandleKline(s.script, kline)
}

func (s *LuaScriptStrategy) HandleBackTestKline(simulationID string, klines []indicator.Kline, kline indicator.Kline) {
	script.RunBacktestHandleKline(s.script, s.userID, simulationID, klines, kline)
}
