package main

import (
	"log"
	asset_management "api-gateway/kitex_server/kitex_gen/asset_management/assetmanagement"
)

func main() {
	svr := asset_management.NewServer(new(AssetManagementImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
