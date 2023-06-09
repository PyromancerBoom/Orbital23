package asset_api

import (
	"context"
	"fmt"

	asset_api "api-gateway/hertz_server/biz/model/asset_api"
	"api-gateway/hertz_server/kitex_gen/asset_management"
	"api-gateway/hertz_server/kitex_gen/asset_management/assetmanagement"

	clientK "github.com/cloudwego/kitex/client"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// QueryAsset handles the GET request to query asset information.
// @router asset/query [GET]
func QueryAsset(ctx context.Context, c *app.RequestContext) {
	var err error

	// Parse and validate the request
	var req asset_api.QueryAssetRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// Print a message indicating that a query request has been received
	fmt.Println("Received Query Request")

	// Create a Kitex client to make an RPC call to the asset management service
	client, err := assetmanagement.NewClient("asset", clientK.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		panic(err)
	}

	// Create an RPC request with the asset ID from the parsed request
	reqRpc := &asset_management.QueryAssetRequest{
		ID: fmt.Sprintf("%d", req.ID),
	}

	// Make an RPC call to the asset management service to query asset information
	respRpc, err := client.QueryAsset(ctx, reqRpc)
	if err != nil {
		panic(err)
	}

	// Check if the asset exists based on the RPC response
	if !respRpc.Exist {
		// If the asset doesn't exist, create a response with a corresponding message
		resp := &asset_api.QueryAssetResponse{
			Msg: fmt.Sprintf("No data for asset ID: %d", req.ID),
		}
		c.JSON(200, resp)
		return
	}

	// Create a response with the asset information from the RPC response
	resp := &asset_api.QueryAssetResponse{
		ID:     respRpc.ID,
		Name:   respRpc.Name,
		Market: respRpc.Market,
	}

	c.JSON(consts.StatusOK, resp)
}

// InsertAsset handles the POST request to insert a new asset.
// @router asset/insert [POST]
func InsertAsset(ctx context.Context, c *app.RequestContext) {
	var err error

	// Parse and validate the request
	var req asset_api.InsertAssetRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	// Print a message indicating that an insert request has been received
	fmt.Println("Received Insert Request")

	// Create a Kitex client to make an RPC call to the asset management service
	client, err := assetmanagement.NewClient("asset", clientK.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		panic(err)
	}

	// Create an RPC request with the asset information from the parsed request
	reqRpc := &asset_management.InsertAssetRequest{
		ID:     req.ID,
		Name:   req.Name,
		Market: req.Market,
	}

	// Make an RPC call to the asset management service to insert the asset
	respRpc, err := client.InsertAsset(ctx, reqRpc)
	if err != nil {
		panic(err)
	}

	// Check if the asset insertion was successful based on the RPC response
	if !respRpc.Ok {
		// If the insertion failed, create a response with the failure status and the error message
		resp := asset_api.InsertAssetResponse{
			OK:  false,
			Msg: respRpc.Msg,
		}
		c.JSON(200, resp)
		return
	}

	// If the insertion was successful, create a response with the success status and the success message
	resp := asset_api.InsertAssetResponse{
		OK:  true,
		Msg: respRpc.Msg,
	}

	c.JSON(200, resp)
}
