package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type Service struct {
	ServiceOwner      string   `json:"serviceOwner"`
	APIKey            string   `json:"apiKey"`
	RegisteredServers []Server `json:"registeredServers"`
}

type Server struct {
	ServiceName         string     `json:"serviceName"`
	ServiceDescription  string     `json:"serviceDescription"`
	ServerAddress       string     `json:"serverAddress"`
	LastPingedAt        string     `json:"lastPingedAt"`
	ServiceVersion      string     `json:"serviceVersion"`
	HealthCheckEndpoint string     `json:"healthCheckEndpoint"`
	Endpoints           []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	IDL    string `json:"idl"`
}

var servicesMap map[string]Service

func Register(ctx context.Context, c *app.RequestContext) {
	var req Service
	reqBody, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, "Request body is missing")
		return
	}
	buf := bytes.NewBuffer(reqBody)

	// Decode the JSON request
	err = json.NewDecoder(buf).Decode(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Check if the ownerId is already registered
	if isAlreadyRegistered(req.ServiceOwner) {
		c.String(consts.StatusBadRequest, "Already registered")
		return
	}

	apiKey := uuid.New().String()
	req.APIKey = apiKey

	// Store the service information
	if servicesMap == nil {
		servicesMap = make(map[string]Service)
	}

	servicesMap[apiKey] = req

	res := make(map[string]string)
	res["apiKey"] = apiKey
	res["Message"] = "Registered successfully. You're good to \"GO\" :D"

	c.JSON(consts.StatusOK, res)
}

func DisplayAll(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, servicesMap)
}

func isAlreadyRegistered(ownerId string) bool {
	for _, service := range servicesMap {
		if service.ServiceOwner == ownerId {
			return true
		}
	}
	return false
}
