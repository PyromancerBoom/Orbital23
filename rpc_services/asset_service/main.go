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

const (
	apikey  = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	gateway = "http://0.0.0.0:4200"
	// gateway = "http://host.docker.internal:4200"
)

var addr = getAddr()

func init() {

	config, err := LoadConfiguration("serviceConfig.json")
	if err != nil {
		log.Fatal(err)
	}

	gatewayClient := NewGatewayClient(apikey, "AssetManagement", gateway)

	advertisedPort := os.Getenv("PORT")

	id, err := gatewayClient.connectServer(config.ServiceURL, advertisedPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	gatewayClient.updateHealthLoop(id, 5)
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
		server.WithReadWriteTimeout(150*time.Second))

	kitexerr := svr.Run()

	if kitexerr != nil {
		log.Println(kitexerr.Error())
	}
}
