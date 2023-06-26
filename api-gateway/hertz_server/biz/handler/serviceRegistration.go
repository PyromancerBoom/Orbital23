package handler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var servicesMap = make(map[string][]map[string]interface{})

type registerRequest struct {
	ServiceOwner string `json:"serviceOwner"`
	// Path              string                   `json:"serviceOwner"`
	RegisteredServers []map[string]interface{} `json:"registeredServers"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	// Parse the request payload
	var req registerRequest

	// Read the request body
	reqBody, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, "Request body is missing")
		return
	}

	buf := bytes.NewBuffer(reqBody)

	err = json.NewDecoder(buf).Decode(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body")
		return
	}
	servicesMap[req.ServiceOwner] = req.RegisteredServers

	c.String(consts.StatusOK, "Services registered successfully")
}

// displayAll returns the hashmap with all the stored details.
func DisplayAll(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, servicesMap)

}
