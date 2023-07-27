// Package used for caching Admin details, IDL files, Kitex Clients for improved performance.
package cache

import (
	"time"

	"go.uber.org/zap"
)

// UpdateCache updates the Admins Cache, IDL Cache and Kitex Generic Clients cache.
// @Returns:
// - error: An error if any
func UpdateCache() error {

	err := updateAdminCache()
	if err != nil {
		return err
	}

	err = updateIDLcache()
	if err != nil {
		return err
	}

	err = updateKitexClientsCache()
	if err != nil {
		return err
	}

	err = updatePathMaskCache()
	if err != nil {
		return err
	}

	return nil
}

// UpdateCacheLoop calls the UpdateCache in a loop to keep updating the stored caches.
// This method is kept for future use, if need be.
// @Returns:
// - error: An error if any
func UpdateCacheLoop(updateInterval time.Duration) error {
	// Run the function immediately to update the cache initially
	if err := UpdateCache(); err != nil {
		zap.L().Error("Error updating IDL cache.", zap.Error(err))
	}

	// Run the function in a loop at the specified interval
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// zap.L().Info("IDL Mappings", zap.Any("IDL Mappings", IDLMappings))
			if err := UpdateCache(); err != nil {
				zap.L().Error("Error updating IDL cache.", zap.Error(err))
			}
		}
	}
}
