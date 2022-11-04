package indicator

import (
	"log"
	"time"

	binanceFutures "github.com/adshao/go-binance/v2/futures"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"github.com/shopspring/decimal"
)

type Kline struct {
	StartTime            int64
	EndTime              int64
	Open                 decimal.Decimal
	Close                decimal.Decimal
	High                 decimal.Decimal
	Low                  decimal.Decimal
	Volume               decimal.Decimal
	TradeNum             int64
	QuoteVolume          decimal.Decimal
	ActiveBuyVolume      decimal.Decimal
	ActiveBuyQuoteVolume decimal.Decimal
	IsFinal              bool
}

func BinanceFKlineToKline(k *binanceFutures.Kline) (*Kline, error) {
	openPrice, err := decimal.NewFromString(k.Open)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	closePrice, err := decimal.NewFromString(k.Close)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	highPrice, err := decimal.NewFromString(k.High)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	lowPrice, err := decimal.NewFromString(k.Low)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	volume, err := decimal.NewFromString(k.Volume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	quoteVolume, err := decimal.NewFromString(k.QuoteAssetVolume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	activeBuyVolume, err := decimal.NewFromString(k.TakerBuyBaseAssetVolume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	activeBuyQuoteVolume, err := decimal.NewFromString(k.TakerBuyQuoteAssetVolume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	return &Kline{
		StartTime:            k.OpenTime,
		EndTime:              k.CloseTime,
		Open:                 openPrice,
		Close:                closePrice,
		High:                 highPrice,
		Low:                  lowPrice,
		Volume:               volume,
		TradeNum:             k.TradeNum,
		QuoteVolume:          quoteVolume,
		ActiveBuyVolume:      activeBuyVolume,
		ActiveBuyQuoteVolume: activeBuyQuoteVolume,
		IsFinal:              true,
	}, nil
}

func BinanceFKlineEventToKline(event binanceFutures.WsKlineEvent) (*Kline, error) {
	openPrice, err := decimal.NewFromString(event.Kline.Open)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	closePrice, err := decimal.NewFromString(event.Kline.Close)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	highPrice, err := decimal.NewFromString(event.Kline.High)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	lowPrice, err := decimal.NewFromString(event.Kline.Low)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	volume, err := decimal.NewFromString(event.Kline.Volume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	quoteVolume, err := decimal.NewFromString(event.Kline.QuoteVolume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	activeBuyVolume, err := decimal.NewFromString(event.Kline.ActiveBuyVolume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	activeBuyQuoteVolume, err := decimal.NewFromString(event.Kline.ActiveBuyQuoteVolume)
	if err != nil {
		log.Println("decimal.NewFromString fail")
		return nil, err
	}

	return &Kline{
		StartTime:            event.Kline.StartTime,
		EndTime:              event.Kline.EndTime,
		Open:                 openPrice,
		Close:                closePrice,
		High:                 highPrice,
		Low:                  lowPrice,
		Volume:               volume,
		TradeNum:             event.Kline.TradeNum,
		QuoteVolume:          quoteVolume,
		ActiveBuyVolume:      activeBuyVolume,
		ActiveBuyQuoteVolume: activeBuyQuoteVolume,
		IsFinal:              event.Kline.IsFinal,
	}, nil
}

// 只加入已經收盤的k線
func AddKline(series *techan.TimeSeries, kline *Kline) {
	if !kline.IsFinal {
		return
	}

	// 累積到1000筆資料時 只保留最新100筆
	if len(series.Candles) >= 1000 {
		newSeries := techan.NewTimeSeries()
		for i := 900; i < len(series.Candles); i++ {
			newSeries.AddCandle(series.Candles[i])
		}
		series.Candles = newSeries.Candles
	}

	startTimeSec := kline.StartTime / 1000
	periodMs := kline.EndTime - kline.StartTime
	period := techan.NewTimePeriod(time.Unix(startTimeSec, 0), time.Duration(periodMs)*time.Millisecond)
	candle := techan.NewCandle(period)
	candle.OpenPrice = big.NewFromString(kline.Open.String())
	candle.ClosePrice = big.NewFromString(kline.Close.String())
	candle.MaxPrice = big.NewFromString(kline.High.String())
	candle.MinPrice = big.NewFromString(kline.Low.String())
	series.AddCandle(candle)
}
