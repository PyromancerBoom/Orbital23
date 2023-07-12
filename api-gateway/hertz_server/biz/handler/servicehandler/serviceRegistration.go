package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	repository "api-gateway/hertz_server/biz/model/repository"
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

	// Register information in Mongo DB along with the apiKey
	for _, data := range req {
		ownerId := data["ownerId"].(string)
		if isAlreadyRegistered(ownerId) {
			c.String(consts.StatusBadRequest, "Owner ID already exists")
			return
		}
		dataMap := bson.M{
			"ownerId": ownerId,
			"apiKey":  apiKey,
		}
		err = repository.StoreClientData(dataMap)
		if err != nil {
			c.String(consts.StatusInternalServerError, "Failed to register client data")
			return
		}
	}

	res := make(map[string]string)
	res["Message"] = "Registered successfully. You're good to GO :D"
	res["api-key"] = apiKey

	c.JSON(consts.StatusOK, res)
}

// Function to check if owner ID already exists using the ownerIds map
func isAlreadyRegistered(ownerId string) bool {
	clientData, err := repository.GetClientDataByOwnerID(ownerId)
	if err != nil {
		return false
	}
	return clientData != nil
}

func DisplayAll(ctx context.Context, c *app.RequestContext) {
	clientDataList, err := repository.GetAllClientData()
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to fetch client data")
		return
	}

	c.JSON(consts.StatusOK, clientDataList)
}
