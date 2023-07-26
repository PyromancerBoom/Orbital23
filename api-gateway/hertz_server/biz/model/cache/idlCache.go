package cache

import (
	"errors"

	"go.uber.org/zap"
)

// HashMap of ServiceName : IDL
var idlMappings map[string]string

// dependent on the admins cache
func UpdateIDLcache() error {

	//we make a new hashmap everytime we update because client may delete an idl file
	//but it will remain in our hashmap if we dont just keep updating the hashmap with keys.
	idlMappings = make(map[string]string)

	for _, admin := range adminsCache {
		for _, service := range admin.Services {
			idlMappings[service.ServiceName] = service.IdlContent
		}
	}

	zap.L().Debug("Successfully cached IDL files to memory.")
	return nil
}

func GetServiceIDL(serviceName string) (string, error) {
	idlstring, ok := idlMappings[serviceName]

	//throw an error if the idl file is not there.
	if !ok {
		zap.L().Info("The service does not exist.")
		return "", errors.New("The service does not exist.")
	}

	return idlstring, nil
}
