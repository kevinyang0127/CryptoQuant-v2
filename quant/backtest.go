package quant

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/exchange"
	"CryptoQuant-v2/indicator"
	"CryptoQuant-v2/simulation"
	"CryptoQuant-v2/strategy"
	"context"
	"fmt"
	"log"
	"time"
)

type BacktestingClient struct {
	userID              string
	strategyID          string
	exchange            string
	symbol              string
	timeframe           string
	startBalance        string
	lever               int
	takerCommissionRate string
	makerCommissionRate string
	startTimeMs         int64
	endTimeMs           int64
	strategyManager     *strategy.Manager
}

func NewBackTestingClient(mongoDB *db.MongoDB, userID string, strategyID string, exchange string, symbol string,
	timeframe string, startBalance string, lever int, takerCommissionRate string, makerCommissionRate string,
	startTimeMs int64, endTimeMs int64) *BacktestingClient {
	return &BacktestingClient{
		userID:              userID,
		strategyID:          strategyID,
		exchange:            exchange,
		symbol:              symbol,
		timeframe:           timeframe,
		startBalance:        startBalance,
		lever:               lever,
		takerCommissionRate: takerCommissionRate,
		makerCommissionRate: makerCommissionRate,
		startTimeMs:         startTimeMs,
		endTimeMs:           endTimeMs,
		strategyManager:     strategy.NewManager(mongoDB),
	}
}

func (b *BacktestingClient) Backtest(ctx context.Context) (simulationID string, err error) {

	if b.startTimeMs >= b.endTimeMs {
		return "", fmt.Errorf("startTimeMs >= endTimeMs")
	}

	simulationKlineCh := make(chan indicator.Kline)
	simulationID, err = simulation.SimulationManager.StartNewSimulation(ctx, simulationKlineCh, b.userID,
		b.startBalance, int64(b.lever), b.takerCommissionRate, b.makerCommissionRate)
	if err != nil {
		log.Println("SimulationManager.StartNewSimulation fail")
		return "", err
	}

	s, err := b.strategyManager.GetStrategyByID(ctx, b.strategyID)
	if err != nil {
		log.Println("strategyManager.GetStrategyByID fail")
		return "", err
	}

	err = b.runBacktesting(ctx, simulationKlineCh, simulationID, s)
	if err != nil {
		log.Println("RunBacktesting fail")
		return "", err
	}

	return simulationID, nil
}

func (b *BacktestingClient) runBacktesting(ctx context.Context, simulationKlineCh chan indicator.Kline, simulationID string, s strategy.Strategy) error {
	ex, err := exchange.GetExchange(b.exchange)
	if err != nil {
		log.Println("GetExchange fail")
		return err
	}
	unitTime, err := b.getTimeframeUnitTime()
	if err != nil {
		log.Println("getTimeframeUnitTime fail")
		return err
	}
	// startTime前的500根
	beforeKlines, err := ex.GetLimitKlineHistoryByTime(ctx, b.symbol, b.timeframe, 500, b.startTimeMs-(500*unitTime).Milliseconds(), b.startTimeMs)
	if err != nil {
		log.Println("GetLimitKlineHistoryByTime fail")
		return err
	}

	go func() {
		durationMs := b.endTimeMs - b.startTimeMs
		// 總共需要多少根k線資料
		count := durationMs / unitTime.Milliseconds()

		sleepCount := 0
		startTimeMs := int64(0)
		endTimeMs := int64(0)
		maxKlineOnce := 1000
		//每次回測1000筆k線，每2萬筆休息30秒
		for i := 0; i < int(count); i += maxKlineOnce {
			if sleepCount >= 20 {
				time.Sleep(30 * time.Second)
				sleepCount = 0
			}
			sleepCount++

			// 超過2000筆只保留最新500筆
			if len(beforeKlines) >= 2000 {
				newSlice := make([]indicator.Kline, 500, 2000)
				copy(newSlice, beforeKlines[len(beforeKlines)-500:])
				beforeKlines = newSlice
			}

			if startTimeMs == 0 {
				startTimeMs = b.startTimeMs
			}
			endTimeMs = startTimeMs + int64(maxKlineOnce)*unitTime.Milliseconds()
			if endTimeMs > b.endTimeMs {
				endTimeMs = b.endTimeMs
			}

			klines, err := ex.GetLimitKlineHistoryByTime(ctx, b.symbol, b.timeframe, maxKlineOnce, startTimeMs, endTimeMs)
			if err != nil {
				log.Println("GetLimitKlineHistoryByTime fail")
				return
			}
			for _, kline := range klines {
				if kline.EndTime < beforeKlines[len(beforeKlines)-1].EndTime {
					continue
				}
				simulationKlineCh <- kline
				s.HandleBackTestKline(simulationID, beforeKlines, kline)
				beforeKlines = append(beforeKlines, kline)
			}

			startTimeMs = endTimeMs
		}

		simulation.SimulationManager.StopSimulation(ctx, simulationID)
	}()

	return nil
}

func (b *BacktestingClient) getTimeframeUnitTime() (time.Duration, error) {
	switch b.timeframe {
	case "1s":
		return time.Second, nil
	case "1m":
		return time.Minute, nil
	case "3m":
		return 3 * time.Minute, nil
	case "5m":
		return 5 * time.Minute, nil
	case "15m":
		return 15 * time.Minute, nil
	case "30m":
		return 30 * time.Minute, nil
	case "1h":
		return time.Hour, nil
	case "2h":
		return 2 * time.Hour, nil
	case "4h":
		return 4 * time.Hour, nil
	case "6h":
		return 6 * time.Hour, nil
	case "8h":
		return 8 * time.Hour, nil
	case "12h":
		return 12 * time.Hour, nil
	case "1d":
		return 24 * time.Hour, nil
	case "3d":
		return 3 * 24 * time.Hour, nil
	case "1w":
		return 7 * 24 * time.Hour, nil
	case "1M":
		return 30 * 24 * time.Hour, nil
	default:
		return time.Duration(0), fmt.Errorf("not support timeframe")
	}
}
