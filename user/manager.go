package user

import (
	"CryptoQuant-v2/db"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "user"
)

type Manager struct {
	mongoDB *db.MongoDB
}

func NewUserManager(mongoDB *db.MongoDB) *Manager {
	return &Manager{
		mongoDB: mongoDB,
	}
}

func (m *Manager) Register(ctx context.Context, userID string, name string) error {
	user := &User{}
	err := m.mongoDB.FindOne(ctx, db.DBNAME, collectionName, bson.D{{"userID", userID}}, user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = m.mongoDB.InsertOne(ctx, db.DBNAME, collectionName, &User{
				UserID: userID,
				Name:   name,
			})
			if err != nil {
				log.Println("mongoDB.InsertOne fail")
				return err
			}
			return nil
		}
		log.Println("mongoDB.FindOne fail")
		return err
	}

	return fmt.Errorf("userID already exist")
}

func (m *Manager) GetUser(ctx context.Context, userID string) (*User, error) {
	user := &User{}
	err := m.mongoDB.FindOne(ctx, db.DBNAME, collectionName, bson.D{{"userID", userID}}, user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("mongoDB.FindOne no result")
		} else {
			log.Println("mongoDB.FindOne fail")
		}
		return nil, err
	}
	return user, nil
}

func (m *Manager) UpsertBinanceKeys(ctx context.Context, userID string, apiKey string, secretKey string) error {
	err := m.mongoDB.UpdateOne(ctx, db.DBNAME, collectionName, bson.D{{"userID", userID}},
		bson.D{{"binanceApiKey", apiKey}, {"binanceSecretKey", secretKey}})
	if err != nil {
		log.Println("mongoDB.UpdateOne fail")
		return err
	}
	return nil
}

func (m *Manager) GetAdminUserID() string {
	return "ADMIN_001"
}
