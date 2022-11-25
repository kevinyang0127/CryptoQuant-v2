package strategy

import (
	"CryptoQuant-v2/market"
	"CryptoQuant-v2/script"
	"log"
)

type LuaScriptStrategy struct {
	luaScriptHandler *script.LuaScriptHandler
	strategyInfo     *StrategyInfo
}

func NewLuaScriptStrategy(luaScriptHandler *script.LuaScriptHandler, strategyInfo *StrategyInfo) *LuaScriptStrategy {
	return &LuaScriptStrategy{
		luaScriptHandler: luaScriptHandler,
		strategyInfo:     strategyInfo,
	}
}

func (s *LuaScriptStrategy) HandleKline(klines []market.Kline, kline market.Kline) error {
	err := s.luaScriptHandler.RunScriptHandleKline(s.strategyInfo.Exchange, s.strategyInfo.UserID, s.strategyInfo.Symbol, s.strategyInfo.Script, klines, kline)
	if err != nil {
		log.Println("luaScriptHandler.RunScriptHandleKline fail")
		return err
	}
	return nil
}

func (s *LuaScriptStrategy) HandleBackTestKline(simulationID string, klines []market.Kline, kline market.Kline) {
	s.luaScriptHandler.RunBacktestHandleKline(s.strategyInfo.StrategyID, s.strategyInfo.UserID, simulationID, s.strategyInfo.Script, klines, kline)
}
