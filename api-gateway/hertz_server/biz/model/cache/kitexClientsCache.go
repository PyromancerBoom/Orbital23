package cache

import (
	"errors"

	"github.com/cloudwego/kitex/client"
	genericClient "github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	consul "github.com/kitex-contrib/registry-consul"
	"go.uber.org/zap"
)

// Hashmap of ServiceName+Path : Kitex Client
var kitexGenericClientsMap map[string]genericClient.Client

// updateKitexClientsCache updates the kitex clients cache.
// It fetches data from the idlmapping cache.
// @Returns:
// - error: An error, if any. Else nil.
func updateKitexClientsCache() error {

	// Reset hashmap before storing new values.
	kitexGenericClientsMap = make(map[string]genericClient.Client)

	// Makes a registry object to be used for all the generic clients
	registry, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		zap.L().Error("Error while getting registry", zap.Error(err))
	}

	// Loop through all key value pair in the idlMapping
	for keyName, serviceDetails := range idlMappings {

		// provider initialisation
		provider, err := generic.NewThriftContentProvider(serviceDetails.IDL, nil)
		if err != nil {
			zap.L().Error("Error while initializing thrict content provider", zap.Error(err))
			return err
		}

		thriftGeneric, err := generic.JSONThriftGeneric(provider)
		if err != nil {
			zap.L().Error("Error while creating JSONThriftGeneric", zap.Error(err))
			return err
		}

		// Fetch hostport from registry later
		kitexGenClient, err := genericClient.NewClient(keyName, thriftGeneric, client.WithResolver(registry))
		if err != nil {
			zap.L().Error("Error while initializing generic client", zap.Error(err))
			return err
		}

		kitexGenericClientsMap[keyName] = kitexGenClient
	}

	zap.L().Debug("Cached Kitex Clients.")
	return nil
}

// GetGenericClient retrieves the kitex client for a given Service and Method.
// If the Service and Method combination does not exist, it returns an error.
// @Params:
// - serviceName: string - The name of the service.
// - methodPath: string - The method path masking the method name.
// @Returns:
// - genericClient.Client: Kitex client for the given Service.
// - error: An error if the given Service and Method combination does not exist. Else, nil.
func GetGenericClient(serviceName string, methodPath string) (genericClient.Client, error) {

	genCli, ok := kitexGenericClientsMap[serviceName+"+"+methodPath]

	if !ok {
		return nil, errors.New("The service/method does not exist.")
	}

	return genCli, nil
}
