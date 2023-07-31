package servicehandler

import (
	"api-gateway/hertz_server/biz/model/cache"
	"bytes"
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
)

type HealthRequest struct {
	ServerID string `json:"ServerID"`
	ApiKey   string `json:"ApiKey"`
}

// HealthCheck is the handler for health check requests
// @Route = /health
func HealthCheck(ctx context.Context, c *app.RequestContext) {

	serverIsHealthy, err := checkIfServiceHealthy("RegistryProxy")
	if err != nil {
		zap.L().Warn(err.Error())
	}

	//even if we get an error from checking server health, we can still deal with it by trying to perform health request by the gateway itself
	if (err != nil) || (!serverIsHealthy) {
		zap.L().Info("Performing health check request.")
		performHealthCheckRequest(ctx, c)
	} else {
		zap.L().Info("Proxying health check request.")
		proxyHealthCheckRequst(ctx, c)
	}

}

func proxyHealthCheckRequst(ctx context.Context, c *app.RequestContext) {
	reqBody, err := c.Body()

	if err != nil {
		zap.L().Error("Error while getting request body", zap.Error(err))
		c.String(consts.StatusBadRequest, "Request body is missing.")
		return
	}

	trimmedReqBody := bytes.TrimSpace(reqBody)
	if len(trimmedReqBody) == 0 {
		zap.L().Warn("Request body is empty")
		c.String(consts.StatusBadRequest, "Request body is empty.")
		return
	}

	jsonString := string(reqBody)

	genClient, err := cache.GetGenericClient("RegistryProxy", "healthCheckServer")
	if err != nil {
		zap.L().Error(err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	response, err := genClient.GenericCall(ctx, "healthCheckServer", jsonString, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		zap.L().Error("Error while making generic call", zap.Error(err))

		if err.Error() == "service discovery error: no service found" {
			c.String(consts.StatusInternalServerError, "Server connection services are currently down.")
		} else {
			c.String(consts.StatusInternalServerError, err.Error())
		}
		return
	}

	c.String(consts.StatusOK, response.(string))
}

func performHealthCheckRequest(ctx context.Context, c *app.RequestContext) {
	var req HealthRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		zap.L().Error(err.Error())
		c.String(consts.StatusExpectationFailed, "Invalid Request.")
		return
	}

	if !authoriseHealthCheckRight(req.ApiKey, req.ServerID) {
		zap.L().Info("Server " + req.ServerID + "Unauthorized.")
		c.String(consts.StatusUnauthorized, "Unauthorized.")
		return
	}

	err = updateAsHealthy(req.ServerID)
	if err != nil {
		zap.L().Error(err.Error())
		c.String(consts.StatusInternalServerError, "Unable to process health update request.")
		return
	}

	res := make(map[string]string)
	res["Status"] = "ok"
	res["Message"] = "Successfully updated server health."

	c.JSON(consts.StatusOK, res)
}
