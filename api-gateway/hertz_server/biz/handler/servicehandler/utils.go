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
	"api-gateway/hertz_server/biz/model/settings"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/discovery"
	consul "github.com/hashicorp/consul/api"
	consul_kitex "github.com/kitex-contrib/registry-consul"
	"go.uber.org/zap"
)

var (
	registryResolver discovery.Resolver
	err              error
	consulClient     *consul.Client
	serverSettings   settings.Settings

	ttl        time.Duration
	ttd        time.Duration
	consulAddr string

	// Define master api key for testing purposees only
	MASTERAPIKEY string
)

func init() {

	err := settings.InitialiseSettings("serverconfig.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	serverSettings = settings.GetSettings()

	ttl = time.Duration(serverSettings.TTD) * time.Second
	ttd = time.Duration(serverSettings.TTD) * time.Second
	consulAddr = serverSettings.ConsulAddress

	MASTERAPIKEY = serverSettings.MasterApiKey

	// Get registry to enable resolving serverIDs
	registryResolver, err = consul_kitex.NewConsulResolver(consulAddr)
	if err != nil {
		zap.L().Error(err.Error())
		log.Fatal("Error while getting registry")
	}

	consulClient, err = consul.NewClient(&consul.Config{})
	if err != nil {
		zap.L().Error(err.Error())
		log.Fatal("Error while fetching a conusl client")
	}
}

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

// Updates a service as Healthy
// @Params
// checkID: The ID of the service to be updated
// @Returns
// error: If any error occurs while updating the service
func updateAsHealthy(checkID string) error {
	maxRetry := 10
	for retry := 0; retry < maxRetry; retry++ {
		err := consulClient.Agent().UpdateTTL(checkID, "online", consul.HealthPassing)
		if err == nil {
			return nil // Health update successful
		}

		zap.L().Error("Failed to update health. Retrying health update in 5 seconds...", zap.Error(err))
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("failed to update health after %d retries", maxRetry)
}

// Checks if health is passing for a given service
func checkIfServiceHealthy(serviceName string) (bool, error) {
	status, _, err := consulClient.Agent().AgentHealthServiceByName("RegistryProxy")
	if err != nil {
		return false, err
	}

	if status == "passing" {
		return true, nil
	}

	return false, nil
}

// Add logic here
// Auhorize only if
// 1: apikey is valid
// and API key has a registered service with the provided name
func authoriseConnect(apiKey string, serviceName string) bool {
	return (apiKey == MASTERAPIKEY)
}

// Validates address
func validateAddress(address string, port string) error {
	_, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	if net.ParseIP(address) == nil {
		return err
	}
	return nil
}

// Registers the server
func registerServer(address string, port string, serverId string, serviceName string, apikey string) error {

	portNum, _ := strconv.Atoi(port)

	check := &consul.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttd.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        serverId,
	}

	req := &consul.AgentServiceRegistration{
		ID:      serverId,
		Name:    serviceName,
		Tags:    []string{serviceName, serverId, apikey},
		Address: address,
		Port:    portNum,
		Check:   check,
	}

	err := consulClient.Agent().ServiceRegister(req)
	if err != nil {
		return err
	}

	// Performs a health check [no need for error checks as this code cannot reach unless auth is valid and registry is online.]
	go updateAsHealthy(serverId)

	return nil
}

// Mocked.
func authoriseHealthCheckRight(apiKey string, serverID string) bool {
	return true
}
