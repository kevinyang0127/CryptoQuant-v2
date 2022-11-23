package market

import "github.com/shopspring/decimal"

type Position struct {
	Symbol    string
	Quantity  decimal.Decimal // 數量，負的代表做空
	OpenPrice decimal.Decimal // 開倉的平均價格
}
