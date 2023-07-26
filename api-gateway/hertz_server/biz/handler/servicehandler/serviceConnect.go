package servicehandler

//handles :/connect endpoint to register a server

//authorise if the api-key is valid

//register the server.

//request here has to have
//1: API KEY
//2: ServiceName
//3: ServerAddress

import (
	"api-gateway/hertz_server/biz/model/cache"
	"bytes"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Request struct {
	ApiKey        string `json:"ApiKey"`
	ServiceName   string `json:"ServiceName"`
	ServerAddress string `json:"ServerAddress"`
	ServerPort    string `json:"ServerPort"`
}

// handler for /connect endpoint
func Connect(ctx context.Context, c *app.RequestContext) {

	serverIsHealthy, err := checkIfServiceHealthy("RegistryProxy")
	if err != nil {
		zap.L().Error(err.Error())
	}

	//even if we get an error from checking server health, we can still deal with it by trying to perform connection request by the gateway itself
	if (err != nil) || (!serverIsHealthy) {
		zap.L().Info("Performing server connection request.")
		performServerConnectionRequest(ctx, c)
	} else {
		zap.L().Info("Proxying server connection request.")
		proxyServerConnectionRequest(ctx, c)
	}

}

func proxyServerConnectionRequest(ctx context.Context, c *app.RequestContext) {
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

	genClient, err := cache.GetGenericClient("RegistryProxy")
	if err != nil {
		zap.L().Error(err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	jsonString := string(reqBody)

	response, err := genClient.GenericCall(ctx, "connectServer", jsonString)
	if err != nil {
		zap.L().Error("Error while making generic call", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	c.String(consts.StatusOK, response.(string))
}

func performServerConnectionRequest(ctx context.Context, c *app.RequestContext) {
	var req Request
	err := c.BindAndValidate(&req)
	if err != nil {
		zap.L().Error(err.Error())
		c.String(consts.StatusExpectationFailed, "Invalid Request.")
		return
	}

	if !authoriseConnect(req.ApiKey, req.ServiceName) {
		zap.L().Info(req.ApiKey + "unauthorized for " + req.ServiceName + "service.")
		c.String(consts.StatusUnauthorized, "Unauthorized.")
		return
	}

	if validateAddress(req.ServerAddress, req.ServerPort) != nil {
		zap.L().Info("Invalid address" + req.ServerAddress + "for server connection.")
		c.String(consts.StatusBadRequest, "Address is invalid.")
		return
	}

	res := make(map[string]string)
	res["Status"] = "ok"
	res["Message"] = "Server Connection Request Accepted."
	res["ServerID"] = uuid.New().String()

	err = registerServer(req.ServerAddress, req.ServerPort, res["ServerID"], req.ServiceName, req.ApiKey)
	if err != nil {
		c.String(consts.StatusInternalServerError, "Unable to connect server.")
		return
	}

	c.JSON(consts.StatusOK, res)
}
