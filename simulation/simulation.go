package simulation

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/market"
	"CryptoQuant-v2/util"
	"context"
	"log"
	"sync"

	"github.com/shopspring/decimal"
)

/*
Simulation 模擬一個用戶執行一個策略的交易，倉位變化及盈利
*/
type Simulation struct {
	mongoDB             *db.MongoDB
	mux                 sync.Mutex
	userID              string
	simulationID        string
	startBalance        decimal.Decimal
	lever               decimal.Decimal
	takerCommissionRate decimal.Decimal
	makerCommissionRate decimal.Decimal
	balance             decimal.Decimal
	positon             *Position
	orders              []*Order
}

type Position struct {
	Quantity  decimal.Decimal // 數量，負的代表做空
	OpenPrice decimal.Decimal // 開倉時的價格，如果有加倉則會算出平均價
}

type OrderType int

const (
	Normal OrderType = iota
	StopLoss
	TakeProfit
)

type Order struct {
	Type      OrderType
	Quantity  decimal.Decimal
	Price     decimal.Decimal
	StopPrice decimal.Decimal //若為停損停利單的觸發價格
}

type TradeLog struct {
	SimulationID string `bson:"simulationID"`
	Timestamp    int64  `bson:"timestamp"`
	Action       string `bson:"action"`
	Side         bool   `bson:"side"`
	Price        string `bson:"price"`
	Quantity     string `bson:"quantity"`
	Commission   string `bson:"commission"`
	Profit       string `bson:"profit"`
	Balance      string `bson:"balance"`
	Msg          string `bson:"msg"`
}

func NewSimulation(mongoDB *db.MongoDB, userID string, startBalance decimal.Decimal,
	lever decimal.Decimal, takerCommissionRate decimal.Decimal, makerCommissionRate decimal.Decimal) *Simulation {
	return &Simulation{
		mongoDB:             mongoDB,
		userID:              userID,
		simulationID:        util.GenIDWithPrefix("sim_", 5),
		startBalance:        startBalance,
		lever:               lever,
		takerCommissionRate: takerCommissionRate,
		makerCommissionRate: makerCommissionRate,
		balance:             startBalance,
		orders:              []*Order{},
	}
}

// 監聽新的kline，處理掛單成交與否
func (s *Simulation) ListenNewKline(ctx context.Context, ch chan market.Kline) {
	for {
		select {
		case <-ctx.Done():
			return
		case kline := <-ch:
			if len(s.orders) != 0 {
				s.checkOrderMatch(ctx, kline)
			}
		}
	}
}

