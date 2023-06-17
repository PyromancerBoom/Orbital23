// Code generated by hertz generator.

package apigateway

import (
	"context"
	"encoding/json"
	"fmt"

	apigateway "api-gateway/hertz_server/biz/model/apigateway"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	client "github.com/cloudwego/kitex/client"
	genericClient "github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"

	idlmap "api-gateway/hertz_server/biz/model/idlmap"
)

// ProcessPostRequest .
// @router /{:serviceName}/{:serviceMethod} [POST]
func ProcessPostRequest(ctx context.Context, c *app.RequestContext) {
	var err error
	var req apigateway.GatewayRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(" ")
	fmt.Println("Reached Here POST")

	serviceName := c.Param("serviceName")
	serviceMethod := c.Param("serviceMethod")

	fmt.Printf("Received generic POST request for service '%s' method '%s'\n", serviceName, serviceMethod)

	// Checking if service and method are valid
	idl, err := idlmap.GetIdlFile(serviceName, serviceMethod)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	fmt.Printf("IDL path '%s'\n", idl)

	// Generic client init
	provider, err := generic.NewThriftFileProvider(idl)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
		client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// Make Json string from request
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	jsonString := string(jsonBytes)

	responseJson, err := genClient.GenericCall(ctx, serviceMethod, jsonString)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// Convert responseJson to []byte
	responseBytes, err := json.Marshal(responseJson)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())

	}

	// Unmarshal responseBytes
	var response apigateway.GatewayResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// Create the final response object
	finalResponse := apigateway.GatewayResponse{
		StatusCode:   response.StatusCode,
		ResponseData: string(responseBytes),
	}

	fmt.Println("Received response from backend service for POST")

	c.JSON(consts.StatusOK, finalResponse)
}

// ProcessGetRequest .
// @router /{:serviceName}/{:serviceMethod} [GET]
func ProcessGetRequest(ctx context.Context, c *app.RequestContext) {
	var err error
	var req apigateway.GatewayRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(" ")
	fmt.Println("Reached Here GET")

	serviceName := c.Param("serviceName")
	serviceMethod := c.Param("serviceMethod")

	fmt.Printf("Received generic GET request for service '%s' method '%s'\n", serviceName, serviceMethod)

	// Checking if service and method are valid
	idl, err := idlmap.GetIdlFile(serviceName, serviceMethod)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())

	}

	fmt.Printf("IDL path '%s'\n", idl)

	// Generic client init
	provider, err := generic.NewThriftFileProvider(idl)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())

	}

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())

	}

	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
		client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())

	}

	// Make Json string from request
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())

	}

	jsonString := string(jsonBytes)

	responseJson, err := genClient.GenericCall(ctx, serviceMethod, jsonString)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())

	}

	// Convert responseJson to []byte
	responseBytes, err := json.Marshal(responseJson)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// Unmarshal responseBytes
	var response apigateway.GatewayResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		panic(err)
	}

	// Create the final response object
	finalResponse := apigateway.GatewayResponse{
		StatusCode:   response.StatusCode,
		ResponseData: string(responseBytes),
	}

	fmt.Println("Received response from backend service for GET")

	c.JSON(consts.StatusOK, finalResponse)
}
