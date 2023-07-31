// Code generated by hertz generator.

package apigateway

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	apigateway "api-gateway/hertz_server/biz/model/apigateway"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"go.uber.org/zap"

	"github.com/cloudwego/kitex/client/callopt"

	cache "api-gateway/hertz_server/biz/model/cache"
)

// ProcessPostRequest .
// @router /{:serviceName}/{:serviceMethod} [POST]
func ProcessPostRequest(ctx context.Context, c *app.RequestContext) {
	var err error
	options := []callopt.Option{callopt.WithRPCTimeout(time.Second * 100), callopt.WithConnectTimeout(time.Millisecond * 150)}

	zap.L().Info("Processing POST request")

	// Parsing and validation
	var req apigateway.GatewayRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		zap.L().Error("Error in parsing and validating request", zap.Error(err))
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	serviceName := c.Param("serviceName")
	path := c.Param("path")

	reqBody, err := c.Body()
	if err != nil {
		zap.L().Error("Error while getting request body", zap.Error(err))
		c.String(consts.StatusBadRequest, "Request body is missing")
		return
	}

	trimmedReqBody := bytes.TrimSpace(reqBody)
	if len(trimmedReqBody) == 0 {
		zap.L().Warn("Request body is empty")
		c.String(consts.StatusBadRequest, "Request body is empty")
		return
	}

	jsonString := string(reqBody)

	//This throws an error if the serviceName and path does not exist.
	genClient, err := cache.GetGenericClient(serviceName, path)
	if err != nil {
		zap.L().Error("Error in fetching client", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	methodName, err := cache.GetExposedMethodFromPath(path)
	if err != nil {
		zap.L().Error("Error fetching method name.", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// Make generic Call and get back response
	zap.L().Info("Making generic Call and getting back response")
	response, err := genClient.GenericCall(ctx, methodName, jsonString, options...)
	if err != nil {
		zap.L().Error("Error while making generic call", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	c.String(consts.StatusOK, response.(string))
	zap.L().Info("(POST) Execution complete")
}

// ProcessGetRequest handles the GET request.
// @router /{:serviceName}/{:serviceMethod} [GET]
func ProcessGetRequest(ctx context.Context, c *app.RequestContext) {
	var err error
	options := []callopt.Option{callopt.WithRPCTimeout(time.Second * 100), callopt.WithConnectTimeout(time.Millisecond * 150)}

	zap.L().Info("Processing GET request")

	// Parsing and validation
	var req apigateway.GatewayRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		zap.L().Error("Error in parsing and validating request", zap.Error(err))
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	serviceName := c.Param("serviceName")
	path := c.Param("path")

	//This throws an error if the serviceName and path does not exist.
	genClient, err := cache.GetGenericClient(serviceName, path)
	if err != nil {
		zap.L().Error("Error in fetching client", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	queryParams := c.QueryArgs()

	// Extract query parameters
	params := make(map[string]interface{})
	queryParams.VisitAll(func(key, value []byte) {
		params[string(key)] = string(value)
	})
	zap.L().Info("QueryArgs fetched", zap.Any("Params: ", params))

	// Make Json string from request
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		zap.L().Error("Error while marshalling json", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}
	jsonString := string(jsonBytes)

	methodName, err := cache.GetExposedMethodFromPath(path)
	if err != nil {
		zap.L().Error("Error fetching method name.", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// Make generic Call and get back response using WithRPC Timeout
	zap.L().Info("Making generic Call and getting back response")
	response, err := genClient.GenericCall(ctx, methodName, jsonString, options...)
	if err != nil {
		zap.L().Error("Error while making generic call", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	c.String(consts.StatusOK, response.(string))
	zap.L().Info("(GET) Execution complete")
}
