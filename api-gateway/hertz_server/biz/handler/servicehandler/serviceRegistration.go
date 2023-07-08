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
var servicesMap map[string]Service

func Register(ctx context.Context, c *app.RequestContext) {
	var req []struct {
		Service Service `json:"Service"`
	}
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

	for _, item := range req {
		service := item.Service

		if isAlreadyRegistered(service.ID) {
			c.String(consts.StatusBadRequest, "Already registered")
			return
		}

		apiKey := uuid.New().String()
		service.Meta.APIKey = apiKey

		// stor the service information
		if servicesMap == nil {
			servicesMap = make(map[string]Service)
		}

		servicesMap[apiKey] = service
	}

	res := make(map[string]string)
	res["Message"] = "Registered successfully. You're good to \"GO\" :D"

	c.JSON(consts.StatusOK, res)
}

func DisplayAll(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, servicesMap)
}

func isAlreadyRegistered(ownerId string) bool {
	for _, service := range servicesMap {
		if service.ID == ownerId {
			return true
		}
	}
	return false
}
