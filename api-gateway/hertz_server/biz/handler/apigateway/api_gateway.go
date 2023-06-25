// Code generated by hertz generator.

package apigateway

import (
	"bytes"
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

// ------------------------------------- NOTE -------------------------------------
// POST Request working
// GET requests have some issue (ProcessGetRequest Method)
// --------------------------------------------------------------------------

// ProcessPostRequest .
// @router /{:serviceName}/{:serviceMethod} [POST]
func ProcessPostRequest(ctx context.Context, c *app.RequestContext) {
	var err error

	// Parsing and validation
	var req apigateway.GatewayRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// fmt.Println(" ")
	// fmt.Println("Reached Here POST")

	serviceName := c.Param("serviceName")
	path := c.Param("path")

	// fmt.Printf("Received generic POST request for service '%s' method '%s'\n", serviceName, serviceMethod)

	reqBody, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, "Request body is missing")
		return
	}

	trimmedReqBody := bytes.TrimSpace(reqBody)
	if len(trimmedReqBody) == 0 {
		c.String(consts.StatusBadRequest, "Request body is empty")
		return
	}

	// // Print request data
	// fmt.Println("Request data:")
	// fmt.Println(string(reqBody))

	// Checking if service and method are valid
	value, err := idlmap.GetService(serviceName, path)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	// fmt.Printf("IDL path '%s'\n", idl)

	// provider initialisation
	provider, err := generic.NewThriftContentProvider(value.IDL, nil)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nProvider Init error \n")
		return
	}

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nJSONThriftGeneric error \n")
		return
	}

	// fetch hostport from registry later
	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
		client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nGeneric client initialisation error \n")
	}

	// Make Json string from request
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nJson Marshalling error \n")
	}

	jsonString := string(jsonBytes)

	// Make generic Call and get back response
	response, err := genClient.GenericCall(ctx, value.Method, jsonString)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nGeneric call error \n")
	}

	c.String(consts.StatusOK, response.(string))
}

// ProcessGetRequest handles the GET request.
// @router /{:serviceName}/{:serviceMethod} [GET]
func ProcessGetRequest(ctx context.Context, c *app.RequestContext) {
	var err error
	// Parsing and validation
	var req apigateway.GatewayRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(" ")
	fmt.Println("Reached Here GET")

	serviceName := c.Param("serviceName")
	path := c.Param("path")

	fmt.Printf("Received generic GET request for service '%s' method '%s'\n", serviceName, path)

	// Checking if service and method are valid
	value, err := idlmap.GetService(serviceName, path)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	}

	fmt.Printf("IDL path '%s'\n", value.IDL)

	// provider initialisation
	provider, err := generic.NewThriftContentProvider(value.IDL, nil)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nProvider Init error \n")
		return
	}

	thriftGeneric, err := generic.JSONThriftGeneric(provider)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nJSONThriftGeneric error \n")
		return
	}

	// fetch hostport from registry later
	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
		client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nGeneric client initialisation error \n")
	}

	queryParams := c.QueryArgs()

	// Make Json string from request
	jsonBytes, err := json.Marshal(queryParams)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nJson Marshalling error \n")
	}

	jsonString := string(jsonBytes)

	// Make generic Call and get back response
	response, err := genClient.GenericCall(ctx, value.Method, jsonString)
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error()+"\nGeneric call error \n")
	}

	c.String(consts.StatusOK, response.(string))
}

// func ProcessGetRequest(ctx context.Context, c *app.RequestContext) {
// 	serviceName := c.Param("serviceName")
// 	path := c.Param("path")

// 	// Checking if service and method are valid
// 	value, err := idlmap.GetIdlFile(serviceName, path)
// 	if err != nil {
// 		c.String(consts.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// Provider initialization
// 	provider, err := generic.NewThriftFileProvider(value.IDL)
// 	if err != nil {
// 		c.String(consts.StatusInternalServerError, err.Error()+"\nProvider Init error \n")
// 		return
// 	}

// 	thriftGeneric, err := generic.JSONThriftGeneric(provider)
// 	if err != nil {
// 		c.String(consts.StatusInternalServerError, err.Error()+"\nJSONThriftGeneric error \n")
// 		return
// 	}

// 	// Fetch hostport from registry later
// 	genClient, err := genericClient.NewClient(serviceName, thriftGeneric,
// 		client.WithHostPorts("127.0.0.1:8888"))
// 	if err != nil {
// 		c.String(consts.StatusInternalServerError, err.Error()+"\nGeneric client initialization error \n")
// 		return
// 	}

// 	// Get the 'id' query parameter from the request URL
// 	id := c.Query("id")

// 	// Create a query request object
// 	queryRequest := map[string]interface{}{
// 		"ID": id,
// 	}
// 	requestBody, err := json.Marshal(queryRequest)
// 	if err != nil {
// 		c.String(consts.StatusInternalServerError, err.Error()+"\nJSON Marshalling error \n")
// 		return
// 	}

// 	// Make the generic call and get the response
// 	response, err := genClient.GenericCall(ctx, value.Method, string(requestBody))
// 	if err != nil {
// 		c.String(consts.StatusInternalServerError, err.Error()+"\nGeneric call error \n")
// 		return
// 	}

// 	c.String(consts.StatusOK, response.(string))
// }
