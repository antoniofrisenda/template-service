package mongo

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient interface {
	GetDB() *mongo.Database
}

var (
	instance  MongoClient
	singleton sync.Once
)

type mongoClient struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoClient(ctx context.Context, uri, dbName string) (MongoClient, error) {
	singleton.Do(func() {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}

		if err := client.Ping(ctx, nil); err != nil {
			panic(err)
		}

		instance = &mongoClient{
			client:   client,
			database: client.Database(dbName),
		}
	})

	if instance == nil {
		return nil, fmt.Errorf("mongo client not Init")
	}

	return instance, nil
}

func (m *mongoClient) GetDB() *mongo.Database {
	return m.database
}
