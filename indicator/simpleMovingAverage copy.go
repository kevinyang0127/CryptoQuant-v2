package indicator

import (
	"context"

	"github.com/sdcoffey/techan"
	"github.com/shopspring/decimal"
)

type SimpleMovingAverage struct {
	Timeframe int                // SMA計算長度
	series    *techan.TimeSeries // k線序列
}

// timeframe -> RSI計算長度
func NewSimpleMovingAverage(ctx context.Context, timeframe int) *SimpleMovingAverage {
	return &SimpleMovingAverage{
		Timeframe: timeframe,
		series:    techan.NewTimeSeries(),
	}
}

func (sma *SimpleMovingAverage) AddKline(kline *Kline) {
	AddKline(sma.series, kline)
}

func (sma *SimpleMovingAverage) Calculate() decimal.Decimal {
	closePrices := techan.NewClosePriceIndicator(sma.series)
	smaInd := techan.NewSimpleMovingAverage(closePrices, sma.Timeframe)
	//計算最新的sma值
	v, _ := decimal.NewFromString(smaInd.Calculate(len(sma.series.Candles) - 1).String())
	return v
}
