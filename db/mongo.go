package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDB() (db *MongoDB, disconnect func(), err error) {
	ctx := context.Background()
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	uri := "mongodb://kevin:123@mongodb:27017/?connect=direct"
	log.Println("uri: ", uri)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
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

// 找不到的話 result == nil
func (mgo *MongoDB) FindOne(ctx context.Context, databaseName string, collectionName string, filter bson.D) (result bson.M, err error) {
	collection := mgo.client.Database(databaseName).Collection(collectionName)

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		log.Println("collection.FindOne fail")
		return nil, err
	}

	return result, nil
}

// 找不到的話 results == nil
func (mgo *MongoDB) Find(ctx context.Context, databaseName string, collectionName string, filter bson.D) ([]bson.D, error) {
	collection := mgo.client.Database(databaseName).Collection(collectionName)

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("collection.Find fail")
		return nil, err
	}
	defer cur.Close(ctx)

	results := []bson.D{}
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Println("cur.Decode fail")
			return nil, err
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Println("cur.Err fail")
		return nil, err
	}

	return results, nil
}
