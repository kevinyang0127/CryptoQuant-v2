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
	klineList, _ := GenFinalKlinePath(kline, 20)
	fmt.Println(len(klineList))
	for _, v := range klineList {
		fmt.Println(v)
	}

	t.Errorf("errrrr")
}
