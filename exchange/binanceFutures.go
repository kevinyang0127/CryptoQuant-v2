package exchange

import (
	"CryptoQuant-v2/market"
	"context"
	"fmt"
	"log"

	binanceFutures "github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
)

type BinanceFuture struct {
	clinet               *binanceFutures.Client
	exchangeInfo         *binanceFutures.ExchangeInfo
	pricePrecisionMap    map[string]int32
	quantityPrecisionMap map[string]int32
}

func newBinanceFuture(ctx context.Context, apiKey, secretKey string) (*BinanceFuture, error) {
	binanceFuture := &BinanceFuture{
		clinet:               binanceFutures.NewClient(apiKey, secretKey),
		pricePrecisionMap:    make(map[string]int32),
		quantityPrecisionMap: make(map[string]int32),
	}

	res, err := binanceFuture.clinet.NewExchangeInfoService().Do(ctx)
	if err != nil {
		log.Println("binanceFuture.clinet.NewExchangeInfoService().Do() fail")
		return nil, err
	}

	binanceFuture.exchangeInfo = res

	return binanceFuture, nil
}

func (bf *BinanceFuture) getPricePrecision(symbol string) int32 {
	v, ok := bf.pricePrecisionMap[symbol]
	if ok {
		return v
	}

	for _, v := range bf.exchangeInfo.Symbols {
		if v.Symbol == symbol {
			bf.pricePrecisionMap[symbol] = int32(v.PricePrecision)
			return int32(v.PricePrecision)
		}
	}
	return 6 // default value
}

func (bf *BinanceFuture) getQuantityPrecision(symbol string) int32 {
	v, ok := bf.quantityPrecisionMap[symbol]
	if ok {
		return v
	}

	for _, v := range bf.exchangeInfo.Symbols {
		if v.Symbol == symbol {
			bf.quantityPrecisionMap[symbol] = int32(v.QuantityPrecision)
			return int32(v.QuantityPrecision)
		}
	}
	return 6 // default value
}

func (bf *BinanceFuture) GetLimitKlineHistory(ctx context.Context, symbol string, timeframe string, limit int) ([]market.Kline, error) {
	klinesService := bf.clinet.NewKlinesService()
	klinesService.Symbol(symbol)
	klinesService.Interval(timeframe)
	klinesService.Limit(limit)
	res, err := klinesService.Do(ctx)
	if err != nil {
		log.Println("klinesService.Do fail")
		return nil, err
	}

	klines := []market.Kline{}
	for i, k := range res {
		if i == len(res)-1 {
			//檔下尚未收盤的k線也會拿到，所以最後一根k線不做事，因為高機率初始化時不會剛好收盤
			break
		}

		kline, err := market.BinanceFKlineToKline(k)
		if err != nil {
			log.Println("market.BinanceFKlineToKline fail")
			return nil, err
		}

		klines = append(klines, *kline)
	}

	return klines, nil
}

func (bf *BinanceFuture) GetLimitKlineHistoryByTime(ctx context.Context, symbol string, timeframe string, limit int, startTimeMs int64, endTimeMs int64) ([]market.Kline, error) {
	klinesService := bf.clinet.NewKlinesService()
	klinesService.Symbol(symbol)
	klinesService.Interval(timeframe)
	klinesService.Limit(limit)
	klinesService.StartTime(startTimeMs)
	klinesService.EndTime(endTimeMs)
	res, err := klinesService.Do(ctx)
	if err != nil {
		log.Println("klinesService.Do fail")
		return nil, err
	}

	klines := []market.Kline{}
	for i, k := range res {
		if i == len(res)-1 {
			//檔下尚未收盤的k線也會拿到，所以最後一根k線不做事，因為高機率初始化時不會剛好收盤
			break
		}

		kline, err := market.BinanceFKlineToKline(k)
		if err != nil {
			log.Println("market.BinanceFKlineToKline fail")
			return nil, err
		}

		klines = append(klines, *kline)
	}

	return klines, nil
}

func (bf *BinanceFuture) CreateMarketOrder(ctx context.Context, symbol string, side bool, quantity decimal.Decimal) error {
	if quantity.IsNegative() {
		return fmt.Errorf("quantity is negative error")
	}

	var sideType binanceFutures.SideType
	if side {
		sideType = binanceFutures.SideTypeBuy
	} else {
		sideType = binanceFutures.SideTypeSell
	}
	createOrderServ := bf.clinet.NewCreateOrderService()
	createOrderServ.Symbol(symbol)
	createOrderServ.Side(sideType)
	createOrderServ.Type(binanceFutures.OrderTypeMarket)
	createOrderServ.Quantity(quantity.StringFixed(bf.getQuantityPrecision(symbol)))
	_, err := createOrderServ.Do(ctx)
	if err != nil {
		log.Println("CreateMarketOrder request fail")
		return err
	}

	return nil
}

