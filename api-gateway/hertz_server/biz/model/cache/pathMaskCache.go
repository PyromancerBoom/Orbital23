package cache

import (
	"errors"

	"go.uber.org/zap"
)

// Mapping of method paths to exposed paths
var pathCache map[string]string

// updatePathMaskCache updates the pathCache.
// It fetches data from adminsCache for caching method paths to exposed method name
// @Returns:
// - error: An error if any
func updatePathMaskCache() error {

	// Clear the pathCache map before populating it again
	pathCache = make(map[string]string)

	// Make the IDL Mappings - "ServiceName+Path" : Method, IDL
	for _, admin := range adminsCache {
		for _, service := range admin.Services {
			for _, path := range service.Paths {
				// Create the key for IDLMappings using ServiceName and MethodPath
				key := path.MethodPath
				value := path.ExposedMethod
				pathCache[key] = value
			}
		}
	}

	zap.L().Debug("Cached IDLs.")
	return nil
}

// GetExposedMethodFromPath retrieves the exposed method name for a given path name.
// If the path does not exist, throws an error.
// @Params:
// - methodPath: string - The masked method path.
// @Returns:
// - string: Name of the exposed method.
// - error: An error if the given methodPath does not exist.
func GetExposedMethodFromPath(methodPath string) (string, error) {
	methodName, ok := pathCache[methodPath]

	if !ok {
		return "", errors.New("The path does not exist.")
	}

	return methodName, nil
}
