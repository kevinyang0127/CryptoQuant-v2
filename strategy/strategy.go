package strategy

import "CryptoQuant-v2/indicator"

type Strategy interface {
	HandleKline(klines []indicator.Kline, kline indicator.Kline)
	HandleBackTestKline(simulationID string, klines []indicator.Kline, kline indicator.Kline)
}
