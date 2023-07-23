package main

import (
	"context"
	"fmt"
	asset_management "rpc_services/asset_service/kitex_gen/asset_management"
	"sync"
)

// AssetManagementImpl implements the last service interface defined in the IDL.
type AssetManagementImpl struct {
	mu sync.Mutex
}

type AssetInfo struct {
	ID     string
	Name   string
	Market string
}

var AssetData = make(map[string]AssetInfo)

// QueryAsset implements the AssetManagementImpl interface.
func (s *AssetManagementImpl) QueryAsset(ctx context.Context, req *asset_management.QueryAssetRequest) (*asset_management.QueryAssetResponse, error) {
	// TODO: Your code here...
	s.mu.Lock()
	defer s.mu.Unlock()

	ast, exist := AssetData[req.ID]
	if !exist {
		return &asset_management.QueryAssetResponse{
			Exist: false,
		}, nil
	}

	// You can use a Logger library here instead of fmt.Println
	fmt.Println("Asset data:")
	fmt.Println(AssetData)

	return &asset_management.QueryAssetResponse{
		Exist:  true,
		ID:     ast.ID,
		Name:   ast.Name,
		Market: ast.Market,
	}, nil
}

// InsertAsset implements the AssetManagementImpl interface.
func (s *AssetManagementImpl) InsertAsset(ctx context.Context, req *asset_management.InsertAssetRequest) (*asset_management.InsertAssetResponse, error) {
	// TODO: Your code here...
	s.mu.Lock()
	defer s.mu.Unlock()

	// Basic data validation
	if req.ID == "" || req.Name == "" || req.Market == "" {
		return nil, fmt.Errorf("missing required fields")
	}

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

	// You can use a Logger library here instead of fmt.Println
	fmt.Println("Asset data:")
	fmt.Println(AssetData)

	return &asset_management.InsertAssetResponse{
		Ok: true,
	}, nil
}
