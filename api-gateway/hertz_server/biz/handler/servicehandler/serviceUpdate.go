package servicehandler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// To be updated (Unintentional pun) later when registration feature is finalised
func Update(ctx context.Context, c *app.RequestContext) {
	// we get the key from header
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

	// ------------- LOGIC FOR SERVICE UPDATE -----------------

	// Sending back a string respone if everything goes well
	c.String(consts.StatusOK, "Service updated successfully. You're good to \"GO\" :D")
}
