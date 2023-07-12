package repository

// Contains all utility functions for operations with MongoDB
// Works with storageStruct.go

import (
	"context"
	"fmt" // Using fmt for error printing. Change to hertz error code later

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Stores the admin data in the database
// Adds document to existing collection
// Returns an error if any
func StoreAdminInfo(adminConfig AdminConfig) error {
	collection := Client.Database(db_name).Collection(collection_name)

	_, err := collection.InsertOne(context.Background(), adminConfig)
	if err != nil {
		return fmt.Errorf("failed to store admin info: %w", err)
	}

	return nil
}

// Updates admin data based on the owner ID
// Params:
// - ownerID: string - The owner ID of the admin data to update
// - isValid: bool - The updated value for the "isValid" field
// Returns:
// - error: An error if any occurred during the update operation
func UpdateAdminInfo(ownerID string, isValid bool) error {
	collection := Client.Database("testDB").Collection("adminCollection")

	filter := bson.M{"OwnerId": ownerID}
	update := bson.M{"$set": bson.M{"isValid": isValid}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update admin info: %w", err)
	}

	return nil
}

// Fetches admin data based on the owner ID
// Returns the admin data and an error if any
func GetAdminInfoByID(ownerID string) (AdminConfig, error) {
	collection := Client.Database(db_name).Collection(collection_name)

	var adminConfig AdminConfig
	filter := bson.M{"OwnerId": ownerID}
	err := collection.FindOne(context.Background(), filter).Decode(&adminConfig)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return AdminConfig{}, fmt.Errorf("admin data not found")
		}
		return AdminConfig{}, fmt.Errorf("failed to get admin info: %w", err)
	}

	return adminConfig, nil
}

// Fetches all the registered client data documents
// Returns a slice of RegisteredServer and an error if any
// Todo : Remove method later
func GetAllClientData() ([]RegisteredServer, error) {
	collection := Client.Database(db_name).Collection(collection_name)

	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch client data: %w", err)
	}
	defer cur.Close(context.Background())

	var registeredServers []RegisteredServer
	for cur.Next(context.Background()) {
		var adminConfig AdminConfig
		if err := cur.Decode(&adminConfig); err != nil {
			return nil, fmt.Errorf("failed to decode client data: %w", err)
		}

		registeredServers = append(registeredServers, adminConfig.Services[0].RegisteredServers...)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while fetching client data: %w", err)
	}

	return registeredServers, nil
}
