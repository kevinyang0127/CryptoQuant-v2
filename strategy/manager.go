package strategy

import (
	"CryptoQuant-v2/db"
	"CryptoQuant-v2/script"
	"CryptoQuant-v2/util"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StrategyStatus int

const (
	Draft StrategyStatus = iota
	Live
	Stop
)

type StrategyInfo struct {
	UserID       string         `bson:"userID"`
	Exchange     string         `bson:"exchange"`
	Symbol       string         `bson:"symbol"`
	Timeframe    string         `bson:"timeframe"`
	Status       StrategyStatus `bson:"status"`
	StrategyID   string         `bson:"strategyID"`
	StrategyName string         `bson:"strategyName"`
	Script       string         `bson:"script"`
}

type Manager struct {
	mongoDB          *db.MongoDB
	strategyMap      map[string]Strategy // key: strategyID
	luaScriptHandler *script.LuaScriptHandler
}

func NewManager(mongoDB *db.MongoDB, luaScriptHandler *script.LuaScriptHandler) *Manager {
	return &Manager{
		mongoDB:          mongoDB,
		strategyMap:      make(map[string]Strategy),
		luaScriptHandler: luaScriptHandler,
	}
}

func (m *Manager) Add(ctx context.Context, userID string, exchange string, symbol string, timeframe string,
	status StrategyStatus, strategyName string, script string) (strategyID string, err error) {

	strategyID = util.GenIDWithPrefix("S_", 10)
	s := &StrategyInfo{
		UserID:       userID,
		Exchange:     exchange,
		Symbol:       symbol,
		Timeframe:    timeframe,
		Status:       status,
		StrategyID:   strategyID,
		StrategyName: strategyName,
		Script:       script,
	}
	_, err = m.mongoDB.InsertOne(ctx, "cryptoQuantV2", "strategy", s)
	if err != nil {
		log.Println("mongoDB.InsertOne fail")
		return "", err
	}

	return strategyID, nil
}

func (m *Manager) GetStrategyByID(ctx context.Context, strategyID string) (Strategy, error) {
	s, ok := m.strategyMap[strategyID]
	if ok {
		return s, nil
	}

	info := &StrategyInfo{}
	err := m.mongoDB.FindOne(ctx, "cryptoQuantV2", "strategy", bson.D{{"strategyID", strategyID}}, info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("FindOne no result")
			log.Println("can't find Strategy by StrategyID = " + strategyID)
		} else {
			log.Println("FindOne fail")
		}
		return nil, err
	}

	s = m.getDefaultStrategy(info)

	// save in cache
	m.strategyMap[strategyID] = s

	return s, nil
}

func (m *Manager) getDefaultStrategy(strategyInfo *StrategyInfo) Strategy {
	return NewLuaScriptStrategy(m.luaScriptHandler, strategyInfo)
}

func (m *Manager) GetStrategyByUserIDAndName(ctx context.Context, userID string, Name string) (Strategy, error) {
	return nil, nil
}

func (m *Manager) GetStrategyInfoByUserID(ctx context.Context, userID string) ([]*StrategyInfo, error) {
	result := []*StrategyInfo{}
	err := m.mongoDB.FindAll(ctx, "cryptoQuantV2", "strategy", bson.D{{"userID", userID}}, &result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("mongoDB.FindAll no result")
			log.Println("can't find StrategyInfo by userID = " + userID)
		} else {
			log.Println("mongoDB.FindAll() fail")
		}
		return nil, err
	}

	return result, nil
}

func (m *Manager) GetStrategyInfo(ctx context.Context, strategyID string) (*StrategyInfo, error) {
	info := &StrategyInfo{}
	err := m.mongoDB.FindOne(ctx, "cryptoQuantV2", "strategy", bson.D{{"strategyID", strategyID}}, info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("FindOne no result")
			log.Println("can't find Strategy by StrategyID = " + strategyID)
		} else {
			log.Println("FindOne fail")
		}
		return nil, err
	}
	return info, nil
}

type StrategyInfoUpdater struct {
	Exchange     string
	Symbol       string
	Timeframe    string
	StrategyName string
	Script       string
}

func (m *Manager) UpdateStrategyInfo(ctx context.Context, strategyID string, strategyInfoUpdater *StrategyInfoUpdater) error {
	info := &StrategyInfo{}
	err := m.mongoDB.FindOne(ctx, "cryptoQuantV2", "strategy", bson.D{{"strategyID", strategyID}}, info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("FindOne no result")
			log.Println("can't find Strategy by StrategyID = " + strategyID)
		} else {
			log.Println("FindOne fail")
		}
		return err
	}

	if info.Status == Live {
		return fmt.Errorf("can't update live strategy")
	}

	updater := bson.D{}
	if strategyInfoUpdater.Exchange != "" {
		updater = append(updater, bson.E{"exchange", strategyInfoUpdater.Exchange})
	}
	if strategyInfoUpdater.Symbol != "" {
		updater = append(updater, bson.E{"symbol", strategyInfoUpdater.Symbol})
	}
	if strategyInfoUpdater.Timeframe != "" {
		updater = append(updater, bson.E{"timeframe", strategyInfoUpdater.Timeframe})
	}
	if strategyInfoUpdater.StrategyName != "" {
		updater = append(updater, bson.E{"strategyName", strategyInfoUpdater.StrategyName})
	}
	if strategyInfoUpdater.Script != "" {
		updater = append(updater, bson.E{"script", strategyInfoUpdater.Script})
	}

	err = m.mongoDB.UpdateOne(ctx, "cryptoQuantV2", "strategy", bson.D{{"strategyID", strategyID}}, updater)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("mongoDB.UpdateOne no result")
			log.Println("can't find StrategyInfo by strategyID = " + strategyID)
		} else {
			log.Println("mongoDB.UpdateOne() fail")
		}
		return err
	}

	return nil
}

func (m *Manager) UpdateStatus(ctx context.Context, strategyID string, status StrategyStatus) error {
	err := m.mongoDB.UpdateOne(ctx, "cryptoQuantV2", "strategy", bson.D{{"strategyID", strategyID}}, bson.D{{"status", status}})
	if err != nil {
		log.Println("mongoDB.UpdateOne fail")
		return err
	}
	return nil
}

func (m *Manager) DeleteStrategy(ctx context.Context, strategyID string) error {
	info := &StrategyInfo{}
	err := m.mongoDB.FindOne(ctx, "cryptoQuantV2", "strategy", bson.D{{"strategyID", strategyID}}, info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("FindOne no result")
			log.Println("can't find Strategy by StrategyID = " + strategyID)
		} else {
			log.Println("FindOne fail")
		}
		return err
	}

	if info.Status == Live {
		return fmt.Errorf("can't delete live strategy")
	}

	err = m.mongoDB.DeleteOne(ctx, "cryptoQuantV2", "strategy", bson.D{{"strategyID", strategyID}})
	if err != nil {
		log.Println("mongoDB.DeleteOne() fail")
		return err
	}

	return nil
}
