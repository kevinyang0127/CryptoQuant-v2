package strategy

import "CryptoQuant-v2/indicator"

type Strategy interface {
	HandleKline(kline indicator.Kline)
}
