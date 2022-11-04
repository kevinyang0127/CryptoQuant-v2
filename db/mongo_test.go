package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMongo() {
	ctx := context.Background()

	mongoDB, disconnect, err := NewMongoDB()
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

	r, err := mongoDB.FindOne(ctx, "cryptoQuantDB", "test", bson.D{{"botID", "1234"}})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("r: ", r)

	r2, err := mongoDB.Find(ctx, "cryptoQuantDB", "test", bson.D{{"botID", "1234"}})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("r2: ", r2)
}
