package indicator

import (
	"context"

	"github.com/sdcoffey/techan"
	"github.com/shopspring/decimal"
)

type ExponentialMovingAverage struct {
	Timeframe int                // SMA計算長度
	series    *techan.TimeSeries // k線序列
}

// timeframe -> RSI計算長度
func NewExponentialMovingAverage(ctx context.Context, timeframe int) *ExponentialMovingAverage {
	return &ExponentialMovingAverage{
		Timeframe: timeframe,
		series:    techan.NewTimeSeries(),
	}
}

func (ema *ExponentialMovingAverage) AddKline(kline *Kline) {
	AddKline(ema.series, kline)
}

func (ema *ExponentialMovingAverage) Calculate() decimal.Decimal {
	closePrices := techan.NewClosePriceIndicator(ema.series)
	emaInd := techan.NewEMAIndicator(closePrices, ema.Timeframe)
	//計算最新的ema值
	v, _ := decimal.NewFromString(emaInd.Calculate(len(ema.series.Candles) - 1).String())
	return v
}
