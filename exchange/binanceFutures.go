package exchange

import (
	"CryptoQuant-v2/market"
	"context"
	"log"

	binanceFutures "github.com/adshao/go-binance/v2/futures"
)

type BinanceFuture struct {
}

func newBinanceFuture() *BinanceFuture {
	return &BinanceFuture{}
}

func (bf *BinanceFuture) GetLimitKlineHistory(ctx context.Context, symbol string, timeframe string, limit int) ([]market.Kline, error) {
	c := binanceFutures.NewClient(roApiKey, roSecretKey)
	klinesService := c.NewKlinesService()
	klinesService.Symbol(symbol)
	klinesService.Interval(timeframe)
	klinesService.Limit(limit)
	res, err := klinesService.Do(ctx)
	if err != nil {
		log.Println("klinesService.Do fail")
		return nil, err
	}

	klines := []market.Kline{}
	for i, k := range res {
		if i == len(res)-1 {
			//檔下尚未收盤的k線也會拿到，所以最後一根k線不做事，因為高機率初始化時不會剛好收盤
			break
		}

		kline, err := market.BinanceFKlineToKline(k)
		if err != nil {
			log.Println("market.BinanceFKlineToKline fail")
			return nil, err
		}

		klines = append(klines, *kline)
	}

	return klines, nil
}

func (bf *BinanceFuture) GetLimitKlineHistoryByTime(ctx context.Context, symbol string, timeframe string, limit int, startTimeMs int64, endTimeMs int64) ([]market.Kline, error) {
	c := binanceFutures.NewClient(roApiKey, roSecretKey)
	klinesService := c.NewKlinesService()
	klinesService.Symbol(symbol)
	klinesService.Interval(timeframe)
	klinesService.Limit(limit)
	klinesService.StartTime(startTimeMs)
	klinesService.EndTime(endTimeMs)
	res, err := klinesService.Do(ctx)
	if err != nil {
		log.Println("klinesService.Do fail")
		return nil, err
	}

	klines := []market.Kline{}
	for i, k := range res {
		if i == len(res)-1 {
			//檔下尚未收盤的k線也會拿到，所以最後一根k線不做事，因為高機率初始化時不會剛好收盤
			break
		}

		kline, err := market.BinanceFKlineToKline(k)
		if err != nil {
			log.Println("market.BinanceFKlineToKline fail")
			return nil, err
		}

		klines = append(klines, *kline)
	}

	return klines, nil
}
