package strategy

import (
	"CryptoQuant-v2/indicator"
	"CryptoQuant-v2/script"
)

type LuaScriptStrategy struct {
	strategyID string
	userID     string
	script     string
}

func NewLuaScriptStrategy(strategyID string, userID string, script string) *LuaScriptStrategy {
	return &LuaScriptStrategy{
		strategyID: strategyID,
		script:     script,
		userID:     userID,
	}
}

func (s *LuaScriptStrategy) HandleKline(klines []indicator.Kline, kline indicator.Kline) {
	script.GetLuaScriptHandler().RunScriptHandleKline(s.script)
}

func (s *LuaScriptStrategy) HandleBackTestKline(simulationID string, klines []indicator.Kline, kline indicator.Kline) {
	script.GetLuaScriptHandler().RunBacktestHandleKline(s.strategyID, s.userID, simulationID, s.script, klines, kline)
}
