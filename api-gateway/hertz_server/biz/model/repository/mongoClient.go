package repository

/*
	About this file :

- var Client *mongo.Client
  - The MongoDB client used for database operations.

- var db_name string
  - The name of the MongoDB database. Used in the utility methods

- var collection_name string
  - The name of the MongoDB collection. Used in the utility methods

- func ConnectToMongoDB() error
  - Establishes a connection to MongoDB.
  - Function called before Hertz server is initialised.
  - Returns an error if the connection fails or nil if the connection is successful.
*/

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client          *mongo.Client
	db_name         string = "testDB"
	collection_name string = "testCollection"
)

// Establishes connection to MongoDB
// @Params:
// - None
// @Returns:
// - error: An error if any
func ConnectToMongoDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	Client = c

	return nil
}
