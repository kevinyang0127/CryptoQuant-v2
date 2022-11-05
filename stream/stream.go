package stream

import (
	"CryptoQuant-v2/indicator"
	"context"
)

type Stream interface {
	// 資料推送源
	wsConnect(ctx context.Context, klineHandler func(streamKey string, kline indicator.Kline)) error
}
