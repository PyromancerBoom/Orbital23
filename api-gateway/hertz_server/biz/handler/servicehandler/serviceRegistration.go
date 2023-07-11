package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"

	db_mongo "api-gateway/hertz_server/biz/model/db_mongo"
)

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
			var clientData db_mongo.ClientData
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

			clientData.ApiKey = apiKey

			// Store the client data information in MongoDB
			db_mongo.StoreClientData(&clientData)
		}
	}

	res := make(map[string]string)
	res["Message"] = "Registered successfully. You're good to GO :D"
	res["api-key"] = apiKey

	c.JSON(consts.StatusOK, res)
}

// Function to check if owner ID already exists using the ownerIds map
func isAlreadyRegistered(ownerId string) bool {
	clientData, err := db_mongo.GetClientDataByOwnerID(ownerId)
	if err != nil {
		return false
	}
	return clientData != nil
}

func DisplayAll(ctx context.Context, c *app.RequestContext) {
	clientDataList, err := db_mongo.GetAllClientData()
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to fetch client data")
		return
	}

	c.JSON(consts.StatusOK, clientDataList)
}
