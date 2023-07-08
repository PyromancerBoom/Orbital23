package servicehandler

//handles :/connect endpoint to register a server

//authorise if the api-key is valid

//register the server.

//request here has to have
//1: API KEY
//2: ServiceName
//3: ServerAddress

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
)

type Request struct {
	ApiKey        string `json:"api-key"`
	ServiceName   string `json:"serviceName"`
	ServerAddress string `json:"serverAddress"`
	ServerPort    string `json:"serverPort"`
}

const (
	MASTER_API_KEY = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	ttl            = 10 * time.Second //declare unhealthy after
	ttd            = 6 * ttl          //remove from registry afer
)

func Connect(ctx context.Context, c *app.RequestContext) {
	var req Request
	err := c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusExpectationFailed, "Invalid Request.")
		return
	}

	if !authoriseConnect(req.ApiKey, req.ServiceName) {
		c.String(consts.StatusUnauthorized, "Unauthorized.")
		return
	}

	if validateAddress(req.ServerAddress, req.ServerPort) != nil {
		c.String(consts.StatusBadRequest, "Address is invalid.")
		return
	}

	res := make(map[string]string)
	res["status"] = "status OK"
	res["message"] = "Server Connection Request Accepted."
	res["serverID"] = uuid.New().String()

	err2 := registerServer(req.ServerAddress, req.ServerPort, res["serverID"], req.ServiceName, req.ApiKey)
	if err2 != nil {
		c.String(consts.StatusInternalServerError, "Unable to register server.")
		return
	}

	c.JSON(consts.StatusOK, res)
}

// Add logic here
// Auhorize only if
// 1: apikey is valid
// and API key has a registered service with the provided name
func authoriseConnect(apiKey string, serviceName string) bool {
	return (apiKey == MASTER_API_KEY) && (serviceName == "UserService" || serviceName == "AssetManagement")
}

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
