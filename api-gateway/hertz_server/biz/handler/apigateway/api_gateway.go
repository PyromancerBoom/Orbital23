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

	// Print request data
	fmt.Println("Request Data:")
	fmt.Println(c.Param("requestData"))

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

	// Make generic Call and get back response
	responseJson, err := genClient.GenericCall(ctx, serviceMethod, jsonString)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	c.JSON(consts.StatusOK, responseJson)
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
		return
	}

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
		client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	// Make JSON string from request
	requestJSON, err := json.Marshal(req)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	// Make generic call and get back response as map[string]interface{}
	responseData, err := genClient.GenericCall(ctx, serviceMethod, string(requestJSON))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	// Convert response data to JSON string
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	c.String(consts.StatusOK, string(responseJSON))
}
