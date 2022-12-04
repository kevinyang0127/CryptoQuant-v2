package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBNAME = "cryptoQuantV2"
)

type DBURI string

const (
	URI      DBURI = "mongodb://kevin:123@mongodb:27017/?connect=direct"
	LocalURI DBURI = "mongodb://kevin:123@127.0.0.1:27017/?connect=direct"
)

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDB(uri DBURI) (db *MongoDB, disconnect func(), err error) {
	ctx := context.Background()
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	log.Println("uri: ", uri)
	clientOptions := options.Client().ApplyURI(string(uri)).SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("mongo.Connect fail")
		return nil, nil, err
	}
	disconnect = func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("client.Ping fail")
		return nil, nil, err
	}
	log.Println("mongoDB connect success!!")

	return &MongoDB{
		client: client,
	}, disconnect, nil
}

func (mgo *MongoDB) InsertOne(ctx context.Context, databaseName string, collectionName string, document interface{}) (objectID interface{}, err error) {
	data, err := bson.Marshal(document)
	if err != nil {
		log.Println("bson.Marshal fail")
		return nil, err
	}
	doc := bson.D{}
	err = bson.Unmarshal(data, &doc)
	if err != nil {
		log.Println("bson.Unmarshal fail")
		return nil, err
	}

	collection := mgo.client.Database(databaseName).Collection(collectionName)
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Println("collection.InsertOne fail")
		return nil, err
	}
	objectID = res.InsertedID
	return objectID, nil
}

// 找不到的話 err == mongo.ErrNoDocuments
func (mgo *MongoDB) FindOne(ctx context.Context, databaseName string, collectionName string, filter bson.D, result interface{}) error {
	collection := mgo.client.Database(databaseName).Collection(collectionName)

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			log.Println("FindOne no result")
		} else {
			log.Println("collection.FindOne fail")
		}
		return err
	}

	return nil
}

/*
results must be pointer to slice,
if there is no result than results will be empty with no error
*/
func (mgo *MongoDB) FindAll(ctx context.Context, databaseName string, collectionName string, filter bson.D, results interface{}) error {
	collection := mgo.client.Database(databaseName).Collection(collectionName)

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("collection.Find fail")
		return err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, results)
	if err != nil {
		log.Println("cur.All fail")
		return err
	}

	return nil
}

// no matched will return error
func (mgo *MongoDB) UpdateOne(ctx context.Context, databaseName string, collectionName string, filter bson.D, update bson.D) error {
	collection := mgo.client.Database(databaseName).Collection(collectionName)

	result, err := collection.UpdateOne(ctx, filter, bson.D{{"$set", update}})
	if err != nil {
		log.Println("collection.UpdateOne fail")
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("collection.UpdateOne no matched")
	}

	return nil
}

// no matched will return error
func (mgo *MongoDB) DeleteOne(ctx context.Context, databaseName string, collectionName string, filter bson.D) error {
	collection := mgo.client.Database(databaseName).Collection(collectionName)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("collection.DeleteOne fail")
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("collection.DeleteOne no matched")
	}

	return nil
}
