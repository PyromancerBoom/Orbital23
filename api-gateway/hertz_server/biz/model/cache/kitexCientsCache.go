package cache

import (
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/client"
	genericClient "github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	consul "github.com/kitex-contrib/registry-consul"
	"go.uber.org/zap"
)

// Hashmap of service name : client, small case cuz accees is only given through getter functions.
var kitexGenericClientsMap map[string]genericClient.Client

// Dependant on the idl mappings cache
func UpdateGenericClientsCache() error {

	//we make a new hashmap everytime we update because client may delete an idl file but it will remain in our hashmap if we dont
	//just keep updating the hashmap with keys.
	kitexGenericClientsMap = make(map[string]genericClient.Client)

	//Makes a registry object to be used for all the generic clients
	registry, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		zap.L().Error("Error while getting registry", zap.Error(err))
	}

	//loop for all key value pair in the idlMapping
	for serviceName, idl := range idlMappings {

		// provider initialisation
		provider, err := generic.NewThriftContentProvider(idl, nil)
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
		genClient, err := genericClient.NewClient(serviceName, thriftGeneric, client.WithResolver(registry))
		if err != nil {
			zap.L().Error("Error while initializing generic client", zap.Error(err))
			return err
		}

		kitexGenericClientsMap[serviceName] = genClient

		fmt.Print(serviceName, idl)
	}

	zap.L().Debug("Successfully cached Kitex Clients to memory.")
	return nil
}

func GetGenericClient(serviceName string) (genericClient.Client, error) {

	genCli, ok := kitexGenericClientsMap[serviceName]

	if !ok {
		return nil, errors.New("The service does not exist.")
	}

	return genCli, nil
}
