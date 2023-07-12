package repository

// Contains all utility functions for operations with MongoDB
// Works with storageStruct.go

import (
	"context"
	"fmt" // Using fmt for error printing. Change to hertz error code later

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Stores the admin data in the database, adds document to existing collection
// @Params:
// - adminConfig: AdminConfig - Data to be inserted
// @Returns
// - error: An error if any
func StoreAdminInfo(adminConfig AdminConfig) error {
	collection := Client.Database(db_name).Collection(collection_name)

	_, err := collection.InsertOne(context.Background(), adminConfig)
	if err != nil {
		return fmt.Errorf("failed to store admin info: %w", err)
	}

	return nil
}

// Updates admin data based on the owner ID
// @Params:
// - ownerID: string - The owner ID of the admin data to update
// - adminConfig: AdminConfig - Data to be updated and overwritten
// @Returns:
// - error: An error if any
func UpdateAdminInfo(ownerID string, adminConfig AdminConfig) error {
	collection := Client.Database("testDB").Collection("adminCollection")

	filter := bson.M{"OwnerId": ownerID}

	// update admin info
	_, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": adminConfig})
	if err != nil {
		return fmt.Errorf("failed to update admin info: %w", err)
	}

	return nil
}

// Fetches admin data based on the owner ID
// The validation based on API key happens outside this
// @Params:
// - ownerID: string - The owner ID of the admin data to fetch
// @Returns:
// - adminConfig: AdminConfig - The admin data
// - error: An error if any
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
