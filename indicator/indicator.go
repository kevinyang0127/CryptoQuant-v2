package indicator

import (
	"log"

	binanceFutures "github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
)

type Kline struct {
	StartTime            int64           `json:"startTime"`
	EndTime              int64           `json:"endTime"`
	Open                 decimal.Decimal `json:"open"`
	Close                decimal.Decimal `json:"close"`
	High                 decimal.Decimal `json:"high"`
	Low                  decimal.Decimal `json:"low"`
	Volume               decimal.Decimal `json:"volume"`
	TradeNum             int64           `json:"tradeNum"`
	QuoteVolume          decimal.Decimal `json:"quoteVolume"`
	ActiveBuyVolume      decimal.Decimal `json:"activeBuyVolume"`
	ActiveBuyQuoteVolume decimal.Decimal `json:"activeBuyQuoteVolume"`
	IsFinal              bool            `json:"isFinal"`
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
