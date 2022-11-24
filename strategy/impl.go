package strategy

import (
	"CryptoQuant-v2/market"
	"CryptoQuant-v2/script"
)

type LuaScriptStrategy struct {
	luaScriptHandler *script.LuaScriptHandler
	strategyID       string
	userID           string
	script           string
}

func NewLuaScriptStrategy(luaScriptHandler *script.LuaScriptHandler, strategyID string, userID string, script string) *LuaScriptStrategy {
	return &LuaScriptStrategy{
		luaScriptHandler: luaScriptHandler,
		strategyID:       strategyID,
		script:           script,
		userID:           userID,
	}
}

func (s *LuaScriptStrategy) HandleKline(klines []market.Kline, kline market.Kline) {
	s.luaScriptHandler.RunScriptHandleKline(s.script)
}

func (s *LuaScriptStrategy) HandleBackTestKline(simulationID string, klines []market.Kline, kline market.Kline) {
	s.luaScriptHandler.RunBacktestHandleKline(s.strategyID, s.userID, simulationID, s.script, klines, kline)
}
