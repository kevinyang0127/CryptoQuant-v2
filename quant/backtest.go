package quant

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/exchange"
	"CryptoQuant-v2/market"
	"CryptoQuant-v2/simulation"
	"CryptoQuant-v2/strategy"
	"CryptoQuant-v2/user"
	"context"
	"fmt"
	"log"
	"time"
)

type BacktestingClient struct {
	userID                string
	strategyID            string
	exchangeName          string
	symbol                string
	timeframe             string
	klineHistoryTimeframe string
	startBalance          string
	lever                 int
	takerCommissionRate   string
	makerCommissionRate   string
	startTimeMs           int64
	endTimeMs             int64
	strategyManager       *strategy.Manager
	exchangeManager       *exchange.Manager
	userManager           *user.Manager
}

func NewBackTestingClient(mongoDB *db.MongoDB, userID string, strategyID string, exchangeName string, symbol string,
	timeframe string, klineHistoryTimeframe string, startBalance string, lever int, takerCommissionRate string, makerCommissionRate string,
	startTimeMs int64, endTimeMs int64) *BacktestingClient {

	userManager := user.NewUserManager(mongoDB)

	return &BacktestingClient{
		userID:                userID,
		strategyID:            strategyID,
		exchangeName:          exchangeName,
		symbol:                symbol,
		timeframe:             timeframe,
		klineHistoryTimeframe: klineHistoryTimeframe,
		startBalance:          startBalance,
		lever:                 lever,
		takerCommissionRate:   takerCommissionRate,
		makerCommissionRate:   makerCommissionRate,
		startTimeMs:           startTimeMs,
		endTimeMs:             endTimeMs,
		strategyManager:       strategy.NewManager(mongoDB),
		exchangeManager:       exchange.NewExchangeManager(userManager),
		userManager:           userManager,
	}
}

func (b *BacktestingClient) Backtest(ctx context.Context) (simulationID string, err error) {

	if b.startTimeMs >= b.endTimeMs {
		return "", fmt.Errorf("startTimeMs >= endTimeMs")
	}

	simulationKlineCh := make(chan market.Kline)
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

func (b *BacktestingClient) runBacktesting(ctx context.Context, simulationKlineCh chan market.Kline, simulationID string, s strategy.Strategy) error {
	ex, err := b.exchangeManager.GetExchange(ctx, b.exchangeName, b.userID)
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

		startTimeMs := int64(0)
		endTimeMs := int64(0)
		maxKlineOnce := 1000
		apiRequestTimes := 0
		//每次回測1000筆k線
		for i := 0; i < int(count); i += maxKlineOnce {
			if apiRequestTimes >= 1000 {
				time.Sleep(55 * time.Second)
				apiRequestTimes = 0
			}

			// 超過2000筆只保留最新500筆
			if len(beforeKlines) >= 2000 {
				newSlice := make([]market.Kline, 500, 2000)
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
				log.Println(err)
				return
			}
			apiRequestTimes++

			for _, kline := range klines {
				if kline.EndTime < beforeKlines[len(beforeKlines)-1].EndTime {
					continue
				}

				// 拿到此根k線的更小時間範圍的數據當作k線走過的歷史
				smallTimeframeKlines, err := ex.GetLimitKlineHistoryByTime(ctx, b.symbol, b.klineHistoryTimeframe, maxKlineOnce, kline.StartTime, kline.EndTime)
				if err != nil {
					log.Println("ex.GetLimitKlineHistoryByTime fail")
					log.Println(err)
					return
				}
				apiRequestTimes++

				klineHistory, err := market.GenFinalKlineHistory(kline, smallTimeframeKlines)
				if err != nil {
					log.Println("market.GenFinalKlineHistory fail")
					log.Println(err)
					return
				}

				for _, kh := range klineHistory {
					simulationKlineCh <- kh
					if kh.IsFinal {
						beforeKlines = append(beforeKlines, kline) // beforeKlines包含目前收盤kline
					}
					s.HandleBackTestKline(simulationID, beforeKlines, kh)
				}
				fmt.Println(kline)
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
