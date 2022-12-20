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
	err := s.luaScriptHandler.RunScriptHandleKline(s.strategyInfo.StrategyID, s.strategyInfo.UserID,
		s.strategyInfo.Exchange, s.strategyInfo.Symbol, s.strategyInfo.Timeframe, s.strategyInfo.Script,
		klines, kline)
	if err != nil {
		log.Println("luaScriptHandler.RunScriptHandleKline fail")
		return err
	}
	return nil
}

func (s *LuaScriptStrategy) HandleBackTestKline(simulationID string, klines []market.Kline, kline market.Kline) {
	s.luaScriptHandler.RunBacktestHandleKline(s.strategyInfo.StrategyID, s.strategyInfo.UserID, simulationID, s.strategyInfo.Script, klines, kline)
}

func (s *LuaScriptStrategy) UpdateStrategyInfo(newInfo *StrategyInfo) error {
	err := s.luaScriptHandler.CleanScriptPrecomplieCache(s.strategyInfo.Script)
	if err != nil {
		log.Println("luaScriptHandler.CleanScriptPrecomplieCache fail")
		return err
	}

	s.strategyInfo = newInfo
	return nil
}
