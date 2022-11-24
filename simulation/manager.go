package simulation

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/market"
	"context"
	"fmt"
	"log"

	"github.com/shopspring/decimal"
)

type Manager struct {
	mongoDB             *db.MongoDB
	simulationMap       map[string]*Simulation
	cancelSimulationMap map[string]context.CancelFunc
}

func NewSimulationManager(mongoDB *db.MongoDB) *Manager {
	return &Manager{
		mongoDB:             mongoDB,
		simulationMap:       make(map[string]*Simulation),
		cancelSimulationMap: make(map[string]context.CancelFunc),
	}
}

func (m *Manager) StartNewSimulation(ctx context.Context, ch chan market.Kline, userID string, startBalance string,
	lever int64, takerCommissionRate string, makerCommissionRate string) (simulationID string, err error) {
	startBalanceD, err := decimal.NewFromString(startBalance)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return "", err
	}
	leverD := decimal.NewFromInt(lever)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return "", err
	}
	takerCommissionRateD, err := decimal.NewFromString(takerCommissionRate)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return "", err
	}
	makerCommissionRateD, err := decimal.NewFromString(makerCommissionRate)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return "", err
	}

	ctx, cancel := context.WithCancel(ctx)
	s := NewSimulation(m.mongoDB, userID, startBalanceD, leverD, takerCommissionRateD, makerCommissionRateD)
	go s.ListenNewKline(ctx, ch)
	m.cancelSimulationMap[s.simulationID] = cancel
	m.simulationMap[s.simulationID] = s

	return s.simulationID, nil
}

func (m *Manager) StopSimulation(ctx context.Context, simulationID string) {
	cancelFunc, ok := m.cancelSimulationMap[simulationID]
	if !ok {
		log.Println("can't find simulation cancel func")
	}
	cancelFunc()
	delete(m.cancelSimulationMap, simulationID)
	log.Println("StopSimulation success")
}

func (m *Manager) Entry(ctx context.Context, simulationID string, side bool, price string, quantity string, isMaker bool, klineTimestamp int64) error {
	priceD, err := decimal.NewFromString(price)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return err
	}
	quantityD, err := decimal.NewFromString(quantity)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return err
	}

	if quantityD.IsNegative() {
		log.Println("input quantity IsNegative")
		return fmt.Errorf("input quantity IsNegative")
	}

	if !side {
		quantityD = quantityD.Mul(decimal.NewFromInt(-1))
	}

	s, ok := m.simulationMap[simulationID]
	if !ok {
		return fmt.Errorf("simulationMap can't find simulationID = %s", simulationID)
	}
	s.Entry(ctx, priceD, quantityD, isMaker, klineTimestamp)
	return nil
}

func (m *Manager) Exit(ctx context.Context, simulationID string, price string, quantity string, isMaker bool, klineTimestamp int64) error {
	priceD, err := decimal.NewFromString(price)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return err
	}
	quantityD, err := decimal.NewFromString(quantity)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return err
	}

	s, ok := m.simulationMap[simulationID]
	if !ok {
		return fmt.Errorf("simulationMap can't find simulationID = %s", simulationID)
	}
	s.Exit(ctx, priceD, quantityD, isMaker, klineTimestamp)
	return nil
}

func (m *Manager) ExitAll(ctx context.Context, simulationID string, price string, isMaker bool, klineTimestamp int64) error {
	priceD, err := decimal.NewFromString(price)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return err
	}

	s, ok := m.simulationMap[simulationID]
	if !ok {
		return fmt.Errorf("simulationMap can't find simulationID = %s", simulationID)
	}

	p := s.GetPosition(ctx)
	if p != nil {
		s.Exit(ctx, priceD, p.Quantity.Abs(), isMaker, klineTimestamp)
	}

	return nil
}

func (m *Manager) Order(ctx context.Context, simulationID string, side bool, price string, quantity string) error {
	priceD, err := decimal.NewFromString(price)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return err
	}
	quantityD, err := decimal.NewFromString(quantity)
	if err != nil {
		log.Println("decimal.NewFromString error")
		return err
	}

	if quantityD.IsNegative() {
		log.Println("input quantity IsNegative")
		return fmt.Errorf("input quantity IsNegative")
	}

	if !side {
		quantityD = quantityD.Mul(decimal.NewFromInt(-1))
	}

	s, ok := m.simulationMap[simulationID]
	if !ok {
		return fmt.Errorf("simulationMap can't find simulationID = %s", simulationID)
	}
	s.Order(ctx, priceD, quantityD)
	return nil
}

func (m *Manager) CloseAllOrder(ctx context.Context, simulationID string) error {
	s, ok := m.simulationMap[simulationID]
	if !ok {
		return fmt.Errorf("simulationMap can't find simulationID = %s", simulationID)
	}
	s.CloseAllOrder(ctx)
	return nil
}

func (m *Manager) GetAllOrder(ctx context.Context, simulationID string) ([]*Order, error) {
	s, ok := m.simulationMap[simulationID]
	if !ok {
		return nil, fmt.Errorf("simulationMap can't find simulationID = %s", simulationID)
	}
	return s.GetAllOrder(ctx), nil
}

func (m *Manager) GetPosition(ctx context.Context, simulationID string) (*Position, error) {
	s, ok := m.simulationMap[simulationID]
	if !ok {
		return nil, fmt.Errorf("simulationMap can't find simulationID = %s", simulationID)
	}
	return s.GetPosition(ctx), nil
}
