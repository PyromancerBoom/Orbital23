package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertTestData() error {
	// Access the "testDB" database and "testCollection" collection
	database := Client.Database("testDB")
	collection := database.Collection("testCollection")

	// Create a document to be inserted
	document := bson.D{
		{"name", "John Doe"},
		{"age", 30},
		{"email", "john.doe@example.com"},
	}

	// Insert the document into the collection
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}
