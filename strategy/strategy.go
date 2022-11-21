package strategy

import "CryptoQuant-v2/market"

type Strategy interface {
	HandleKline(klines []market.Kline, kline market.Kline)
	HandleBackTestKline(simulationID string, klines []market.Kline, kline market.Kline)
}
