package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	registry_proxy_service "registry_proxy/kitex_gen/registry_proxy_service"
	"strconv"
	"time"

	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
)

const (
	MASTER_API_KEY = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	ttl            = 60 * time.Second // Declare unhealthy after
	ttd            = 3 * ttl          // Remove from registry afer
	consulAddr     = "127.0.0.1:8500"
)

var consulClient *consul.Client

func init() {
	var err error
	consulClient, err = consul.NewClient(&consul.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
}

// RegistryProxyImpl implements the last service interface defined in the IDL.
type RegistryProxyImpl struct{}

// ConnectServer implements the RegistryProxyImpl interface.
func (s *RegistryProxyImpl) ConnectServer(ctx context.Context, req *registry_proxy_service.ConnectRequest) (resp *registry_proxy_service.ConnectResponse, err error) {
	// TODO: Your code here...

	if !authoriseConnect(req.ApiKey, req.ServiceName) {
		return &registry_proxy_service.ConnectResponse{
			Status:   "failed",
			Message:  "Server connection unauthorized.",
			ServerID: "",
		}, nil
	}

	err2 := validateAddress(req.ServerAddress, req.ServerPort)
	if err2 != nil {
		return &registry_proxy_service.ConnectResponse{
			Status:   "failed",
			Message:  "Server address is invalid.",
			ServerID: "",
		}, nil
	}

	serverId := uuid.New().String()

	err = registerServer(req.ServerAddress, req.ServerPort, serverId, req.ServiceName, req.ApiKey)
	if err != nil {
		return &registry_proxy_service.ConnectResponse{
			Status:   "failed",
			Message:  "Unable to register server.",
			ServerID: "",
		}, err
	}

	return &registry_proxy_service.ConnectResponse{
		Status:   "ok",
		Message:  "Successfully connected server to gateway.",
		ServerID: serverId,
	}, nil
}

// HealthCheckServer implements the RegistryProxyImpl interface.
func (s *RegistryProxyImpl) HealthCheckServer(ctx context.Context, req *registry_proxy_service.HealtRequest) (resp *registry_proxy_service.HealthResponse, err error) {
	// TODO: Your code here...

	//The fact that user has the serverID means that he is authorised to perform the health check.
	err = updateAsHealthy(req.ServerID)
	if err != nil {
		//Note, error will also be thrown if the server ID is invalid.

		//throw the error
		return nil, err
	}

	return &registry_proxy_service.HealthResponse{
		Status:  "ok",
		Message: "Successfully updated server health",
	}, nil
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

	client, err := consul.NewClient(&consul.Config{})
	if err != nil {
		return err
	}

	err2 := client.Agent().ServiceRegister(req)
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
