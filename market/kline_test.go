package market

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestXxx(t *testing.T) {
	kline := Kline{
		High:    decimal.NewFromFloat(1000),
		Low:     decimal.NewFromFloat(500),
		Open:    decimal.NewFromFloat(900),
		Close:   decimal.NewFromFloat(800),
		IsFinal: true,
	}
	// GenFinalKlinePath(kline, 10)
	klineList, _ := GenFakeFinalKlinePath(kline, 20)
	fmt.Println(len(klineList))
	for _, v := range klineList {
		fmt.Println(v)
	}

	t.Errorf("errrrr")
}

func TestGenFinalKlineHistory(t *testing.T) {
	smallKline1 := Kline{
		StartTime:            100,
		EndTime:              199,
		High:                 decimal.NewFromFloat(1135.20),
		Low:                  decimal.NewFromFloat(1131.70),
		Open:                 decimal.NewFromFloat(1132.22),
		Close:                decimal.NewFromFloat(1133.74),
		Volume:               decimal.NewFromFloat(10),
		QuoteVolume:          decimal.NewFromFloat(10),
		ActiveBuyVolume:      decimal.NewFromFloat(10),
		ActiveBuyQuoteVolume: decimal.NewFromFloat(10),
		IsFinal:              true,
	}

	smallKline2 := Kline{
		StartTime:            200,
		EndTime:              299,
		High:                 decimal.NewFromFloat(1134.58),
		Low:                  decimal.NewFromFloat(1132.92),
		Open:                 decimal.NewFromFloat(1133.74),
		Close:                decimal.NewFromFloat(1132.98),
		Volume:               decimal.NewFromFloat(10),
		QuoteVolume:          decimal.NewFromFloat(10),
		ActiveBuyVolume:      decimal.NewFromFloat(10),
		ActiveBuyQuoteVolume: decimal.NewFromFloat(10),
		IsFinal:              true,
	}

	smallKline3 := Kline{
		StartTime:            300,
		EndTime:              399,
		High:                 decimal.NewFromFloat(1133.72),
		Low:                  decimal.NewFromFloat(1130.36),
		Open:                 decimal.NewFromFloat(1132.97),
		Close:                decimal.NewFromFloat(1130.60),
		Volume:               decimal.NewFromFloat(10),
		QuoteVolume:          decimal.NewFromFloat(10),
		ActiveBuyVolume:      decimal.NewFromFloat(10),
		ActiveBuyQuoteVolume: decimal.NewFromFloat(10),
		IsFinal:              true,
	}

	finalKline := Kline{
		StartTime:            100,
		EndTime:              399,
		High:                 decimal.NewFromFloat(1135.20),
		Low:                  decimal.NewFromFloat(1130.36),
		Open:                 decimal.NewFromFloat(1132.22),
		Close:                decimal.NewFromFloat(1130.60),
		Volume:               decimal.NewFromFloat(30),
		QuoteVolume:          decimal.NewFromFloat(30),
		ActiveBuyVolume:      decimal.NewFromFloat(30),
		ActiveBuyQuoteVolume: decimal.NewFromFloat(30),
		IsFinal:              true,
	}

	// GenFinalKlinePath(kline, 10)
	klineList, _ := GenFinalKlineHistory(finalKline, []Kline{smallKline1, smallKline2, smallKline3})
	fmt.Println(len(klineList))
	for _, v := range klineList {
		fmt.Println(v)
	}

	t.Errorf("errrrr")
}
