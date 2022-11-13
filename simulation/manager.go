package simulation

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/indicator"
	"context"
	"fmt"
	"log"

	"github.com/shopspring/decimal"
)

var SimulationManager *Manager

type Manager struct {
	mongoDB             *db.MongoDB
	simulationMap       map[string]*Simulation
	cancelSimulationMap map[string]context.CancelFunc
}

func RunNewManager(mongoDB *db.MongoDB) *Manager {
	if SimulationManager != nil {
		return SimulationManager
	}

	SimulationManager = &Manager{
		mongoDB:             mongoDB,
		simulationMap:       make(map[string]*Simulation),
		cancelSimulationMap: make(map[string]context.CancelFunc),
	}

	return SimulationManager
}

func (m *Manager) StartNewSimulation(ctx context.Context, ch chan indicator.Kline, userID string, startBalance string,
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
