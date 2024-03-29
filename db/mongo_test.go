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
	log.Println("before update result: ", result)

	err = mongoDB.UpdateOne(ctx, "cryptoQuantDB", "test", bson.D{{"botID", "1234"}}, bson.D{{"balance", 90.0}})
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	err = mongoDB.FindOne(ctx, "cryptoQuantDB", "test", bson.D{{"botID", "1234"}}, result)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}
	log.Println("update result: ", result)

	findAllResults := []*TradeLog{}
	err = mongoDB.FindAll(ctx, "cryptoQuantDB", "test", bson.D{{"botID", "12345"}}, &findAllResults)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}
	log.Println("findAllResults: ", findAllResults)

	t.Error("123456")
}
