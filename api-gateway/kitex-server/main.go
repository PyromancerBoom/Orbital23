package main

import (
	asset_management "kitex-server/kitex_gen/asset_management/assetmanagement"
	"log"
)

func main() {
	svr := asset_management.NewServer(new(AssetManagementImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
