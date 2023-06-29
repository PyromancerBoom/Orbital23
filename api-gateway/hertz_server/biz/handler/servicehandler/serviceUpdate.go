package servicehandler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func Update(ctx context.Context, c *app.RequestContext) {
	// // we get the key from header
	// apiKey := string(c.GetHeader("apikey"))

	// // now we get the body
	// reqBody, err := c.Body()
	// if err != nil {
	// 	c.String(consts.StatusBadRequest, "Request body error")
	// 	return
	// }

	// // and now we check if api key exists
	// service, ok := servicesMap[apiKey]
	// if !ok {
	// 	c.String(consts.StatusBadRequest, "Invalid API key")
	// 	return
	// }

	// // for updation, we make a temp map and assign it to the registered one
	// var updatedServiceTemp Service
	// err = json.Unmarshal(reqBody, &updatedServiceTemp)
	// if err != nil {
	// 	c.String(consts.StatusBadRequest, "Failed to parse request body")
	// 	return
	// }

	// service.ServiceOwner = updatedServiceTemp.ServiceOwner
	// service.RegisteredServers = updatedServiceTemp.RegisteredServers

	// servicesMap[apiKey] = service

	// // Sending back a string respone if everything goes well
	// c.String(consts.StatusOK, "Service updated successfully. You're good to \"GO\" :D")
}
