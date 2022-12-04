package quant

import (
	"CryptoQuant-v2/exchange"
	"CryptoQuant-v2/market"
	"CryptoQuant-v2/strategy"
	"CryptoQuant-v2/stream"
	"CryptoQuant-v2/user"
	"context"
	"fmt"
	"log"
	"sync"

	mapset "github.com/deckarep/golang-set"
)

/*
Platform負責的事：
新增並監聽各個資料流(只限k線資料)
收到資料後呼叫strategy來處理資料
*/

type Platform struct {
	mux             sync.Mutex
	strategyManager *strategy.Manager
	exchangeManager *exchange.Manager
	userManager     *user.Manager
	runningStream   mapset.Set            // streamKey set
	liveStrategyID  map[string]mapset.Set // key: streamKey, value: strategyID set
}

func NewPlatform(strategyManager *strategy.Manager, exchangeManager *exchange.Manager, userManager *user.Manager) *Platform {
	return &Platform{
		strategyManager: strategyManager,
		exchangeManager: exchangeManager,
		userManager:     userManager,
		runningStream:   mapset.NewSet(),
		liveStrategyID:  make(map[string]mapset.Set),
	}
}

/*
執行strategy，並監聽所需資料推送源
*/
func (p *Platform) RunStrategy(ctx context.Context, strategyID string) error {
	info, err := p.strategyManager.GetStrategyInfo(ctx, strategyID)
	if err != nil {
		log.Println("strategyManager.GetStrategyInfo fail")
		return err
	}

	streamName, err := p.getStreamName(info.Exchange)
	if err != nil {
		log.Println("p.getStreamName fail")
		return err
	}

	streamParam := stream.KlineStreamParam{
		Name:      streamName,
		Symbol:    info.Symbol,
		Timeframe: info.Timeframe,
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	go p.ListenAndHandleKlineStream(ctx, info.Exchange, streamParam)

	streamKey := stream.GenKlineStreamKey(streamParam)
	_, ok := p.liveStrategyID[streamKey]
	if !ok {
		p.liveStrategyID[streamKey] = mapset.NewSet()
	}

	if info.Status != strategy.Live {
		err = p.strategyManager.UpdateStatus(ctx, strategyID, strategy.Live)
		if err != nil {
			log.Println("strategyManager.UpdateStatus fail")
			return err
		}
	}

	p.liveStrategyID[streamKey].Add(strategyID)

	return nil
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
func (p *Platform) ListenAndHandleKlineStream(ctx context.Context, exchange string, param stream.KlineStreamParam) {
	streamKey := stream.GenKlineStreamKey(param)
	if p.runningStream.Contains(streamKey) {
		// 已經在監聽該資料流了
		return
	}

	ch := make(chan market.Kline)
	err := stream.KlineStreamManager.Subscribe(ctx, param, "PLATFORM", ch)
	if err != nil {
		log.Println("KlineStreamManager.Subscribe fail")
		log.Println(err)
		return
	}

	p.runningStream.Add(streamKey)

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
			liveStrategyIDSet := p.liveStrategyID[streamKey]
			for _, liveStrategyID := range liveStrategyIDSet.ToSlice() {
				strategyID := fmt.Sprintf("%v", liveStrategyID)
				s, err := p.strategyManager.GetStrategyByID(ctx, strategyID)
				if err != nil {
					log.Println("strategyManager.GetStrategyByID fail")
					log.Println(err)
					continue
				}
				err = s.HandleKline(klines, kline)
				if err != nil {
					log.Println("s.HandleKline fail")
					log.Println(err)
					continue
				}
			}
			p.mux.Unlock()
		}
	}
}

// 取得最新的數根k線
func (p *Platform) getLimitKlineHistory(ctx context.Context, exchangeName string, symbol string, timeframe string, limit int) ([]market.Kline, error) {
	ex, err := p.exchangeManager.GetExchange(ctx, exchangeName, p.userManager.GetAdminUserID())
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

// 停止正在運行的strategy，使其不會被呼叫
func (p *Platform) StopStrategy(ctx context.Context, strategyID string) error {
	info, err := p.strategyManager.GetStrategyInfo(ctx, strategyID)
	if err != nil {
		log.Println("strategyManager.GetStrategyInfo fail")
		return err
	}

	streamName, err := p.getStreamName(info.Exchange)
	if err != nil {
		log.Println("p.getStreamName fail")
		return err
	}

	streamParam := stream.KlineStreamParam{
		Name:      streamName,
		Symbol:    info.Symbol,
		Timeframe: info.Timeframe,
	}
	streamKey := stream.GenKlineStreamKey(streamParam)

	p.mux.Lock()
	defer p.mux.Unlock()

	if info.Status == strategy.Live {
		err = p.strategyManager.UpdateStatus(ctx, strategyID, strategy.Stop)
		if err != nil {
			log.Println("strategyManager.UpdateStatus fail")
			return err
		}
	}

	_, ok := p.liveStrategyID[streamKey]
	if ok {
		p.liveStrategyID[streamKey].Remove(strategyID)
		// TODO close stream when no strategy listen this stream
		// if p.liveStrategyID[streamKey].Cardinality() == 0 {
		// 	p.runningStream.Remove(streamKey)
		// }
	}

	return nil
}
