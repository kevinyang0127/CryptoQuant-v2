package exchange

import (
	"CryptoQuant-v2/user"
	"context"
	"fmt"
	"log"
)

type Manager struct {
	userManager       *user.Manager
	userExchangeCache map[string]Exchange
}

func NewExchangeManager(userManager *user.Manager) *Manager {
	return &Manager{
		userManager:       userManager,
		userExchangeCache: make(map[string]Exchange),
	}
}

func (m *Manager) GetExchange(ctx context.Context, exchangeName string, userID string) (Exchange, error) {
	key := m.genUserExchangeKey(exchangeName, userID)
	ex, ok := m.userExchangeCache[key]
	if ok {
		return ex, nil
	}

	user, err := m.userManager.GetUser(ctx, userID)
	if err != nil {
		log.Println("userManager.GetUser fail")
		return nil, err
	}

	switch GetExchangeName(exchangeName) {
	case BINANCE_FUTURE:
		bf, err := newBinanceFuture(ctx, user.BinanceApiKey, user.BinanceSecretKey)
		if err != nil {
			log.Println("newBinanceFuture fail")
			return nil, err
		}
		m.userExchangeCache[key] = bf
		return bf, nil
	default:
		return nil, fmt.Errorf("don't support exchange(%s)", exchangeName)
	}
}

func (m *Manager) genUserExchangeKey(exchangeName, userID string) string {
	return fmt.Sprintf("%s_%s", userID, exchangeName)
}

func (m *Manager) RemoveUserExchangeCache(ctx context.Context, userID string) {
	delete(m.userExchangeCache, userID)
}
