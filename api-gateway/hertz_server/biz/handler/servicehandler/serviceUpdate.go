package servicehandler

import (
	repository "api-gateway/hertz_server/biz/model/repository"
	"bytes"
	"context"
	"encoding/json"

	"api-gateway/hertz_server/biz/model/cache"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"
)

// Handler for "/update"
// Updates info sent through a post request
func Update(ctx context.Context, c *app.RequestContext) {
	var req []map[string]interface{}
	// Fetch ownerId from params
	ownerId := c.Query("ownerid")
	zap.L().Info("ownerId fetched from params", zap.String("ownerId", ownerId))
	// Fetch api key from header
	apiKey := string(c.GetHeader("apikey"))
	zap.L().Info("apiKey fetched from header", zap.String("apiKey", apiKey))

	// Check if owner ID is valid
	if !ownerIdExists(ownerId) {
		c.String(consts.StatusBadRequest, "Invalid owner ID")
		zap.L().Error("Invalid owner ID", zap.String("ownerId", ownerId))
		return
	}

	// Validate api key
	if !apiKeyValid(apiKey, ownerId) {
		c.String(consts.StatusBadRequest, "Invalid API key")
		zap.L().Error("Invalid API key", zap.String("apiKey", apiKey))
		return
	}

	// Getting the request Body
	reqBody, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, "Request body is missing"+err.Error())
		zap.L().Error("Request body is missing", zap.Error(err))
		return
	}
	buf := bytes.NewBuffer(reqBody)

	// Decode the JSON request
	err = json.NewDecoder(buf).Decode(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body"+err.Error())
		zap.L().Error("Failed to parse request body", zap.Error(err))
		return
	}

	zap.L().Info("Request body decoded successfully")

	// Prepare the updated admin info
	adminConfig := repository.AdminConfig{
		ApiKey:    apiKey,
		OwnerName: req[0]["OwnerName"].(string),
		OwnerId:   req[0]["OwnerId"].(string),
	}

	zap.L().Info("Updated admin info prepared successfully")

	// Unmarshal the Services field from the request into the Services field of adminConfig
	servicesJSON, err := json.Marshal(req[0]["Services"])
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse Services field")
		zap.L().Error("Failed to parse Services field", zap.Error(err))
		return
	}
	err = json.Unmarshal(servicesJSON, &adminConfig.Services)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse Services field"+err.Error())
		zap.L().Error("Failed to parse Services field", zap.Error(err))
		return
	}

	err = repository.UpdateAdminInfo(ownerId, adminConfig)
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to update admin")
		zap.L().Error("Failed to update admin", zap.Error(err))
		return
	}

	zap.L().Info("Admin info updated successfully")

	go cache.UpdateIDLcache()

	response := make(map[string]string)
	response["Message"] = "Updated successfully. You're good to GO :D"

	c.JSON(consts.StatusOK, response)
}
