package stream

import (
	"CryptoQuant-v2/indicator"
	"context"
	"fmt"
	"log"
	"sync"
)

var KlineStreamManager *Manager

func init() {
	KlineStreamManager = &Manager{
		streamMap:    make(map[string]Stream),
		subscribeMap: make(map[string][]*KlineStreamSubscriber),
	}
}

type StreamName string

const (
	BinanceFKlineStream    StreamName = "BinanceFKline"
	BinanceFUserDataStream StreamName = "BinanceFUserData"
)

func IsKlineStream(streamName StreamName) bool {
	switch streamName {
	case BinanceFKlineStream:
		return true
	}
	return false
}

type KlineStreamParam struct {
	Name      StreamName
	Symbol    string
	Timeframe string // for klineStream
}

func GenKlineStreamKey(param KlineStreamParam) string {
	return fmt.Sprintf("%s_%s_%s", string(param.Name), param.Symbol, param.Timeframe)
}

type KlineStreamSubscriber struct {
	SubscriberID string
	SubscriberCh chan<- indicator.Kline
}

// kline stream manager
type Manager struct {
	mux          sync.Mutex
	streamMap    map[string]Stream                   // key: StreamName_symbol_timeframe , Ex: BinanceFKline_ETHUSDT_15m
	subscribeMap map[string][]*KlineStreamSubscriber // key: StreamName_symbol_timeframe, value: Subscriber List
}

func (m *Manager) createStream(ctx context.Context, param KlineStreamParam) (Stream, error) {
	var newStream Stream
	switch param.Name {
	case BinanceFKlineStream:
		newStream = newBinanceFKlineStream(param)
	default:
		return nil, fmt.Errorf("create stream fail, can't find StreamName(%s) or param error", param.Name)
	}

	err := newStream.wsConnect(ctx, m.klineHandler)
	if err != nil {
		log.Println("newStream.wsConnect fail")
		log.Println(err)
		return nil, err
	}
	return newStream, nil
}

func (m *Manager) klineHandler(streamKey string, kline indicator.Kline) {
	m.mux.Lock()
	defer m.mux.Unlock()

	subscribeList := m.subscribeMap[streamKey]
	for _, s := range subscribeList {
		select {
		case s.SubscriberCh <- kline:
		default:
			log.Printf("publish to subscriber(key = %s) fail", streamKey)
		}
	}
}

func (m *Manager) Subscribe(ctx context.Context, param KlineStreamParam, subscriberID string, subscriberCh chan<- indicator.Kline) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	key := GenKlineStreamKey(param)

	_, ok := m.streamMap[key]
	if !ok {
		s, err := m.createStream(ctx, param)
		if err != nil {
			log.Println("create stream fail")
			return err
		}
		m.streamMap[key] = s
		m.subscribeMap[key] = []*KlineStreamSubscriber{}
	}

	subscribeList, ok := m.subscribeMap[key]
	if !ok {
		return fmt.Errorf("can't find subscribe list")
	}

	for _, s := range subscribeList {
		if s.SubscriberID == subscriberID {
			log.Println("subscriber already subscribe")
			return nil
		}
	}
	m.subscribeMap[key] = append(subscribeList, &KlineStreamSubscriber{
		SubscriberID: subscriberID,
		SubscriberCh: subscriberCh,
	})

	return nil
}

func (m *Manager) Unsubscribe(ctx context.Context, param KlineStreamParam, subscriberID string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	key := GenKlineStreamKey(param)
	subscribeList, ok := m.subscribeMap[key]
	if !ok {
		return fmt.Errorf("can't find subscribe list")
	}

	for i, s := range subscribeList {
		if s.SubscriberID == subscriberID {
			close(s.SubscriberCh)
			m.subscribeMap[key] = append(subscribeList[:i], subscribeList[i+1:]...)
			return nil
		}
	}

	log.Println("subscriber didn't subscribe stream")
	return nil
}
