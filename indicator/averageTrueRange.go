package indicator

import (
	"context"

	"github.com/sdcoffey/techan"
	"github.com/shopspring/decimal"
)

type AverageTrueRange struct {
	Timeframe int                // atr計算長度
	series    *techan.TimeSeries // k線序列
}

// timeframe -> atr計算長度
func NewAverageTrueRange(ctx context.Context, timeframe int) *AverageTrueRange {
	return &AverageTrueRange{
		Timeframe: timeframe,
		series:    techan.NewTimeSeries(),
	}
}

func (atr *AverageTrueRange) AddKline(kline *Kline) {
	AddKline(atr.series, kline)
}

func (atr *AverageTrueRange) Calculate() decimal.Decimal {
	atrInd := techan.NewAverageTrueRangeIndicator(atr.series, atr.Timeframe)
	v, _ := decimal.NewFromString(atrInd.Calculate(len(atr.series.Candles) - 1).String())
	return v
}
