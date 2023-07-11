package database

// Contains all utility functions for operations with MongoDB
// Works with storageStruct.go

import (
	"context"
	"encoding/json"
	"fmt" // Using fmt for error printing. Change to hertz error code later

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// stores the client data in the database
// Returns nothing
func (db *Database) StoreClientData(clientData *ClientData) error {
	// Marshal the client data into JSON bytes
	clientDataBytes, err := json.Marshal(clientData)
	if err != nil {
		return fmt.Errorf("failed to marshal client data: %v", err)
	}

	// Unmarshal the client data JSON bytes into a BSON map
	var clientDataMap bson.M
	err = bson.UnmarshalExtJSON(clientDataBytes, true, &clientDataMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal client data: %v", err)
	}

	// Insert the client data into the database collection
	_, err = db.collection.InsertOne(context.Background(), clientDataMap)
	if err != nil {
		return fmt.Errorf("failed to store client data: %v", err)
	}

	return nil
}

// GetClientDataByOwnerID fetches client data based on the owner ID
// Returns ClientData
func (db *Database) GetClientDataByOwnerID(ownerID string) (*ClientData, error) {
	// Create a filter based on the owner ID
	filter := bson.M{"OwnerId": ownerID}

	// Find the client data in the database collection
	result := db.collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("client data not found")
		}
		return nil, fmt.Errorf("failed to fetch client data: %v", result.Err())
	}

	// Decode the client data BSON map into a clientDataMap variable
	var clientDataMap bson.M
	err := result.Decode(&clientDataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode client data: %v", err)
	}

	// Marshal the clientDataMap into JSON bytes
	clientDataBytes, err := bson.MarshalExtJSON(clientDataMap, true, true)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal client data: %v", err)
	}

	// Unmarshal the client data JSON bytes into a ClientData struct
	var clientData ClientData
	err = json.Unmarshal(clientDataBytes, &clientData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal client data: %v", err)
	}

	return &clientData, nil
}

// fetches all the registered client data
// TODO : Make private later
// This function only to be used for server administration purposes
func (db *Database) GetAllClientData() ([]*ClientData, error) {
	// Find all client data in the database collection
	cursor, err := db.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch client data: %v", err)
	}
	defer cursor.Close(context.Background())

	var clientDataList []*ClientData
	// Iterate over the cursor and decode each client data into a ClientData struct
	for cursor.Next(context.Background()) {
		var clientDataMap bson.M
		err := cursor.Decode(&clientDataMap)
		if err != nil {
			return nil, fmt.Errorf("failed to decode client data: %v", err)
		}

		// Marshal the clientDataMap into JSON bytes
		clientDataBytes, err := bson.MarshalExtJSON(clientDataMap, true, true)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal client data: %v", err)
		}

		// Unmarshal the client data JSON bytes into a ClientData struct
		var clientData ClientData
		err = json.Unmarshal(clientDataBytes, &clientData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal client data: %v", err)
		}

		clientDataList = append(clientDataList, &clientData)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return clientDataList, nil
}
