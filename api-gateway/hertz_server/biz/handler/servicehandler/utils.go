package servicehandler

/*
This package contains utility methods for the servicehandler package.

- func ownerIdExists(ownerId string) bool
  Checks if the ownerID is already registered in the database.


- func apiKeyValid(apiKey string, ownerId string) bool
  Checks if the provided API key is valid for the given ownerID.
*/

import (
	repository "api-gateway/hertz_server/biz/model/repository"

	"go.uber.org/zap"
)

// Define master api key for testing purposees only
const MASTERAPIKEY = "masterapikey"

// Utility Function to check if ownerID is already registered in the database
// @Params:
// - ownerId: string - The owner ID to check
// @Returns:
// - bool: true if already registered, false otherwise
func ownerIdExists(ownerId string) bool {
	_, err := repository.GetAdminInfoByID(ownerId)
	if err != nil {
		zap.L().Error("Owner ID does not exist: ", zap.Error(err))
		return false
	}

	return true
}

// Method to check if api key is valid for some owner id
// @Params:
// - apiKey: string - The api key to check
// - ownerId: string - The owner id to check
// @Returns:
// - bool: true if valid, false otherwise
func apiKeyValid(apiKey string, ownerId string) bool {

	// Temporary implementation to accept master key for easy testing
	// TODO: Remove on final code
	if apiKey == MASTERAPIKEY {
		return true
	}
	// ----------
	adminConfig, err := repository.GetAdminInfoByID(ownerId)
	if err != nil {
		zap.L().Error("Admin config not found: ", zap.Error(err))
		return false
	}

	if adminConfig.ApiKey == apiKey {
		return true
	}

	return false
}
