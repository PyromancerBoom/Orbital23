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
	return err == nil
}
