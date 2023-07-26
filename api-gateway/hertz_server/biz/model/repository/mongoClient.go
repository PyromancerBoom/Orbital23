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
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	Client          *mongo.Client
	db_name         string = "api_gateway_db"
	collection_name string = "admin_services"
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
		zap.L().Error("Failed to connect to MongoDB", zap.Error(err))
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	Client = c

	zap.L().Debug("Connected to MongoDB")

	return nil
}

// Perform MongoDB health check by pinging the server every 30 sec normally.
// If health check fails, ping every 5 seconds until MongoDB server comes back online.
func MongoHealthCheck() {
	normalPingInterval := 30 * time.Second
	failedPingInterval := 5 * time.Second
	maxFailedPingDuration := 3 * time.Minute

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
