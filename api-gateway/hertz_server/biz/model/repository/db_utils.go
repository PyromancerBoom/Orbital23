package repository

// Contains all utility functions for operations specific only to MongoDB
// Includes methods for storing, updating, deleting, etc

/*
This package contains the following utility method:

-- func StoreAdminInfo(adminConfig AdminConfig) error
Stores the admin data in the database by adding a document to the existing collection.

-- func DeleteAdminInfo(ownerID string, ownerName string) error
Deletes the admin data document from the database based on the owner ID and owner name.
UpdateAdminInfo

-- func UpdateAdminInfo(ownerID string, adminConfig AdminConfig) error
Updates the admin data in the database based on the owner ID by overwriting the existing data with the provided adminConfig data.
GetAdminInfoByID

-- func GetAdminInfoByID(ownerID string) (AdminConfig, error)
Fetches the admin data from the database based on the owner ID and returns the corresponding AdminConfig data.
GetApiKey

-- func GetApiKey(ownerID string) (string, error)
Fetches the API key from the database based on the owner ID and returns the API key as a string.

*/

import (
	"context"
	"fmt" // Using fmt for error printing. Change to hertz error code later

	"github.com/hertz-contrib/logger/zap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
		zap.L().Errorf("failed to store admin info: ", zap.Error(err))
		return fmt.Errorf("failed to store admin info: %w", err)
	}

	zap.L().Info("Stored admin info", zap.Any("adminConfig", adminConfig))
	return nil
}

// Delete the admin data document in the database
// No authentication implemented yet for this.
// @Params:
// - ownerID: string - The owner ID of the admin data to delete
// - ownerName: string - Name of the owner
// @Returns:
// - error: An error if any
func DeleteAdminInfo(ownerID string, ownerName string) error {
	collection := Client.Database(db_name).Collection(collection_name)

	filter := bson.M{"ownerid": ownerID, "ownername": ownerName}

	// delete admin info
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		zap.L().Error("failed to delete admin info", zap.Error(err))
		return fmt.Errorf("failed to delete admin info: %w", err)
	}

	zap.L().Info("Deleted admin info for", zap.String("ownerID", ownerID), zap.String("ownerName", ownerName))
	return nil
}

// Updates admin data based on the owner ID
// @Params:
// - ownerID: string - The owner ID of the admin data to update
// - adminConfig: AdminConfig - Data to be updated and overwritten
// @Returns:
// - error: An error if any
func UpdateAdminInfo(ownerID string, adminConfig AdminConfig) error {
	collection := Client.Database(db_name).Collection(collection_name)

	filter := bson.M{"ownerid": ownerID}

	// update admin info
	_, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": adminConfig})
	if err != nil {
		zap.L().Error("failed to update admin info", zap.Error(err))
		return fmt.Errorf("failed to update admin info: %w", err)
	}

	zap.L().Info("Updated admin info for", zap.String("ownerID", ownerID), zap.String("admin config:", adminConfig))
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
	filter := bson.M{"ownerid": ownerID}
	err := collection.FindOne(context.Background(), filter).Decode(&adminConfig)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Info("Admin data not found for ownerID", zap.String("ownerID", ownerID))
			return AdminConfig{}, fmt.Errorf("admin data not found")
		}
		zap.L().Error("Failed to get admin info", zap.String("ownerID", ownerID), zap.Error(err))
		return AdminConfig{}, fmt.Errorf("failed to get admin info: %w", err)
	}

	zap.L().Info("Fetched admin info for ownerID", zap.String("ownerID", ownerID), zap.Any("adminConfig", adminConfig))
	return adminConfig, nil
}

// Fetch API Key from the database
// @Params:
// - ownerID: string - The owner ID of the admin data to fetch
// @Returns:
// - string: The API key
// - error: An error if any
func GetApiKey(ownerID string) (string, error) {
	collection := Client.Database(db_name).Collection(collection_name)

	var adminConfig AdminConfig
	filter := bson.M{"ownerid": ownerID}
	err := collection.FindOne(context.Background(), filter).Decode(&adminConfig)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Info("Admin data not found for ownerID", zap.String("ownerID", ownerID))
			return "", fmt.Errorf("admin data not found")
		}
		zap.L().Error("Failed to get admin info", zap.String("ownerID", ownerID), zap.Error(err))
		return "", fmt.Errorf("failed to get admin info: %w", err)
	}

	zap.L().Info("Fetch operation completed")
	return adminConfig.ApiKey, nil
}
