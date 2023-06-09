package asset_api

import (
	"context"
	"fmt"

	asset_api "api-gateway/hertz_server/biz/model/asset_api"

	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Generic client initialisation
func NewGenericClient(destServiceName string) genericclient.Client {
	genericCli := genericclient.NewClient(destServiceName, generic.BinaryThriftGeneric())
	return genericCli
}

// QueryStudent .
// @router asset/query [GET]
func QueryAsset(ctx context.Context, c *app.RequestContext) {
	var err error
	var req asset_api.QueryAssetRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// Use generic client
	// genericclient.NewClient in func NewGenericClient handles all errors
	client := NewGenericClient()

	// The resulting reqRpc object represents the request to be
	// sent to the "QueryAsset" method, with the
	// asset ID provided as a string value.
	// It allows the client to specify which asset it wants
	// to query in the remote "asset" service.
	requestRpc := &QueryAssetRequest{
		ID: fmt.Sprintf("%d", req.ID),
	}

	// Responsible for making an RPC Call
	// to the "QueryAsset" method of the "asset" service
	// using the Kitex client.
	responseRpc, err := client.QueryAsset(ctx, requestRpc)
	if err != nil {
		panic(err)
	}

	if !responseRpc.Exist {
		resp := &asset_api.QueryAssetResponse{
			Msg: fmt.Sprintf("No data for asset ID: %d", req.ID),
		}
		// The response object (resp) is serialized as JSON
		// using c.JSON(200, resp) and sent as the HTTP response.
		c.JSON(consts.StatusOK, resp)
		return
	}

	// Constructs response for client
	// A new instance of asset_api.QueryAssetResponse is
	// created, and its fields ID, Name, and Market are
	// assigned values from the corresponding fields of
	// the respRpc response object.
	resp := &asset_api.QueryAssetResponse{
		ID:     responseRpc.ID.(int),
		Name:   responseRpc.Name.(string),
		Market: responseRpc.Market.(string),
	}

	c.JSON(consts.StatusOK, resp)
}

// InsertAsset is an HTTP handler function for handling the "asset/insert" route with a POST method.
// It receives an asset_api.InsertAssetRequest object from the client, validates the request,
// and performs an RPC call to the "InsertAsset" method of the "asset" service using a generic client.
// The response from the RPC call is converted into an asset_api.InsertAssetResponse object
// and sent back to the client as an HTTP response.
func InsertAsset(ctx context.Context, c *app.RequestContext) {
	// Initialize variables
	var err error
	var req asset_api.InsertAssetRequest

	// Bind and validate the request data received from the client
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	fmt.Println("Received Insert Request")

	// Use a generic client to make the RPC call
	client := NewGenericClient("asset")

	// Create the request object to be sent to the "InsertAsset" method
	requestRpc := &InsertAssetRequest{
		ID:     req.ID,
		Name:   req.Name,
		Market: req.Market,
	}

	// Make an RPC call to the "InsertAsset" method using the generic client
	responseRpc, err := client.InsertAsset(ctx, requestRpc)
	if err != nil {
		panic(err)
	}

	// Check if the RPC call was successful
	if !responseRpc.Ok {
		// If the response indicates an error, construct the error response
		resp := asset_api.InsertAssetResponse{
			OK:  false,
			Msg: responseRpc.Msg.(string),
		}
		c.JSON(200, resp)
		return
	}

	// Construct the successful response
	resp := asset_api.InsertAssetResponse{
		OK:  true,
		Msg: responseRpc.Msg.(string),
	}

	c.JSON(200, resp)
}
