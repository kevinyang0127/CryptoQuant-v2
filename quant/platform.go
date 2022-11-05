package quant

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/indicator"
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
收到資料後呼叫各個有訂閱的strategy來處理資料
*/

type Platform struct {
	mux             sync.Mutex
	platformID      string
	strategyManager *strategy.Manager
	runningStream   map[string]bool     //key: streamKey
	strategyIDCache map[string][]string // key: streamKey, value: strategyID list
}

func NewPlatform(mongoDB *db.MongoDB) *Platform {
	return &Platform{
		platformID:      util.GenIDWithPrefix("P_", 5),
		strategyManager: strategy.NewManager(mongoDB),
		runningStream:   make(map[string]bool),
		strategyIDCache: make(map[string][]string),
	}
}

/*
新增strategy，並監聽所需資料推送源
*/
func (p *Platform) AddStrategy(ctx context.Context, userID string, exchange string, symbol string,
	timeframe string, status strategy.StrategyStatus, strategyName string, script string) error {
	streamName, err := p.getStreamName(exchange)
	if err != nil {
		log.Println("p.getStreamName fail")
		return err
	}

	streamParam := stream.KlineStreamParam{
		Name:      streamName,
		Symbol:    symbol,
		Timeframe: timeframe,
	}

	go p.ListenNewKlineStream(ctx, streamParam)

	p.mux.Lock()
	defer p.mux.Unlock()

	strategyID, err := p.strategyManager.Add(ctx, userID, exchange, symbol, timeframe, status, strategyName, script)
	if err != nil {
		log.Println("strategyManager.Add fail")
		return err
	}

	streamKey := stream.GenKlineStreamKey(streamParam)
	p.strategyIDCache[streamKey] = append(p.strategyIDCache[streamKey], strategyID)

	return nil
}

func (p *Platform) getStreamName(exchange string) (stream.StreamName, error) {
	switch exchange {
	case "BINANCE_FUTURE":
		return stream.BinanceFKlineStream, nil
	default:
		log.Println("unsupport exchange name")
		return "", fmt.Errorf("unsupport exchange name")
	}
}

/*
監聽新的k線資料流
*/
func (p *Platform) ListenNewKlineStream(ctx context.Context, param stream.KlineStreamParam) {
	streamKey := stream.GenKlineStreamKey(param)
	isRunning := p.runningStream[streamKey]
	if isRunning {
		// 已經在監聽該資料流了
		return
	}

	ch := make(chan indicator.Kline)
	err := stream.KlineStreamManager.Subscribe(ctx, param, p.platformID, ch)
	if err != nil {
		log.Println("KlineStreamManager.Subscribe fail")
		log.Println(err)
		return
	}

	p.runningStream[streamKey] = true

	for {
		select {
		case <-ctx.Done():
			return
		case kline := <-ch:
			p.mux.Lock()
			strategyIDList := p.strategyIDCache[streamKey]
			for _, strategyID := range strategyIDList {
				info, err := p.strategyManager.GetStrategyInfo(ctx, strategyID)
				if err != nil {
					log.Println("strategyManager.GetStrategyInfo fail")
					log.Println(err)
					continue
				}

				if info.Status == strategy.Live {
					s, err := p.strategyManager.GetStrategyByID(ctx, info.StrategyID)
					if err != nil {
						log.Println(err)
						break
					}
					s.HandleKline(kline)
					break
				}
			}
			p.mux.Unlock()
		}
	}
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
