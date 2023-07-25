package servicehandler

import (
	repository "api-gateway/hertz_server/biz/model/repository"
	"bytes"
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	genericClient "github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

type HealthRequest struct {
	ServerID string `json:"ServerID"`
	APIKey   string `json:"ApiKey"`
}

// handler for /health endpoint
func HealthCheck(ctx context.Context, c *app.RequestContext) {

	serverIsHealthy, err := checkIfServiceHealthy("RegistryProxy")
	if err != nil {
		zap.L().Error(err.Error())
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

	// Checking if service is valid
	idl, err := repository.GetServiceIDL(healthCheckServiceName)
	if err != nil {
		zap.L().Error("Error getting Registry Proxy IDL", zap.Error(err))
		c.String(consts.StatusInternalServerError, "Unable to perform health check.")
		return
	}

	// provider initialisation
	provider, err := generic.NewThriftContentProvider(idl, nil)
	if err != nil {
		zap.L().Error("Error while initializing provider for Registry Proxy", zap.Error(err))
		c.String(consts.StatusInternalServerError, "Unable to perform health check.")
		return
	}

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		zap.L().Error("Error while creating JSONThriftGeneric", zap.Error(err))
		c.String(consts.StatusInternalServerError, "Unable to perform health check.")
		return
	}

	// Fetch hostport from registry later
	genClient, err := genericClient.NewClient(healthCheckServiceName, thriftGeneric,
		client.WithResolver(registryResolver))
	if err != nil {
		zap.L().Error("Error while initializing generic client", zap.Error(err))
		c.String(consts.StatusInternalServerError, "Unable to perform health check.")
		return
	}

	jsonString := string(reqBody)

	// Make generic Call and get back response
	zap.L().Info("Making generic Call and getting back response")
	response, err := genClient.GenericCall(ctx, healthcheckMethodName, jsonString, callopt.WithRPCTimeout(5*time.Second))
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

	if !authoriseHealthCheckRight(req.APIKey, req.ServerID) {
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
	res["Status"] = "status OK"
	res["Message"] = "Successfully Updated the heatlh of server"

	c.JSON(consts.StatusOK, res)
}
