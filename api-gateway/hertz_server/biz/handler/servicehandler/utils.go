package servicehandler

// Utility methods for servicehandler

import repository "api-gateway/hertz_server/biz/model/repository"

// Utility Function to check if ownerID is already registered in the database
// @Params:
// - ownerId: string - The owner ID to check
// @Returns:
// - bool: true if already registered, false otherwise
func ownerIdExists(ownerId string) bool {
	_, err := repository.GetAdminInfoByID(ownerId)
	if err != nil {
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
	adminConfig, err := repository.GetAdminInfoByID(ownerId)
	if err != nil {
		return false
	}

	if adminConfig.ApiKey == apiKey {
		return true
	}

	return false
}
