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

// Handler for "/register"
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

	// Check if ownerId exists
	if ownerIdExists(req[0]["OwnerId"].(string)) {
		c.String(consts.StatusBadRequest, "Owner ID already exists")
		return
	}

	apiKey := uuid.New().String()

	// Store the admin info in the database
	adminConfig := repository.AdminConfig{
		ApiKey:    apiKey,
		OwnerName: req[0]["OwnerName"].(string),
		OwnerId:   req[0]["OwnerId"].(string),
	}

	// Unmarshal the Services field from the request into the Services field of adminConfig
	servicesJSON, err := json.Marshal(req[0]["Services"])
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse Services field")
		return
	}
	err = json.Unmarshal(servicesJSON, &adminConfig.Services)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse Services field")
		return
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
