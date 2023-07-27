package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

const (
	ttl        = 10 * time.Second
	ttd        = 6 * ttl
	consulAddr = "127.0.0.1:8500"

	// Define master api key for testing purposees only
	MASTERAPIKEY = "master_api_key_uuid"
)

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
