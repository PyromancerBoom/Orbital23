package repository

import (
	"errors"

	"go.uber.org/zap"
)

// HashMap of ServiceName : IDL
var IDLMappings map[string]string

func UpdateIDLcache() error {
	//update the admins cache
	err := UpdateAdminCache()
	if err != nil {
		zap.L().Info("Error updating Admins cache.")
		return err
	}

	IDLMappings = make(map[string]string)
	for _, admin := range AdminsCache {
		for _, service := range admin.Services {
			IDLMappings[service.ServiceName] = service.IdlContent
		}
	}

	zap.L().Debug("Successfully cached IDL files to memory.")
	return nil
}

func GetServiceIDL(serviceName string) (string, error) {
	idlstring, ok := IDLMappings[serviceName]

	//throw an error if the idl file is not there.
	if !ok {
		zap.L().Info("The service does not exist.")
		return "", errors.New("The service does not exist.")
	}

	return idlstring, nil
}
