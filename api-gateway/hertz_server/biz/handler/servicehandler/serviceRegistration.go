package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type RegisteredServer struct {
	ServerUrl string `json:"ServerUrl"`
	Port      int    `json:"Port"`
}

type Service struct {
	ServiceId          string             `json:"ServiceId"`
	ServiceName        string             `json:"ServiceName"`
	ExposedMethod      string             `json:"ExposedMethod"`
	Path               string             `json:"Path"`
	IdlContent         string             `json:"IdlContent"`
	Version            string             `json:"version"`
	ServiceDescription string             `json:"ServiceDescription"`
	ServerCount        int                `json:"ServerCount"`
	RegisteredServers  []RegisteredServer `json:"RegisteredServers"`
}

type Config struct {
	ApiKey    string    `json:"ApiKey"`
	OwnerName string    `json:"OwnerName"`
	OwnerId   string    `json:"OwnerId"`
	Services  []Service `json:"Services"`
}

// To be stored in DB later and cached in the gateway
var registrationData map[string]Config

// Make a Set of OwnerIds
var ownerIds map[string]bool

func Register(ctx context.Context, c *app.RequestContext) {
	var req []struct {
		Service Service `json:"Service"`
	}

	// Getting the request Body
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

	// Add logic for registration here

	// Check if owner ID already exists
	if isAlreadyRegistered(reqBody.OwnerId) {
		c.String(consts.StatusBadRequest, "Owner ID already exists")
		return
	}

	// Generate a UUID Api key for that owner
	apiKey := uuid.New().String()
	registrationData[reqBody.OwnerId].ApiKey = apiKey

	// If not, add the owner Id to set and Hashmap
	ownerIds[reqBody.OwnerId] = true

	// Add the owner ID as the key and the entire request body as the value
	registrationData[reqBody.OwnerId] = reqBody

	res := make(map[string]string)
	res["Message"] = "Registered successfully. You're good to GO :D"
	res["Api Key"] = apiKey

	c.JSON(consts.StatusOK, res)
}

// Function to check if owner ID already exists using the ownderId set
func isAlreadyRegistered(ownerId string) bool {
	_, ok := ownerIds[ownerId]
	return ok
}

func DisplayAll(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, registrationData)
}
