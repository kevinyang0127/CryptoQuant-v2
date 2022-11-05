package stream

/*
	資料推送源：幣安U本位合約-k線
	https://binance-docs.github.io/apidocs/futures/cn/#k-6
*/

import (
	"CryptoQuant-v2/indicator"
	"context"
	"log"

	binanceFutures "github.com/adshao/go-binance/v2/futures"
)

func newBinanceFKlineStream(param KlineStreamParam) *BinanceFKline {
	return &BinanceFKline{
		Param: param,
	}
}

type BinanceFKline struct {
	Param KlineStreamParam
}

func (s *BinanceFKline) wsConnect(ctx context.Context, klineHandler func(streamKey string, kline indicator.Kline)) error {
	wsKlineHandler := func(event *binanceFutures.WsKlineEvent) {
		kline, err := indicator.BinanceFKlineEventToKline(*event)
		if err != nil {
			log.Println("indicator.BinanceFKlineEventToKline fail")
			return
		}
		key := GenKlineStreamKey(s.Param)
		klineHandler(key, *kline)
	}

	errHandler := func(err error) {
		log.Println(err)
	}

	doneC, _, err := binanceFutures.WsKlineServe(s.Param.Symbol, s.Param.Timeframe, wsKlineHandler, errHandler)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		<-doneC
		log.Println("binanceFutures.WsKlineServe is closed")
		//s.closeStream(ctx)
	}()

	return nil
}
