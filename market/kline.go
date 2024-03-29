package market

import (
	"fmt"
	"math/rand"
	"time"

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

/*
輸入一根已經收盤的k線，輸出此根k線的產生過程
[notFinalKline, notFinalKline, notFinalKline, ..., finalKline]
precision決定在high和low之間要分成多少點位
*/
func GenFakeFinalKlinePath(finalKline Kline, precision int) ([]Kline, error) {
	if !finalKline.IsFinal {
		return nil, fmt.Errorf("kline is not final")
	}

	if precision <= 1 {
		return []Kline{finalKline}, nil
	}

	d := finalKline.High.Sub(finalKline.Low).Div(decimal.NewFromInt(int64(precision)))

	positions := []decimal.Decimal{}
	for i := 0; i <= precision; i++ {
		if i == precision {
			positions = append(positions, finalKline.High)
		} else {
			positions = append(positions, finalKline.Low.Add(d.Mul(decimal.NewFromInt(int64(i)))))
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(positions), func(i, j int) { positions[i], positions[j] = positions[j], positions[i] })
	pricePath := []decimal.Decimal{}
	pricePath = append(pricePath, finalKline.Open)
	pricePath = append(pricePath, positions...)
	pricePath = append(pricePath, finalKline.Close)

	pricePath = smoothPath(pricePath, d)

	startKline := Kline{
		StartTime: finalKline.StartTime,
		EndTime:   finalKline.EndTime,
		TradeNum:  finalKline.TradeNum,
		Open:      finalKline.Open,
		Close:     finalKline.Open,
		High:      finalKline.Open,
		Low:       finalKline.Open,
		IsFinal:   false,
	}
	klineList := []Kline{startKline}
	for i, nowPrice := range pricePath {
		if i == 0 || i == len(pricePath)-1 {
			continue
		}
		newKline := drawNewKline(klineList[i-1], nowPrice)
		klineList = append(klineList, newKline)
	}
	klineList = append(klineList, finalKline)

	return klineList, nil
}

func drawNewKline(prevKline Kline, nowPrice decimal.Decimal) Kline {
	newKline := prevKline
	newKline.Close = nowPrice
	if nowPrice.GreaterThanOrEqual(prevKline.High) {
		newKline.High = nowPrice
	} else if nowPrice.LessThanOrEqual(prevKline.Low) {
		newKline.Low = nowPrice
	}
	return newKline
}

func smoothPath(path []decimal.Decimal, d decimal.Decimal) []decimal.Decimal {
	prev := path[0]
	for i := 0; i < len(path); i++ {
		for prev.Add(d).LessThan(path[i]) {
			path = append(path[:i+1], path[i:]...)
			path[i] = prev.Add(d)
			prev = prev.Add(d)
			i++
		}
		for prev.Sub(d).GreaterThan(path[i]) {
			path = append(path[:i+1], path[i:]...)
			path[i] = prev.Sub(d)
			prev = prev.Sub(d)
			i++
		}
		prev = path[i]
	}
	return path
}

// 輸入一根時間框架較高的k線 和 他的開始時間到結束時間內的小時間框架的k線序列，輸出高時間框架的k線變化序列
func GenFinalKlineHistory(finalKline Kline, smallTimeframeKlines []Kline) ([]Kline, error) {
	finalKlineHistory := []Kline{}
	for i, sk := range smallTimeframeKlines {
		if sk.StartTime < finalKline.StartTime || sk.EndTime > finalKline.EndTime {
			return nil, fmt.Errorf("smallTimeframeKline's startTime or endTime out of finalKline")
		}

		if i == len(smallTimeframeKlines)-1 {
			finalKlineHistory = append(finalKlineHistory, finalKline)
			break
		}

		newKline := Kline{
			StartTime:            finalKline.StartTime,
			EndTime:              finalKline.EndTime,
			Open:                 sk.Open,
			Close:                sk.Close,
			High:                 sk.High,
			Low:                  sk.Low,
			Volume:               sk.Volume,
			QuoteVolume:          sk.QuoteVolume,
			ActiveBuyVolume:      sk.ActiveBuyVolume,
			ActiveBuyQuoteVolume: sk.ActiveBuyQuoteVolume,
			IsFinal:              false,
		}
		if len(finalKlineHistory) > 0 {
			newKline.Open = finalKlineHistory[i-1].Open
			if finalKlineHistory[i-1].High.GreaterThan(newKline.High) {
				newKline.High = finalKlineHistory[i-1].High
			}
			if finalKlineHistory[i-1].Low.LessThan(newKline.Low) {
				newKline.Low = finalKlineHistory[i-1].Low
			}
			newKline.Volume = finalKlineHistory[i-1].Volume.Add(newKline.Volume)
			newKline.QuoteVolume = finalKlineHistory[i-1].QuoteVolume.Add(newKline.QuoteVolume)
			newKline.ActiveBuyVolume = finalKlineHistory[i-1].ActiveBuyVolume.Add(newKline.ActiveBuyVolume)
			newKline.ActiveBuyQuoteVolume = finalKlineHistory[i-1].ActiveBuyQuoteVolume.Add(newKline.ActiveBuyQuoteVolume)
		}
		finalKlineHistory = append(finalKlineHistory, newKline)
	}
	return finalKlineHistory, nil
}
