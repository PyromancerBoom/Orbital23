package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type ClientData struct {
	ApiKey    string    `json:"ApiKey"`
	OwnerName string    `json:"OwnerName"`
	OwnerId   string    `json:"OwnerId"`
	Services  []Service `json:"Services"`
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

type RegisteredServer struct {
	ServerUrl string `json:"ServerUrl"`
	Port      int    `json:"Port"`
}

// To be stored in DB later and cached in the gateway
var servicesMap map[string]ClientData

// Make a Set of OwnerIds
var ownerIdSet map[string]bool

func Register(ctx context.Context, c *app.RequestContext) {
	var req []map[string]interface{}

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

	apiKey := uuid.New().String()

	// Iterate over the request items
	for _, item := range req {
		// Check if "ClientData" key exists in the item
		if clientDataRaw, ok := item["ClientData"]; ok {
			// Convert the value to bytes for decoding
			clientDataBytes, err := json.Marshal(clientDataRaw)
			if err != nil {
				c.String(consts.StatusBadRequest, "Failed to parse ClientData")
				return
			}

			// Decode the ClientData JSON
			var clientData ClientData
			err = json.Unmarshal(clientDataBytes, &clientData)
			if err != nil {
				c.String(consts.StatusBadRequest, "Failed to parse ClientData")
				return
			}

			// Check if owner ID already exists
			if isAlreadyRegistered(clientData.OwnerId) {
				c.String(consts.StatusBadRequest, "Owner ID already exists")
				return
			}

			// if not, then add owner ID to the set ownerIdSet
			// assignment to entry in nil map
			if ownerIdSet == nil {
				ownerIdSet = make(map[string]bool)
			}
			ownerIdSet[clientData.OwnerId] = true

			clientData.ApiKey = apiKey

			// Store the client data information
			if servicesMap == nil {
				servicesMap = make(map[string]ClientData)
			}

			servicesMap[apiKey] = clientData
		}
	}

	res := make(map[string]string)
	res["Message"] = "Registered successfully. You're good to GO :D"
	res["api-key"] = apiKey

	c.JSON(consts.StatusOK, res)
}

// Function to check if owner ID already exists using the ownerIds map
func isAlreadyRegistered(ownerId string) bool {
	_, ok := ownerIdSet[ownerId]
	return ok
}

func DisplayAll(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, servicesMap)
}
