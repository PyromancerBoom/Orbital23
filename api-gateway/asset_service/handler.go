package main

import (
	asset_management "api-gateway/asset_service/kitex_gen/asset_management"
	"context"
	"fmt"
)

// AssetManagementImpl implements the last service interface defined in the IDL.
type AssetManagementImpl struct{}

type AssetInfo struct {
	ID     string
	Name   string
	Market string
}

var AssetData = make(map[string]AssetInfo, 5)

// QueryAsset implements the AssetManagementImpl interface.
func (s *AssetManagementImpl) QueryAsset(ctx context.Context, req *asset_management.QueryAssetRequest) (resp *asset_management.QueryAssetResponse, err error) {
	// TODO: Your code here...
	// TODO: Your code here...

	// fmt.Println("\nReached START Query Asset")

	// ast, exist := AssetData[req.ID]
	// if !exist {
	// 	return &asset_management.QueryAssetResponse{
	// 		Exist: false,
	// 	}, nil
	// }

	// fmt.Println("Asset data:")
	// fmt.Println(AssetData)

	// resp = &asset_management.QueryAssetResponse{
	// 	Exist:  true,
	// 	ID:     ast.ID,
	// 	Name:   ast.Name,
	// 	Market: ast.Market,
	// }
	resp = &asset_management.QueryAssetResponse{
		Exist:  true,
		ID:     "2",
		Name:   "CompanyName",
		Market: "MarketHere",
	}

	fmt.Println("\nReached end Query Asset")

	return resp, nil
}

// InsertAsset implements the AssetManagementImpl interface.
func (s *AssetManagementImpl) InsertAsset(ctx context.Context, req *asset_management.InsertAssetRequest) (resp *asset_management.InsertAssetResponse, err error) {
	// TODO: Your code here...
	// TODO: Your code here...
	fmt.Println("\nReached Start InsertAsset")
	_, exist := AssetData[req.ID]
	if exist {
		return &asset_management.InsertAssetResponse{
			Ok:  false,
			Msg: "the id already exists",
		}, nil
	}

	AssetData[req.ID] = AssetInfo{
		ID:     req.ID,
		Name:   req.Name,
		Market: req.Market,
	}

	fmt.Println("Asset data:")
	fmt.Println(AssetData)

	fmt.Println("\nReached END InsertAsset")

	return &asset_management.InsertAssetResponse{
		Ok: true,
	}, nil
}
