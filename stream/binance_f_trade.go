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

func newBinanceFTradeStream(symbol string) *BinanceFTrade {
	return &BinanceFTrade{
		symbol:        symbol,
		SubscriberMap: make(map[string]chan<- []byte),
	}
}

type BinanceFTrade struct {
	mux           sync.Mutex
	SubscriberMap map[string]chan<- []byte
	symbol        string
}

func (bft *BinanceFTrade) subscribe(ctx context.Context, subscriberKey string, subscriberCh chan<- []byte) error {
	bft.mux.Lock()
	defer bft.mux.Unlock()

	_, ok := bft.SubscriberMap[subscriberKey]
	if ok {
		return fmt.Errorf("subscriberKey: %s already subscribe", subscriberKey)
	}
	bft.SubscriberMap[subscriberKey] = subscriberCh
	return nil
}

func (bft *BinanceFTrade) unsubscribe(ctx context.Context, subscriberKey string) {
	bft.mux.Lock()
	defer bft.mux.Unlock()
	close(bft.SubscriberMap[subscriberKey])
	delete(bft.SubscriberMap, subscriberKey)
}

func (bft *BinanceFTrade) publish(ctx context.Context, data []byte) {
	bft.mux.Lock()
	defer bft.mux.Unlock()
	for key := range bft.SubscriberMap {
		ch := bft.SubscriberMap[key]
		select {
		case ch <- data:
		default:
			log.Printf("publish to subscriber(key = %s) fail", key)
		}
	}
}

func (bft *BinanceFTrade) wsConnect(ctx context.Context) error {
	wsTradeHandler := func(event *binanceFutures.WsAggTradeEvent) {
		data, err := json.Marshal(event)
		if err != nil {
			log.Println("json.Marshal() fail")
			return
		}
		bft.publish(ctx, data)
	}

	errHandler := func(err error) {
		log.Println(err)
	}

	doneC, _, err := binanceFutures.WsAggTradeServe(bft.symbol, wsTradeHandler, errHandler)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		<-doneC
		log.Println("binanceFutures.WsAggTradeServe is closed")
		bft.closeStream(ctx)
	}()

	return nil
}

func (s *BinanceFTrade) closeStream(ctx context.Context) {
	log.Println("(BinanceFTrade) start unsubscribe all subscriber")
	for key := range s.SubscriberMap {
		s.unsubscribe(ctx, key)
	}
	log.Println("(BinanceFTrade) unsubscribe all subscriber done")
}
