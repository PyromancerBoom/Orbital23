package servicehandler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func Update(ctx context.Context, c *app.RequestContext) {
	// // Get the API key from the header
	// apiKey := string(c.GetHeader("apikey"))
	// // Check if the API key exists
	// clientData, ok := servicesMap[apiKey]
	// if !ok {
	// 	c.String(consts.StatusBadRequest, "Invalid API key")
	// 	return
	// }

	// // Get the request body
	// reqBody, err := c.Body()
	// if err != nil {
	// 	c.String(consts.StatusBadRequest, "Request body error")
	// 	return
	// }

	// // Decode the JSON request
	// var updatedData ClientData
	// err = json.Unmarshal(reqBody, &updatedData)
	// if err != nil {
	// 	c.String(consts.StatusBadRequest, "Failed to parse request body")
	// 	return
	// }

	// // Update the client data with the new values
	// clientData.OwnerName = updatedData.OwnerName
	// clientData.OwnerId = updatedData.OwnerId
	// clientData.Services = updatedData.Services

	// // Update the services map
	// servicesMap[apiKey] = clientData

	// // Sending back a string response if everything goes well
	// c.String(consts.StatusOK, "Service updated successfully. You're good to \"GO\" :D")
}
