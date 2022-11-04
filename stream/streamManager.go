package stream

import (
	"context"
	"fmt"
	"log"
)

var StreamManager *Manager

func init() {
	StreamManager = &Manager{
		StreamMap: make(map[string]Stream),
	}
}

type StreamName string

const (
	BinanceFTradeStream    StreamName = "BinanceFTrade"
	BinanceFKlineStream    StreamName = "BinanceFKline"
	BinanceFUserDataStream StreamName = "BinanceFUserData"
)

type StreamParam struct {
	Name      StreamName
	Symbol    string
	Interval  string
	ListenKey string
}

type Manager struct {
	StreamMap map[string]Stream // key = StreamName_symbol , Ex: BinanceFTrade_ETHUSDT
}

func (m *Manager) getKey(name StreamName, symbol string) string {
	return fmt.Sprintf("%s_%s", string(name), symbol)
}

func (m *Manager) createStream(ctx context.Context, param StreamParam) (Stream, error) {
	var newStream Stream
	switch param.Name {
	case BinanceFTradeStream:
		newStream = newBinanceFTradeStream(param.Symbol)
	case BinanceFKlineStream:
		newStream = newBinanceFKlineStream(param.Symbol, param.Interval)
	case BinanceFUserDataStream:
		newStream = newBinanceFUserDataStream(param.ListenKey)
	default:
		newStream = nil
	}
	if newStream != nil {
		err := newStream.wsConnect(ctx)
		if err != nil {
			//ctx.Err()
			return nil, err
		}
		return newStream, nil
	}
	return nil, fmt.Errorf("create stream fail, can't find StreamName(%s)", param.Name)
}

// 若不需要 interval 則帶入""即可
func (m *Manager) getStream(ctx context.Context, param StreamParam) (Stream, error) {
	key := m.getKey(param.Name, param.Symbol)
	stream, ok := m.StreamMap[key]
	if !ok {
		stream, err := m.createStream(ctx, param)
		if err != nil {
			log.Println("m.createStream fail")
			return nil, err
		}
		m.StreamMap[key] = stream
		return stream, nil
	}
	return stream, nil
}

func (m *Manager) SubscribeStream(ctx context.Context, param StreamParam, subscriberKey string, subscriberCh chan<- []byte) error {
	s, err := m.getStream(ctx, param)
	if err != nil {
		log.Println("m.getStream fail")
		return err
	}
	err = s.subscribe(ctx, subscriberKey, subscriberCh)
	if err != nil {
		log.Println("s.Subscribe fail")
		return err
	}
	return nil
}

func (m *Manager) UnsubscribeStream(ctx context.Context, param StreamParam, subscriberKey string) error {
	s, err := m.getStream(ctx, param)
	if err != nil {
		log.Println("m.getStream fail")
		return err
	}
	s.unsubscribe(ctx, subscriberKey)
	return nil
}
