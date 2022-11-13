package exchange

import (
	"CryptoQuant-v2/indicator"
	"context"
	"fmt"
)

const (
	// read-only api key & secret key
	roApiKey    = "6XTRkEGIUNXblIyauSpKh72Bm1IfRzj3xR5qy9SVONzMjpJtUZLYd2rcmnWYEiDf"
	roSecretKey = "OQX5ivZje1GiQtnDs573saHeqxVdWzJMxcaqxMQC2uMMIs2bE1GsmiKwQ8zJonaO"
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
	GetLimitKlineHistory(ctx context.Context, symbol string, timeframe string, limit int) ([]indicator.Kline, error)
	GetLimitKlineHistoryByTime(ctx context.Context, symbol string, timeframe string, limit int, startTimeMs int64, endTimeMs int64) ([]indicator.Kline, error)
}

func GetExchange(exchangeName string) (Exchange, error) {
	switch GetExchangeName(exchangeName) {
	case BINANCE_FUTURE:
		return newBinanceFuture(), nil
	default:
		return nil, fmt.Errorf("don't support exchange(%s)", exchangeName)
	}
}
