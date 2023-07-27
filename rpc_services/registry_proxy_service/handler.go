package main

import (
	"context"
	"log"
	registry_proxy_service "registry_proxy/kitex_gen/registry_proxy_service"

	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
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

		zap.L().Info("Unauthorized connection for server using " + req.ApiKey)
		return &registry_proxy_service.ConnectResponse{
			Status:   "failed",
			Message:  "Server connection unauthorized.",
			ServerID: "",
		}, nil
	}

	err = validateAddress(req.ServerAddress, req.ServerPort)
	if err != nil {

		zap.L().Info("Invalid server address provided for server using " + req.ApiKey)
		return &registry_proxy_service.ConnectResponse{
			Status:   "failed",
			Message:  "Server address is invalid.",
			ServerID: "",
		}, nil
	}

	serverId := uuid.New().String()

	err = registerServer(req.ServerAddress, req.ServerPort, serverId, req.ServiceName, req.ApiKey)
	if err != nil {

		zap.L().Error("Error registering server with key "+req.ApiKey, zap.Error(err))
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

	if !authoriseHealthCheckRight(req.ApiKey, req.ServerID) {
		zap.L().Info("Unauthorized health update request for server id " + req.ServerID)
		return &registry_proxy_service.HealthResponse{
			Status:  "failed",
			Message: "Unauthorized to update health.",
		}, nil
	}

	//The fact that user has the serverID means that he is authorised to perform the health check.
	err = updateAsHealthy(req.ServerID)
	if err != nil {
		//Note, error will also be thrown if the server ID is invalid.
		zap.L().Error("Error occured while updating health", zap.Error(err))
		//throw the error
		return nil, err
	}

	return &registry_proxy_service.HealthResponse{
		Status:  "ok",
		Message: "Successfully updated server health",
	}, nil
}
