package servicehandler

import (
	repository "api-gateway/hertz_server/biz/model/repository"
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Handler for "/update"
// Updates info sent through a post request
func Update(ctx context.Context, c *app.RequestContext) {
	var req []map[string]interface{}
	// Fetch ownerId from params
	ownerId := c.Query("ownerid")
	// Fetch api key from header
	apiKey := string(c.GetHeader("apikey"))

	// Check if owner ID is valid
	if !ownerIdExists(ownerId) {
		c.String(consts.StatusBadRequest, "Invalid owner ID")
		return
	}

	// Validate api key
	if !apiKeyValid(apiKey, ownerId) {
		c.String(consts.StatusBadRequest, "Invalid API key")
		return
	}

	// Getting the request Body
	reqBody, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, "Request body is missing"+err.Error())
		return
	}
	buf := bytes.NewBuffer(reqBody)

	// Decode the JSON request
	err = json.NewDecoder(buf).Decode(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body"+err.Error())
		return
	}

	// Prepare the updated admin info
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
		c.String(consts.StatusBadRequest, "Failed to parse Services field"+err.Error())
		return
	}

	err = repository.UpdateAdminInfo(ownerId, adminConfig)
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to update admin")
		return
	}

	// Return the reponse
	response := make(map[string]string)
	response["Message"] = "Updated successfully. You're good to GO :D"

	c.JSON(consts.StatusOK, response)
}
