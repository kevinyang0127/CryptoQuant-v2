package quant

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/exchange"
	"CryptoQuant-v2/market"
	"CryptoQuant-v2/strategy"
	"CryptoQuant-v2/stream"
	"CryptoQuant-v2/util"
	"context"
	"fmt"
	"log"
	"sync"
)

/*
Platform負責的事：
新增並監聽各個資料流(只限k線資料)
收到資料後呼叫strategy來處理資料
*/

type Platform struct {
	mux             sync.Mutex
	platformID      string
	strategyManager *strategy.Manager
	runningStream   map[string]bool     //key: streamKey
	liveStrategyID  map[string][]string // key: streamKey, value: strategyID list
}

func NewPlatform(mongoDB *db.MongoDB) *Platform {
	return &Platform{
		platformID:      util.GenIDWithPrefix("P_", 5),
		strategyManager: strategy.NewManager(mongoDB),
		runningStream:   make(map[string]bool),
		liveStrategyID:  make(map[string][]string),
	}
}

/*
新增strategy，並監聽所需資料推送源
*/
func (p *Platform) AddStrategy(ctx context.Context, userID string, exchange string, symbol string,
	timeframe string, status strategy.StrategyStatus, strategyName string, script string) (strategyID string, err error) {

	strategyID, err = p.strategyManager.Add(ctx, userID, exchange, symbol, timeframe, status, strategyName, script)
	if err != nil {
		log.Println("strategyManager.Add fail")
		return "", err
	}

	if status == strategy.Live {
		streamName, err := p.getStreamName(exchange)
		if err != nil {
			log.Println("p.getStreamName fail")
			return "", err
		}

		streamParam := stream.KlineStreamParam{
			Name:      streamName,
			Symbol:    symbol,
			Timeframe: timeframe,
		}

		p.mux.Lock()
		defer p.mux.Unlock()

		go p.ListenNewKlineStream(ctx, exchange, streamParam)

		streamKey := stream.GenKlineStreamKey(streamParam)
		p.liveStrategyID[streamKey] = append(p.liveStrategyID[streamKey], strategyID)
	}

	return strategyID, nil
}

func (p *Platform) getStreamName(exchangeName string) (stream.StreamName, error) {
	switch exchange.GetExchangeName(exchangeName) {
	case exchange.BINANCE_FUTURE:
		return stream.BinanceFKlineStream, nil
	default:
		log.Println("unsupport exchange name")
		return "", fmt.Errorf("unsupport exchange name")
	}
}

/*
監聽新的k線資料流
*/
func (p *Platform) ListenNewKlineStream(ctx context.Context, exchange string, param stream.KlineStreamParam) {
	streamKey := stream.GenKlineStreamKey(param)
	isRunning := p.runningStream[streamKey]
	if isRunning {
		// 已經在監聽該資料流了
		return
	}

	ch := make(chan market.Kline)
	err := stream.KlineStreamManager.Subscribe(ctx, param, p.platformID, ch)
	if err != nil {
		log.Println("KlineStreamManager.Subscribe fail")
		log.Println(err)
		return
	}

	p.runningStream[streamKey] = true

	// 取得前500根k線資料
	klines, err := p.getLimitKlineHistory(ctx, exchange, param.Symbol, param.Timeframe, 500)
	if err != nil {
		log.Println("getLimitKlineHistory fail")
		log.Println(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case kline := <-ch:
			if kline.IsFinal {
				klines = append(klines, kline)
			}

			// 超過2000筆只保留最新500筆
			if len(klines) >= 2000 {
				newSlice := make([]market.Kline, 500, 2000)
				copy(newSlice, klines[len(klines)-500:])
				klines = newSlice
			}

			p.mux.Lock()
			liveStrategyIDList := p.liveStrategyID[streamKey]
			for _, liveStrategyID := range liveStrategyIDList {
				s, err := p.strategyManager.GetStrategyByID(ctx, liveStrategyID)
				if err != nil {
					log.Println("strategyManager.GetStrategyByID fail")
					log.Println(err)
					continue
				}
				s.HandleKline(klines, kline)
			}
			p.mux.Unlock()
		}
	}
}

// 取得最新的數根k線
func (p *Platform) getLimitKlineHistory(ctx context.Context, exchangeName string, symbol string, timeframe string, limit int) ([]market.Kline, error) {
	ex, err := exchange.GetExchange(exchangeName)
	if err != nil {
		log.Println("GetExchange fail")
		return nil, err
	}
	klines, err := ex.GetLimitKlineHistory(ctx, symbol, timeframe, limit)
	if err != nil {
		log.Println("GetHistoryKline fail")
		return nil, err
	}
	return klines, nil
}

/*
停止已經存在的strategy，使其不會被呼叫
*/
func (p *Platform) StopStrategy() {

}

/*
恢復已經被停止的strategy
*/
func (p *Platform) RecoverStrategy() {

}

/*
移除已經存在的strategy
*/
func (p *Platform) RemoveStrategy() {

}
