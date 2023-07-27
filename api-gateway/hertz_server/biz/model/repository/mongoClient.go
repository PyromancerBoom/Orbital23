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
	"api-gateway/hertz_server/biz/model/settings"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	Client          *mongo.Client
	serverSettings  settings.Settings
	db_name         string
	collection_name string
	db_url          string
)

func init() {
	err := settings.InitialiseSettings("serverconfig.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	serverSettings = settings.GetSettings()
	db_name = serverSettings.DbName
	collection_name = serverSettings.DbColletionName
	db_url = serverSettings.DbUrl
}

// Establishes connection to MongoDB
// @Params:
// - None
// @Returns:
// - error: An error if any
func ConnectToMongoDB() error {
	clientOptions := options.Client().ApplyURI(db_url)
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		zap.L().Error("Failed to connect to MongoDB", zap.Error(err))
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	Client = c

	zap.L().Debug("Connected to MongoDB")

	return nil
}

// Function to close MongoDB Client before server is shutdown
// @Params:
// - None
// @Returns:
// - error: An error if any
func CloseMongoDB() error {
	err := Client.Disconnect(context.Background())
	if err != nil {
		zap.L().Error("Failed to disconnect from MongoDB", zap.Error(err))
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	zap.L().Debug("Disconnected from MongoDB")

	return nil
}

// Perform MongoDB health check by pinging the server every 30 sec normally.
// If health check fails, ping every 5 seconds until MongoDB server comes back online.
func MongoHealthCheck() {
	normalPingInterval := time.Duration(serverSettings.DbPingInterval) * time.Second
	failedPingInterval := time.Duration(serverSettings.DbFailedPingInterval) * time.Second
	maxFailedPingDuration := time.Duration(serverSettings.DbMaxFailPingDuration) * time.Minute

	ticker := time.NewTicker(normalPingInterval)
	pingInterval := normalPingInterval
	failedPingStart := time.Time{}

	for {
		select {
		case <-ticker.C:
			err := Client.Ping(context.Background(), nil)
			if err != nil {
				zap.L().Error("Failed to ping MongoDB server", zap.Error(err))
				if pingInterval != failedPingInterval {
					zap.L().Info("Switching to fast ping interval")
					pingInterval = failedPingInterval
					ticker.Reset(pingInterval)
					failedPingStart = time.Now()
				} else if time.Since(failedPingStart) > maxFailedPingDuration {
					zap.L().Warn("MongoDB server is still offline after 3 minutes, stopping health checks")
					ticker.Stop()
					return
				}
			} else {
				if pingInterval != normalPingInterval {
					zap.L().Info("MongoDB server is back online, switching to normal ping interval")
					pingInterval = normalPingInterval
					ticker.Reset(pingInterval)
				}
			}
		}
	}
}
