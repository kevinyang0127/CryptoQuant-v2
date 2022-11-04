package indicator

import (
	"context"

	"github.com/sdcoffey/techan"
	"github.com/shopspring/decimal"
)

/*
 參考trading view 特殊指標
 EMA & MA Crossover (by HPotter)
*/

type EmaMaCrossover struct {
	emaTimeframe int                // EMA計算長度
	maTimeframe  int                // MA計算長度
	series       *techan.TimeSeries // k線序列
}

// timeframe -> RSI計算長度
func NewEmaMaCrossover(ctx context.Context, emaTimeframe int, maTimeframe int) *EmaMaCrossover {
	return &EmaMaCrossover{
		emaTimeframe: emaTimeframe,
		maTimeframe:  maTimeframe,
		series:       techan.NewTimeSeries(),
	}
}

func (c *EmaMaCrossover) AddKline(kline *Kline) {
	AddKline(c.series, kline)
}

func (c *EmaMaCrossover) Calculate() (maVal decimal.Decimal, emaVal decimal.Decimal) {
	closePrices := techan.NewClosePriceIndicator(c.series)
	maInd := techan.NewSimpleMovingAverage(closePrices, c.maTimeframe)
	maVal, _ = decimal.NewFromString(maInd.Calculate(len(c.series.Candles) - 1).String())

	emaInd := techan.NewEMAIndicator(maInd, c.emaTimeframe)
	emaVal, _ = decimal.NewFromString(emaInd.Calculate(len(c.series.Candles) - 1).String())
	return
}
