package db

import (
	"context"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMongo(t *testing.T) {
	ctx := context.Background()

	mongoDB, disconnect, err := NewMongoDB(LocalURI)
	if err != nil {
		log.Println(err)
		return
	}
	defer disconnect()

	type TradeLog struct {
		BotID     string  `bson:"botID"`
		OpenClose bool    `bson:"openClose"`
		Side      bool    `bson:"side"`
		Timestamp int64   `bson:"timestamp"`
		Balance   float64 `bson:"balance"`
	}

	tradeLog := TradeLog{
		BotID:     "1234",
		OpenClose: true,
		Side:      true,
		Timestamp: time.Now().Unix(),
		Balance:   10000.0,
	}

	_, err = mongoDB.InsertOne(ctx, "cryptoQuantDB", "test", tradeLog)
	if err != nil {
		log.Println(err)
		return
	}

	result := &TradeLog{}
	err = mongoDB.FindOne(ctx, "cryptoQuantDB", "test", bson.D{{"botID", "1234"}}, result)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}
	log.Println("result: ", result)

	// r2, err := mongoDB.Find(ctx, "cryptoQuantDB", "test", bson.D{{"botID", "1234"}})
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println("r2: ", r2)

	t.Error("123")
}
