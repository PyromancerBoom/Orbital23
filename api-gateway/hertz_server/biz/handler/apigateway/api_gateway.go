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

	consul "github.com/kitex-contrib/registry-consul"

	"go.uber.org/zap"

	client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	genericClient "github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/generic"

	cache "api-gateway/hertz_server/biz/model/cache"
)

var resolver discovery.Resolver

// init is called during package initialization and sets up the resolver.
func init() {
	// Get registry to enable resolving serverIDs
	var err error
	resolver, err = consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		zap.L().Error("Error while getting registry", zap.Error(err))
	}
}

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

	// Checking if service and method are valid
	value, err := cache.GetServiceDetails(serviceName, path)
	if err != nil {
		zap.L().Error("Error while getting service", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}
	zap.L().Info("Checked that service and method are valid")
	// Print the service using zap
	zap.L().Info("Service: ", zap.Any("Details: ", value))

	// provider initialisation
	provider, err := generic.NewThriftContentProvider(value.IDL, nil)
	if err != nil {
		zap.L().Error("Error while initializing provider", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	zap.L().Info("Provider initialised")

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		zap.L().Error("Error while creating JSONThriftGeneric", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	// Fetch hostport from registry later
	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
		client.WithResolver(resolver))
	if err != nil {
		zap.L().Error("Error while initializing generic client", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	jsonString := string(reqBody)

	// Make generic Call and get back response
	zap.L().Info("Making generic Call and getting back response")
	response, err := genClient.GenericCall(ctx, value.Method, jsonString, options...)
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

	// Checking if service and method are valid
	value, err := cache.GetServiceDetails(serviceName, path)
	if err != nil {
		zap.L().Error("Error while getting service", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// Provider initialisation
	provider, err := generic.NewThriftContentProvider(value.IDL, nil)
	if err != nil {
		zap.L().Error("Error while initializing provider", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		zap.L().Error("Error while creating JSONThriftGeneric", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	// Fetch hostport from registry later
	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
		client.WithResolver(resolver))
	if err != nil {
		zap.L().Error("Error while initializing generic client", zap.Error(err))
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

	// Make generic Call and get back response using WithRPC Timeout
	zap.L().Info("Making generic Call and getting back response")
	response, err := genClient.GenericCall(ctx, value.Method, jsonString, options...)
	if err != nil {
		zap.L().Error("Error while making generic call", zap.Error(err))
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// response, err := genClient.GenericCall(ctx, value.Method, jsonString)
	// if err != nil {
	// 	zap.L().Error("Error while making generic call", zap.Error(err))
	// 	c.String(consts.StatusInternalServerError, err.Error())
	// }

	c.String(consts.StatusOK, response.(string))
	zap.L().Info("(GET) Execution complete")
}
