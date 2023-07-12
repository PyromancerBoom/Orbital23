package servicehandler

import (
	"bytes"
	"context"
	"encoding/json"

	"api-gateway/hertz_server/biz/model/repository"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

var db *repository.Database

func Init() error {
	var err error
	db, err = repository.ConnectDB("http://localhost:32769/", "testDB", "testCollection")
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}

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
		clientData := &repository.ClientData{
			OwnerId: ownerId,
			ApiKey:  apiKey,
		}

		// Use the db instance to call the StoreClientData method
		err := db.StoreClientData(clientData)
		if err != nil {
			c.String(consts.StatusInternalServerError, "Failed to store client data")
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
	clientData, err := db.GetClientDataByOwnerID(ownerId)
	if err != nil {
		return false
	}
	return clientData != nil
}

func DisplayAll(ctx context.Context, c *app.RequestContext) {
	clientDataList, err := db.GetAllClientData()
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to fetch client data")
		return
	}

	c.JSON(consts.StatusOK, clientDataList)
}
