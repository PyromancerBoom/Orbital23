package database

// Package for DB Management

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig contains the configuration for connecting to MongoDB
type MongoDBConfig struct {
	URI      string
	Database string
}

// MongoDBClient holds the MongoDB client and database connection
type MongoDBClient struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection string
}

// ConnectMongoDB establishes a connection to MongoDB using the provided configuration
func ConnectMongoDB(config MongoDBConfig) (*MongoDBClient, error) {
	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db := client.Database(config.Database)
	return &MongoDBClient{
		Client:     client,
		Database:   db,
		Collection: "services", // Replace with the actual collection name
	}, nil
}

// Close closes the MongoDB client connection
func (m *MongoDBClient) Close() {
	err := m.Client.Disconnect(context.Background())
	if err != nil {
		log.Println(err)
	}
}

// StoreData stores the provided data in the specified collection
func (m *MongoDBClient) StoreData(collectionName string, data interface{}) error {
	collection := m.Database.Collection(collectionName)

	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetAllData retrieves all documents from the specified collection
func (m *MongoDBClient) GetAllData(collectionName string, result interface{}) error {
	collection := m.Database.Collection(collectionName)

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		return err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), result); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
