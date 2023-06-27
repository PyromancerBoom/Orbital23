package handler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

var servicesMap = make(map[string]interface{})

type registerRequest struct {
	ServiceOwner      string                       `json:"serviceOwner"`
	ApiKey            string                       `json:"apiKey"`
	RegisteredServers []RegisteredServerWithAPIKey `json:"registeredServers"`
}

type RegisteredServerWithAPIKey struct {
	RegisteredServer
	ServiceID string `json:"serviceId"`
}

type RegisteredServer struct {
	ServerAddress string               `json:"serverAddress"`
	LastPingedAt  string               `json:"lastPingedAt"`
	Service       string               `json:"service"`
	Endpoints     []RegisteredEndpoint `json:"endpoints"`
}

type RegisteredEndpoint struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Idl    string `json:"idl"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	// Parse the request payload
	var req registerRequest

	// Read the request body
	reqBody, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, "Request body is missing")
		return
	}

	buf := bytes.NewBuffer(reqBody)

	err = json.NewDecoder(buf).Decode(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Generate a unique API Key
	apiKey := uuid.New().String()

	servicesMap[req.ApiKey] = apiKey

	// Create a map to store the registered servers
	registeredServers := make(map[string]interface{})

	for _, registeredServer := range req.RegisteredServers {
		// Generate a unique Service ID using UUID
		serviceID := uuid.New().String()

		registeredServerWithAPIKey := RegisteredServerWithAPIKey{
			RegisteredServer: registeredServer.RegisteredServer,
			ServiceID:        serviceID,
		}
		registeredServers[serviceID] = registeredServerWithAPIKey
	}

	response := make(map[string]interface{})
	response["Status"] = "Registered successfully"
	response["apiKey"] = apiKey
	response["Message"] = "You're good to go!"

	c.JSON(consts.StatusOK, response)
}

// displayAll returns the hashmap with all the stored details.
func DisplayAll(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, servicesMap)
}
