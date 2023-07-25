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
	"errors"
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

const (
	MASTER_API_KEY         = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	ttl                    = 10 * time.Second // Declare unhealthy after
	ttd                    = 6 * ttl          // Remove from registry afer
	consulAddr             = "127.0.0.1:8500"
	healthCheckServiceName = "RegistryProxy"
	connectMethodName      = "connectServer"
	healthcheckMethodName  = "healthCheckServer"
)

var (
	registryResolver discovery.Resolver
	err              error
	consulClient     *consul.Client
)

func init() {
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

func authoriseConnect(apiKey string, serviceName string) bool {
	return (apiKey == MASTER_API_KEY)
}

func validateAddress(address string, port string) error {
	_, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	fmt.Println(address)

	ip := net.ParseIP(address)
	if ip == nil {
		return errors.New("Invalid Address")
	}

	return nil
}

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

	err2 := consulClient.Agent().ServiceRegister(req)
	if err2 != nil {
		return err
	}

	//performs a health check [no need for error checks as this code cannot reach unless auth is valid and registry is online.]
	go updateAsHealthy(serverId)

	return nil
}

func updateAsHealthy(checkID string) error {
	err := consulClient.Agent().UpdateTTL(checkID, "online", consul.HealthPassing)
	if err != nil {
		return err // Health update successful
	}

	return nil
}

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

func authoriseHealthCheckRight(apiKey string, serverID string) bool {
	//the fact that adin has serverID is enough to know he is authenticated. If the ID is invalid, an error will be thrown.
	return true
}
