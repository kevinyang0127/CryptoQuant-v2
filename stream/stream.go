package stream

import (
	"CryptoQuant-v2/market"
	"context"
)

type Stream interface {
	// 資料推送源
	wsConnect(ctx context.Context, klineHandler func(streamKey string, kline market.Kline)) error
}