func (s *Simulation) checkOrderMatch(ctx context.Context, kline market.Kline) {
	doneOrderIndexs := make(map[int]bool)
	someOrderDone := false
	for i, order := range s.orders {
		switch order.Type {
		case Normal:
			// 買入只要訂單價錢高於目前價格，或賣出只要訂單價錢低於目前價格都會觸發成交
			if (order.Quantity.IsPositive() && order.Price.GreaterThanOrEqual(kline.Low)) ||
				(order.Quantity.IsNegative() && order.Price.LessThanOrEqual(kline.High)) {
				// 掛單價位在k線範圍當中，就會以掛單價位成交，不然就是以收盤價成交
				matchPrice := kline.Close
				isMaker := false
				if order.Price.GreaterThanOrEqual(kline.Low) && order.Price.LessThanOrEqual(kline.High) {
					matchPrice = order.Price
					isMaker = true
				}

				if s.positon == nil || s.positon.Quantity.Mul(order.Quantity).IsPositive() {
					// 開倉/加倉
					log.Println("Normal order Entry, order = ", order)
					s.Entry(ctx, matchPrice, order.Quantity, isMaker, kline.EndTime)
				} else {
					// 關倉/減倉
					log.Println("Normal order Exit, order = ", order)
					s.Exit(ctx, matchPrice, order.Quantity.Abs(), isMaker, kline.EndTime)
				}
				doneOrderIndexs[i] = true
				someOrderDone = true
			}
		case StopLoss:
			// 多單當k線最低價低於觸發價 或 空單當k線最高價高於觸發價 時，會觸發成交
			// 限價停損掛單是賣出代表原本是買入
			if (order.Quantity.IsNegative() && kline.Low.LessThanOrEqual(order.StopPrice)) ||
				(order.Quantity.IsPositive() && kline.High.GreaterThanOrEqual(order.StopPrice)) {
				//FIXME: 如果希望達到order.StopPrice後掛單掛上合理的order.Price呢？
				//ex: 做多時 停損觸發在100元賣出，但掛單掛120元賣出（若掛單掛在100元以下賣出都會以100元馬上成交）
				log.Println("StopLoss order exit, order = ", order)
				s.Exit(ctx, order.StopPrice, order.Quantity.Abs(), true, kline.EndTime)
				doneOrderIndexs[i] = true
				someOrderDone = true
			}
		case TakeProfit:
			// 多單當k線最高價高於觸發價 或 空單當k線最低價低於觸發價 時，會觸發成交
			// 限價停利掛單是賣出代表原本是買入
			if order.Quantity.IsNegative() && kline.High.GreaterThanOrEqual(order.StopPrice) ||
				(order.Quantity.IsPositive() && kline.Low.LessThanOrEqual(order.StopPrice)) {
				//FIXME: 如果希望達到order.StopPrice後掛單掛上合理的order.Price呢？
				//ex: 做多時 停利觸發在100元賣出，但掛單掛120元賣出（若掛單掛在100元以下賣出都會以100元馬上成交）
				log.Println("TakeProfit order exit, order = ", order)
				s.Exit(ctx, order.StopPrice, order.Quantity.Abs(), true, kline.EndTime)
				doneOrderIndexs[i] = true
				someOrderDone = true
			}
		}
	}

	if someOrderDone {
		// remove done order
		newOrders := []*Order{}
		for i := 0; i < len(s.orders); i++ {
			if !doneOrderIndexs[i] {
				newOrders = append(newOrders, s.orders[i])
			}
		}
		s.orders = newOrders
	}
}

// 開倉/加倉, quantity為正代表買入，為負代表賣出， isMaker-> true為maker, false為taker
func (s *Simulation) Entry(ctx context.Context, price decimal.Decimal, quantity decimal.Decimal, isMaker bool, klineTimestamp int64) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if price.IsNegative() {
		log.Println("SimulationEngine Entry fail, price IsNegative")
		return
	}

	spend := price.Mul(quantity).Div(s.lever).Abs()
	fee := decimal.Zero
	if isMaker {
		fee = spend.Mul(s.makerCommissionRate)
	} else {
		fee = spend.Mul(s.takerCommissionRate)
	}
	if spend.Add(fee).GreaterThan(s.balance) {
		tradelog := &TradeLog{
			Balance: s.balance.StringFixed(6),
			Msg:     "balance not enough",
		}
		_, err := s.mongoDB.InsertOne(ctx, db.DBNAME, "simulationLog", tradelog)
		if err != nil {
			log.Println("mongoDB.InsertOne fail")
			log.Println(err)
		}
		return
	}

	// 已經開倉，但加倉方向相反
	if s.positon != nil && s.positon.Quantity.Mul(quantity).IsNegative() {
		log.Println("SimulationEngine Entry fail, add position wrong side")
		return
	}

	s.balance = s.balance.Sub(spend).Sub(fee)

	if s.positon != nil {
		s.positon.OpenPrice = s.positon.OpenPrice.Mul(s.positon.Quantity).Add(price.Mul(quantity)).Div(s.positon.Quantity.Add(quantity))
		s.positon.Quantity = s.positon.Quantity.Add(quantity)
	} else {
		s.positon = &Position{
			Quantity:  quantity,
			OpenPrice: price,
		}
	}

	tradelog := &TradeLog{
		SimulationID: s.simulationID,
		Timestamp:    klineTimestamp,
		Action:       "ENTRY",
		Side:         quantity.IsPositive(),
		Price:        price.StringFixed(6),
		Quantity:     quantity.StringFixed(6),
		Commission:   fee.StringFixed(6),
		Profit:       decimal.Zero.StringFixed(6),
		Balance:      s.balance.StringFixed(6),
		Msg:          "Entry success",
	}
	_, err := s.mongoDB.InsertOne(ctx, db.DBNAME, "simulationLog", tradelog)
	if err != nil {
		log.Println("mongoDB.InsertOne fail")
		log.Println(err)
	}
}

