package repository

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database stores the MongoDB client and database information
type Database struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// establishes a connection to the MongoDB server
func ConnectDB(uri, dbName, collectionName string) (*Database, error) {
	fmt.Println("Reached Connection")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB server: %v", err)
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &Database{
		client:     client,
		database:   db,
		collection: collection,
	}, nil
}

// clos e the database connection
func (db *Database) Close() {
	err := db.client.Disconnect(context.Background())
	if err != nil {
		log.Printf("Error disconnecting from MongoDB: %v\n", err)
	}
}
