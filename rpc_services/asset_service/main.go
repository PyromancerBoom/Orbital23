package main

import (
	"log"
	"net"
	"os"
	asset_management "rpc_services/asset_service/kitex_gen/asset_management/assetmanagement"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	server "github.com/cloudwego/kitex/server"
)

var addr = getAddr()

func init() {
	config, err := LoadConfiguration("serviceConfig.json")
	if err != nil {
		log.Fatal(err)
	}

	gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName)
	advertisedPort := os.Getenv("PORT")

	id, err := gatewayClient.connectServer(config.ServiceURL, advertisedPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	go gatewayClient.updateHealthLoop(id, 5)
}

func main() {
	config, err := LoadConfiguration("serviceConfig.json")
	if err != nil {
		log.Fatal(err)
	}

	url := config.URL
	port := config.Port
	addrDocker, _ := net.ResolveTCPAddr("tcp", url+":"+port)

	svr := asset_management.NewServer(new(AssetManagementImpl),
		server.WithServiceAddr(addrDocker),
		server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
		server.WithReadWriteTimeout(100*time.Second))

	kitexerr := svr.Run()

	if kitexerr != nil {
		log.Println(kitexerr.Error())
	}
}