// 關倉/減倉, quantity只能為正，會自動處理成相反方向
func (s *Simulation) Exit(ctx context.Context, price decimal.Decimal, quantity decimal.Decimal, isMaker bool, klineTimestamp int64) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if price.IsNegative() {
		log.Println("SimulationEngine Exit fail, price IsNegative")
		return
	}

	if quantity.IsNegative() {
		log.Println("SimulationEngine Exit fail, quantity IsNegative")
		return
	}

	if s.positon == nil {
		log.Println("SimulationEngine Exit fail, no position to close")
		return
	}

	if quantity.Abs().GreaterThan(s.positon.Quantity.Abs()) {
		log.Println("SimulationEngine Exit fail, not enough position to close")
		return
	}

	if s.positon.Quantity.IsPositive() {
		quantity = quantity.Mul(decimal.NewFromInt(-1))
	}

	fee := decimal.Zero
	if isMaker {
		fee = price.Mul(quantity).Div(s.lever).Abs().Mul(s.makerCommissionRate)
	} else {
		fee = price.Mul(quantity).Div(s.lever).Abs().Mul(s.takerCommissionRate)
	}

	s.balance = s.balance.Sub(fee)

	originSpend := s.positon.OpenPrice.Mul(s.positon.Quantity).Div(s.lever).Abs()
	openSide := s.positon.Quantity.IsPositive()
	profit := decimal.Zero
	// 賣出-買入
	if openSide {
		profit = price.Mul(quantity).Abs().Sub(s.positon.OpenPrice.Mul(quantity.Abs()))
	} else {
		profit = s.positon.OpenPrice.Mul(quantity).Abs().Sub(price.Mul(quantity))
	}

	s.balance = s.balance.Add(originSpend).Add(profit)

	s.positon.Quantity = s.positon.Quantity.Add(quantity)
	if s.positon.Quantity.IsZero() {
		s.positon = nil
	}

	tradelog := &TradeLog{
		SimulationID: s.simulationID,
		Timestamp:    klineTimestamp,
		Action:       "EXIT",
		Side:         quantity.IsPositive(),
		Price:        price.StringFixed(6),
		Quantity:     quantity.StringFixed(6),
		Commission:   fee.StringFixed(6),
		Profit:       profit.StringFixed(6),
		Balance:      s.balance.StringFixed(6),
		Msg:          "Exit success",
	}
	_, err := s.mongoDB.InsertOne(ctx, db.DBNAME, "simulationLog", tradelog)
	if err != nil {
		log.Println("mongoDB.InsertOne fail")
		log.Println(err)
	}
}

// 限價掛單, quantity為正代表買入，為負代表賣出
func (s *Simulation) Order(ctx context.Context, price decimal.Decimal, quantity decimal.Decimal) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.orders = append(s.orders, &Order{
		Type:      Normal,
		Quantity:  quantity,
		Price:     price,
		StopPrice: decimal.Zero,
	})
}

// 限價停損掛單, quantity為正代表買入，為負代表賣出
func (s *Simulation) StopLossOrder(ctx context.Context, price decimal.Decimal, quantity decimal.Decimal, stopPrice decimal.Decimal) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.orders = append(s.orders, &Order{
		Type:      StopLoss,
		Quantity:  quantity,
		Price:     price,
		StopPrice: stopPrice,
	})
}

// 限價停利掛單, quantity為正代表買入，為負代表賣出
func (s *Simulation) TakeProfitOrder(ctx context.Context, price decimal.Decimal, quantity decimal.Decimal, stopPrice decimal.Decimal) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.orders = append(s.orders, &Order{
		Type:      TakeProfit,
		Quantity:  quantity,
		Price:     price,
		StopPrice: stopPrice,
	})
}

// 撤銷所有掛單
func (s *Simulation) CloseAllOrder(ctx context.Context) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.orders = []*Order{}
}

// 取得所有掛單
func (s *Simulation) GetAllOrder(ctx context.Context) []*Order {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.orders
}

// 取得目前倉位
func (s *Simulation) GetPosition(ctx context.Context) *Position {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.positon
}

// 取得目前餘額
func (s *Simulation) GetBalance(ctx context.Context) decimal.Decimal {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.balance
}