func (bf *BinanceFuture) CreateLimitOrder(ctx context.Context, symbol string, side bool, price decimal.Decimal, quantity decimal.Decimal) error {
	if price.IsNegative() {
		return fmt.Errorf("price is negative error")
	}
	if quantity.IsNegative() {
		return fmt.Errorf("quantity is negative error")
	}

	var sideType binanceFutures.SideType
	if side {
		sideType = binanceFutures.SideTypeBuy
	} else {
		sideType = binanceFutures.SideTypeSell
	}
	createOrderServ := bf.clinet.NewCreateOrderService()
	createOrderServ.Symbol(symbol)
	createOrderServ.Side(sideType)
	createOrderServ.Type(binanceFutures.OrderTypeLimit)
	createOrderServ.TimeInForce(binanceFutures.TimeInForceTypeGTC)
	createOrderServ.Price(price.StringFixed(bf.getPricePrecision(symbol)))
	createOrderServ.Quantity(quantity.StringFixed(bf.getQuantityPrecision(symbol)))
	_, err := createOrderServ.Do(ctx)
	if err != nil {
		log.Println("CreateLimitOrder request fail")
		return err
	}

	return nil
}

func (bf *BinanceFuture) CancelAllOpenOrders(ctx context.Context, symbol string) error {
	cancelAllOpenOrdersServ := bf.clinet.NewCancelAllOpenOrdersService()
	cancelAllOpenOrdersServ.Symbol(symbol)
	err := cancelAllOpenOrdersServ.Do(ctx)
	if err != nil {
		log.Println("cancelAllOpenOrdersServ request fail")
		return err
	}

	return nil
}

func (bf *BinanceFuture) GetPosition(ctx context.Context, symbol string) (*market.Position, error) {
	getPositionRiskServ := bf.clinet.NewGetPositionRiskService()
	getPositionRiskServ.Symbol(symbol)
	res, err := getPositionRiskServ.Do(ctx)
	if err != nil {
		log.Println("GetPositionRiskServ request fail")
		return nil, err
	}
	for _, v := range res {
		if v.Symbol == symbol && v.PositionSide == "BOTH" {
			qty, err := decimal.NewFromString(v.PositionAmt)
			if err != nil {
				log.Println("decimal.NewFromString fail")
				log.Println(err)
				continue
			}
			entryPrice, err := decimal.NewFromString(v.EntryPrice)
			if err != nil {
				log.Println("decimal.NewFromString fail")
				log.Println(err)
				continue
			}
			return &market.Position{
				Symbol:    v.Symbol,
				Quantity:  qty,
				OpenPrice: entryPrice,
			}, nil
		}
	}
	return nil, nil
}

func (bf *BinanceFuture) CreateStopLossOrder(ctx context.Context, symbol string, side bool, price decimal.Decimal, quantity decimal.Decimal, stopPrice decimal.Decimal) error {
	if price.IsNegative() {
		return fmt.Errorf("price is negative error")
	}
	if stopPrice.IsNegative() {
		return fmt.Errorf("stopPrice is negative error")
	}
	if quantity.IsNegative() {
		return fmt.Errorf("quantity is negative error")
	}

	var sideType binanceFutures.SideType
	if side {
		sideType = binanceFutures.SideTypeBuy
	} else {
		sideType = binanceFutures.SideTypeSell
	}
	createOrderServ := bf.clinet.NewCreateOrderService()
	createOrderServ.Symbol(symbol)
	createOrderServ.Side(sideType)
	createOrderServ.Type(binanceFutures.OrderTypeStop)
	createOrderServ.Price(price.StringFixed(bf.getPricePrecision(symbol)))
	createOrderServ.Quantity(quantity.StringFixed(bf.getQuantityPrecision(symbol)))
	createOrderServ.StopPrice(stopPrice.StringFixed(bf.getPricePrecision(symbol)))
	_, err := createOrderServ.Do(ctx)
	if err != nil {
		log.Println("CreateStopLossOrder request fail")
		return err
	}

	return nil
}

func (bf *BinanceFuture) CreateTakeProfitOrder(ctx context.Context, symbol string, side bool, price decimal.Decimal, quantity decimal.Decimal, stopPrice decimal.Decimal) error {
	if price.IsNegative() {
		return fmt.Errorf("price is negative error")
	}
	if stopPrice.IsNegative() {
		return fmt.Errorf("stopPrice is negative error")
	}
	if quantity.IsNegative() {
		return fmt.Errorf("quantity is negative error")
	}

	var sideType binanceFutures.SideType
	if side {
		sideType = binanceFutures.SideTypeBuy
	} else {
		sideType = binanceFutures.SideTypeSell
	}
	createOrderServ := bf.clinet.NewCreateOrderService()
	createOrderServ.Symbol(symbol)
	createOrderServ.Side(sideType)
	createOrderServ.Type(binanceFutures.OrderTypeTakeProfit)
	createOrderServ.Price(price.StringFixed(bf.getPricePrecision(symbol)))
	createOrderServ.Quantity(quantity.StringFixed(bf.getQuantityPrecision(symbol)))
	createOrderServ.StopPrice(stopPrice.StringFixed(bf.getPricePrecision(symbol)))
	_, err := createOrderServ.Do(ctx)
	if err != nil {
		log.Println("CreateTakeProfitOrder request fail")
		return err
	}

	return nil
}
