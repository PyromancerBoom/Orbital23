package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	repository "api-gateway/hertz_server/biz/model/repository"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

// Handler for /register
// Registers info sent through a post request
// Generates API key for admin
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

	// Store the admin info in the database
	adminConfig := repository.AdminConfig{
		ApiKey:    apiKey,
		OwnerName: req[0]["OwnerName"].(string),
		OwnerId:   req[0]["OwnerId"].(string),
		Services:  []repository.Service{},
	}

	err = repository.StoreAdminInfo(adminConfig)
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to register admin")
		return
	}

	response := make(map[string]string)
	response["Message"] = "Registered successfully. You're good to GO :D"
	response["api-key"] = apiKey

	c.JSON(consts.StatusOK, response)
}

// Utility Function to check if ownerID is already registered in the database
// @Params:
// - ownerId: string - The owner ID to check
// @Returns:
// - bool: true if already registered, false otherwise
func isAlreadyRegistered(ownerId string) bool {
	_, err := repository.GetAdminInfoByID(ownerId)
	if err != nil {
		return false
	}
	return true
}
