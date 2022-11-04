package stream

/*
	資料推送源：幣安U本位合約-歸集交易
	https://binance-docs.github.io/apidocs/futures/cn/#1e66c0284e
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	binanceFutures "github.com/adshao/go-binance/v2/futures"
)

func newBinanceFKlineStream(symbol string, interval string) *BinanceFKline {
	return &BinanceFKline{
		symbol:        symbol,
		interval:      interval,
		SubscriberMap: make(map[string]chan<- []byte),
	}
}

type BinanceFKline struct {
	mux           sync.Mutex
	SubscriberMap map[string]chan<- []byte
	symbol        string // ETHUSDT BTCUSDT ...
	interval      string // 1m 5m 1h 4h 1d ...
}

func (s *BinanceFKline) subscribe(ctx context.Context, subscriberKey string, subscriberCh chan<- []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.SubscriberMap[subscriberKey]
	if ok {
		return fmt.Errorf("subscriberKey: %s already subscribe", subscriberKey)
	}
	s.SubscriberMap[subscriberKey] = subscriberCh
	return nil
}

func (s *BinanceFKline) unsubscribe(ctx context.Context, subscriberKey string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	close(s.SubscriberMap[subscriberKey])
	delete(s.SubscriberMap, subscriberKey)
}

func (s *BinanceFKline) publish(ctx context.Context, data []byte) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for key := range s.SubscriberMap {
		ch := s.SubscriberMap[key]
		select {
		case ch <- data:
		default:
			log.Printf("publish to subscriber(key = %s) fail", key)
		}
	}
}

func (s *BinanceFKline) wsConnect(ctx context.Context) error {
	wsKlineHandler := func(event *binanceFutures.WsKlineEvent) {
		data, err := json.Marshal(event)
		if err != nil {
			log.Println("json.Marshal() fail")
			return
		}
		s.publish(ctx, data)
	}

	errHandler := func(err error) {
		log.Println(err)
	}

	doneC, _, err := binanceFutures.WsKlineServe(s.symbol, s.interval, wsKlineHandler, errHandler)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		<-doneC
		log.Println("binanceFutures.WsUserDataServe is closed")
		s.closeStream(ctx)
	}()

	return nil
}

func (s *BinanceFKline) closeStream(ctx context.Context) {
	log.Println("(BinanceFKline) start unsubscribe all subscriber")
	for key := range s.SubscriberMap {
		s.unsubscribe(ctx, key)
	}
	log.Println("(BinanceFKline) unsubscribe all subscriber done")
}
