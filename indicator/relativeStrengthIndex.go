package indicator

import (
	"context"

	"github.com/sdcoffey/techan"
	"github.com/shopspring/decimal"
)

type RelativeStrengthIndex struct {
	Timeframe int                // RSI計算長度
	series    *techan.TimeSeries // k線序列
}

// timeframe -> RSI計算長度
func NewRelativeStrengthIndex(ctx context.Context, timeframe int) *RelativeStrengthIndex {
	return &RelativeStrengthIndex{
		Timeframe: timeframe,
		series:    techan.NewTimeSeries(),
	}
}

func (rsi *RelativeStrengthIndex) AddKline(kline *Kline) {
	AddKline(rsi.series, kline)
}

func (rsi *RelativeStrengthIndex) Calculate() decimal.Decimal {
	closePrices := techan.NewClosePriceIndicator(rsi.series)
	rsiInd := techan.NewRelativeStrengthIndexIndicator(closePrices, rsi.Timeframe)
	//計算最新的rsi值
	v, _ := decimal.NewFromString(rsiInd.Calculate(len(rsi.series.Candles) - 1).String())
	return v
}
