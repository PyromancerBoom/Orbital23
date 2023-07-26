package main

import (
	"log"
	asset_management "rpc_services/asset_service/kitex_gen/asset_management/assetmanagement"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/server"
)

/*
How to host in local ? Simply uncomment the code in main func as per the comments. Do not delete.
*/

func main() {
	// Load configurations for Server
	config, err := LoadConfiguration("serviceConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	// --------------------- Uncomment below to host in local ---------------------
	var addr = getAddr() // Function in utils.go

	gatewayAddress := "http://0.0.0.0:4200"
	gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName, gatewayAddress)

	id, err := gatewayClient.connectServer(addr.IP.String(), strconv.Itoa(addr.Port))
	if err != nil {
		log.Fatal(err.Error())
	}

	print(id)
	go gatewayClient.updateHealthLoop(id, 5)

	svr := asset_management.NewServer(new(AssetManagementImpl),
		server.WithServiceAddr(addr),
		server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
		server.WithReadWriteTimeout(100*time.Second))

	kitexerr := svr.Run()

	if kitexerr != nil {
		log.Println(kitexerr.Error())
	}

	// --------------------- Uncomment below to host in docker ---------------------

	// gatewayAddress = "http://host.docker.internal:4200"
	// gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName, gatewayAddress)
	// // advertisedPort := os.Getenv("PORT")

	// advertisedPort = GetFreePort()

	// id, err := gatewayClient.connectServer(config.ServiceURL, advertisedPort)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// go gatewayClient.updateHealthLoop(id, 5)

	// url := config.URL
	// port := config.Port
	// addrDocker, _ := net.ResolveTCPAddr("tcp", url+":"+port)

	// svr := asset_management.NewServer(new(AssetManagementImpl),
	// 	server.WithServiceAddr(addrDocker),
	// 	server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
	// 	server.WithReadWriteTimeout(100*time.Second))

	// kitexerr := svr.Run()

	// if kitexerr != nil {
	// 	log.Println(kitexerr.Error())
	// }
}
