package mongo

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoClient interface {
	Connect() error
	Disconnect() error
	GetConnection() *mongo.Database
}

// singleton
var (
	instance IMongoClient
	once     sync.Once
)

type MongoClient struct {
	uri      string
	dbName   string
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoClient(uri, dbName string) IMongoClient {
	once.Do(func() {
		instance = &MongoClient{
			uri:    uri,
			dbName: dbName,
		}
	})
	return instance
}

func (m *MongoClient) Connect() error {
	if m.client != nil { 
		return nil //già connesso
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	m.client = client
	m.database = client.Database(m.dbName)
	return nil
}

func (m *MongoClient) Disconnect() error {
	if m.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		panic(err)
	}

	m.client = nil
	m.database = nil
	return nil
}

func (m *MongoClient) GetConnection() *mongo.Database {
	return m.database
}
