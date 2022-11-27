package exchange

import (
	"CryptoQuant-v2/market"
	"context"

	"github.com/shopspring/decimal"
)

type ExchangeName string

const (
	UNKNOWN        ExchangeName = "UNKNOWN"
	BINANCE_FUTURE ExchangeName = "BINANCE_FUTURE"
)

func GetExchangeName(exchangeName string) ExchangeName {
	switch exchangeName {
	case "BINANCE_FUTURE":
		return BINANCE_FUTURE
	default:
		return UNKNOWN
	}
}

type Exchange interface {
	GetLimitKlineHistory(ctx context.Context, symbol string, timeframe string, limit int) ([]market.Kline, error)
	GetLimitKlineHistoryByTime(ctx context.Context, symbol string, timeframe string, limit int, startTimeMs int64, endTimeMs int64) ([]market.Kline, error)

	// 目前倉位，沒有倉位-> return nil, nil
	GetPosition(ctx context.Context, symbol string) (*market.Position, error)

	// 創建市價訂單(taker)
	CreateMarketOrder(ctx context.Context, symbol string, side bool, quantity decimal.Decimal) error

	// 創建限價訂單(maker)
	CreateLimitOrder(ctx context.Context, symbol string, side bool, price decimal.Decimal, quantity decimal.Decimal) error

	// 取消所有掛單委託
	CancelAllOpenOrders(ctx context.Context, symbol string) error

	// 創建限價止損訂單
	CreateStopLossOrder(ctx context.Context, symbol string, side bool, price decimal.Decimal, quantity decimal.Decimal, stopPrice decimal.Decimal) error

	// 創建限價止盈訂單
	CreateTakeProfitOrder(ctx context.Context, symbol string, side bool, price decimal.Decimal, quantity decimal.Decimal, stopPrice decimal.Decimal) error
}
