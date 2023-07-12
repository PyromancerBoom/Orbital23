package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertTestData() error {
	// Set up the MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Access the "testDB" database and "testCollection" collection
	database := client.Database("testDB")
	collection := database.Collection("testCollection")

	// Create a document to be inserted
	document := bson.D{
		{"name", "John Doe"},
		{"age", 30},
		{"email", "john.doe@example.com"},
	}

	// Insert the document into the collection
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	// Disconnect from the MongoDB instance
	err = client.Disconnect(context.Background())
	if err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	return nil
}
