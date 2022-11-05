package stream

/*
	資料推送源：幣安U本位合約-帳戶信息推送
	https://binance-docs.github.io/apidocs/futures/cn/#balance-position
	https://binance-docs.github.io/apidocs/futures/cn/#060a012f0b
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	binanceFutures "github.com/adshao/go-binance/v2/futures"
)

func newBinanceFUserDataStream(listenKey string) *BinanceFUserData {
	return &BinanceFUserData{
		listenKey:     listenKey,
		SubscriberMap: make(map[string]chan<- []byte),
	}
}

type BinanceFUserData struct {
	mux           sync.Mutex
	SubscriberMap map[string]chan<- []byte
	listenKey     string
}

func (s *BinanceFUserData) subscribe(ctx context.Context, subscriberKey string, subscriberCh chan<- []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.SubscriberMap[subscriberKey]
	if ok {
		return fmt.Errorf("subscriberKey: %s already subscribe", subscriberKey)
	}
	s.SubscriberMap[subscriberKey] = subscriberCh
	return nil
}

func (s *BinanceFUserData) unsubscribe(ctx context.Context, subscriberKey string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	close(s.SubscriberMap[subscriberKey])
	delete(s.SubscriberMap, subscriberKey)
}

func (s *BinanceFUserData) publish(ctx context.Context, data []byte) {
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

func (s *BinanceFUserData) wsConnect(ctx context.Context) error {
	wsUserDataHandler := func(event *binanceFutures.WsUserDataEvent) {
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

	doneC, _, err := binanceFutures.WsUserDataServe(s.listenKey, wsUserDataHandler, errHandler)
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

func (s *BinanceFUserData) closeStream(ctx context.Context) {
	log.Println("(BinanceFUserData) start unsubscribe all subscriber")
	for key := range s.SubscriberMap {
		s.unsubscribe(ctx, key)
	}
	log.Println("(BinanceFUserData) unsubscribe all subscriber done")
}
