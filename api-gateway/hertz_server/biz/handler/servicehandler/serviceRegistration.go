package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	repository "api-gateway/hertz_server/biz/model/repository"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler for "/register"
// Registers info sent through a post request
// Generates API key for admin
func Register(ctx context.Context, c *app.RequestContext) {
	var req []map[string]interface{}

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

	// Check if ownerId exists
	if ownerIdExists(req[0]["OwnerId"].(string)) {
		c.String(consts.StatusBadRequest, "Owner ID already exists")
		zap.L().Error("Owner ID already exists", zap.String("ownerId", req[0]["OwnerId"].(string)))
		return
	}

	apiKey := uuid.New().String()
	zap.L().Info("API key generated successfully", zap.String("apiKey", apiKey))

	// Store the admin info in the database
	adminConfig := repository.AdminConfig{
		ApiKey:    apiKey,
		OwnerName: req[0]["OwnerName"].(string),
		OwnerId:   req[0]["OwnerId"].(string),
	}

	zap.L().Info("Admin info prepared successfully")

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

	err = repository.StoreAdminInfo(adminConfig)
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to register admin")
		zap.L().Error("Failed to register admin", zap.Error(err))
		return
	}

	zap.L().Info("Admin info stored successfully")

	response := make(map[string]string)
	response["Message"] = "Registered successfully. You're good to GO :D"
	response["api-key"] = apiKey

	c.JSON(consts.StatusOK, response)
}
