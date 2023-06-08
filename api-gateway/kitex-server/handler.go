package main

import (
	asset_management "api-gateway/kitex-server/kitex_gen/asset_management"
	"context"
)

// AssetManagementImpl implements the last service interface defined in the IDL.
type AssetManagementImpl struct{}

// QueryAsset implements the AssetManagementImpl interface.
func (s *AssetManagementImpl) QueryAsset(ctx context.Context, req *asset_management.QueryAssetRequest) (resp *asset_management.QueryAssetResponse, err error) {
	// TODO: Your code here...
	return
}

// InsertAsset implements the AssetManagementImpl interface.
func (s *AssetManagementImpl) InsertAsset(ctx context.Context, req *asset_management.InsertAssetRequest) (resp *asset_management.InsertAssetResponse, err error) {
	// TODO: Your code here...
	return
}
