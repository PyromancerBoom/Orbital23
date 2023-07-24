package repository

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ServiceDetails represents the details of the service, including IDL.
type ServiceDetails struct {
	Service string
	Path    string
	Method  string // Added the ExposedMethod field
	IDL     string
}

// HashMap of ServiceName_Path: IDL
var IDLMappings map[string]ServiceDetails

const updateInterval = time.Second * 5

// UpdateIDLcache updates the IDL cache by populating the IDLMappings map.
// It fetches data from the AdminsCache and stores service names and their corresponding IDL content.
// @Returns:
// - error: An error if any
func updateIDLcache() error {
	// Update the admins cache
	err := UpdateAdminCache()
	if err != nil {
		zap.L().Error("Error updating Admins cache.", zap.Error(err))
		return err
	}

	// Clear the IDLMappings map before populating it again
	IDLMappings = make(map[string]ServiceDetails)

	for _, admin := range AdminsCache {
		for _, service := range admin.Services {
			IDLMappings[service.ServiceName+"_"+service.Path] = ServiceDetails{
				Service: service.ServiceName,
				Path:    service.Path,
				Method:  service.ExposedMethod,
				IDL:     service.IdlContent,
			}
		}
	}
	return nil
}

// UpdateIDLcacheLoop calls the updateIDLcache in a loop concurrently to keep updating
// the IDL Mappings in an interval (updateInterval)
// It fetches data from the AdminsCache and stores service names and their corresponding IDL content.
// @Returns:
// - error: An error if any
func UpdateIDLcacheLoop() error {
	// Run the function immediately to update the cache initially
	if err := updateIDLcache(); err != nil {
		zap.L().Error("Error updating IDL cache.", zap.Error(err))
	}

	// Run the function in a loop at the specified interval
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// zap.L().Info("IDL Mappings", zap.Any("IDL Mappings", IDLMappings))
			if err := updateIDLcache(); err != nil {
				zap.L().Error("Error updating IDL cache.", zap.Error(err))
			}
			// else {
			// 	zap.L().Debug("Cached")
			// }
		}
	}
}

// GetServiceDetails retrieves the complete ServiceDetails struct for a given service name and path from the IDLMappings cache.
// If the service name and path combination is not found in the cache, it returns an error.
// @Params:
// - serviceName: string - The name of the service for which to retrieve the details.
// - path: string - The path of the service for which to retrieve the details.
// @Returns:
// - ServiceDetails: The complete ServiceDetails struct for the specified service name and path.
// - error: An error if the service name and path combination does not exist in the cache.
func GetServiceDetails(serviceName string, path string) (ServiceDetails, error) {
	idlDetails, ok := IDLMappings[serviceName+"_"+path]

	// Throw an error if the IDL file is not found.
	if !ok {
		err := fmt.Errorf("The service does not exist.")
		zap.L().Error(err.Error(), zap.String("serviceName", serviceName), zap.String("servicePath", path))
		return ServiceDetails{}, err
	}

	return idlDetails, nil
}

// GetServiceIDL retrieves the IDL content for a given service name from the IDLMappings cache.
// If the service name is not found in the cache, it returns an error.
// @Params:
// - serviceName: string - The name of the service for which to retrieve the IDL content.
// @Returns:
// - string: The IDL content for the specified service name.
// - error: An error if the service name does not exist in the cache.
// func GetServiceIDL(serviceName string) (string, error) {
// 	idlstring, ok := IDLMappings[serviceName]

// 	// Throw an error if the IDL file is not found.
// 	if !ok {
// 		zap.L().Error("The service does not exist.", zap.String("serviceName", serviceName))
// 		return "", errors.New("The service does not exist.")
// 	}

// 	return idlstring, nil
// }
