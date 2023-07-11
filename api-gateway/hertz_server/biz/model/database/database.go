package database

// Package for managing database operations for easier management of code

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// contains the configuration for connecting to MongoDB
type MongoDBConfig struct {
	URI      string // URI is the connection string for MongoDB
	Database string // Database is the name of the MongoDB database
}

// holds the MongoDB client and database connection
type MongoDBClient struct {
	Client     *mongo.Client // Mongo client
	Database   *mongo.Database
	Collection string
}

// establishes a connection to MongoDB using the provided configuration.
func ConnectToDB(config MongoDBConfig) (*MongoDBClient, error) {
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

// closes the mongodb client connection
func (m *MongoDBClient) CloseClientConn() {
	err := m.Client.Disconnect(context.Background())
	if err != nil {
		log.Println(err)
	}
}

// stores the provided data in the specified collection.
func (m *MongoDBClient) StoreData(collectionName string, data interface{}) error {
	collection := m.Database.Collection(collectionName)

	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// retrieves all documents from the specified collection.
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
