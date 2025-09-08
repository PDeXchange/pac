package helpers

import (
	"context"

	"github.com/PDeXchange/pac/test/e2e/tests/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	db = "pac"
)

type MongoClient struct {
	client *mongo.Client
}

func GetMongoClient() (*MongoClient, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(config.Current.MongoDBURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoClient{client: client}, nil
}

func (m *MongoClient) Disconnect() error {
	return m.client.Disconnect(context.Background())
}

func (m *MongoClient) DropCollection(collectionName string) error {
	collection := m.client.Database(db).Collection(collectionName)
	return collection.Drop(context.Background())
}

func (m *MongoClient) InsertOne(collectionName string, body map[string]any) error {
	collection := m.client.Database(db).Collection(collectionName)

	var bsonDoc bson.D
	for key, value := range body {
		bsonDoc = append(bsonDoc, bson.E{Key: key, Value: value})
	}

	_, err := collection.InsertOne(context.Background(), bsonDoc)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoClient) FindOne(collectionName string) (bson.M, error) {
	collection := m.client.Database(db).Collection(collectionName)

	var result bson.M
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&result)
	return result, err
}
