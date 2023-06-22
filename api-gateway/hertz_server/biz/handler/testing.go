// I made this custom handler to test out if we can add new routes.
// and since this code worked. We can use similar logic as below to register services in service registry

package handler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var dataMap = make(map[string]map[string]string)

type registerRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// Register stores the provided ID, name, and path in a hashmap.
func Register(ctx context.Context, c *app.RequestContext) {
	// Parse the request payload
	var req registerRequest

	// Read the request body
	reqBody, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, "Request body is missing")
		return
	}

	// Use a buffer to create an io.Reader for json.NewDecoder
	buf := bytes.NewBuffer(reqBody)

	// Decode the JSON request
	err = json.NewDecoder(buf).Decode(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Store the details in the hashmap
	if _, ok := dataMap[req.ID]; !ok {
		dataMap[req.ID] = make(map[string]string)
	}
	dataMap[req.ID]["name"] = req.Name
	dataMap[req.ID]["path"] = req.Path

	c.String(consts.StatusOK, "Registered successfully")
}

// displayAll returns the hashmap with all the stored details.
func DisplayAll(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, dataMap)

}
