package main

import (
	"log"
	"net"
	"os"
	asset_management "rpc_services/asset_service/kitex_gen/asset_management/assetmanagement"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/server"
)

/*
How to host in local? simply set the isDockerised in .config to false!
*/

func main() {
	// Load configurations for Server
	config, err := LoadConfiguration("serviceConfig.json")
	if err != nil {
		log.Fatal(err)
	}

	if config.IsDockerised == false {

		// --------------------- Server is hosted in local ---------------------

		//if the port is set to 0, get a random port.
		var addr *net.TCPAddr
		if config.ServerPort == "0" {
			addr = getAddr() // Function in utils.go
		} else {
			addr, err = MakeAddress("127.0.0.1", config.ServerPort) // Function in utils.go
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName, config.GatewayAddress)

		id, err := gatewayClient.ConnectServer(addr.IP.String(), strconv.Itoa(addr.Port))
		if err != nil {
			log.Fatal(err.Error())
		}

		go gatewayClient.UpdateHealthLoop(id, config.HealthCheckFrequency)

		svr := asset_management.NewServer(new(AssetManagementImpl),
			server.WithServiceAddr(addr),
			server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
			server.WithReadWriteTimeout(100*time.Second))

		kitexerr := svr.Run()

		if kitexerr != nil {
			log.Println(kitexerr.Error())
		}

	} else {

		// --------------------- Server is hosted in docker ---------------------

		// This is supposed to be the default, however, please set it in the config file.
		// gatewayAddress = "http://host.docker.internal:4200"

		gatewayAddress := config.GatewayAddress
		gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName, gatewayAddress)

		advertisedPort := os.Getenv("PORT")

		id, err := gatewayClient.ConnectServer(config.ServiceURL, advertisedPort)
		if err != nil {
			log.Fatal(err.Error())
		}

		go gatewayClient.UpdateHealthLoop(id, config.HealthCheckFrequency)

		url := config.DockerUrl
		port := config.DockerPort
		addrDocker, err := net.ResolveTCPAddr("tcp", url+":"+port)
		if err != nil {
			log.Fatal(err.Error())
		}

		svr := asset_management.NewServer(new(AssetManagementImpl),
			server.WithServiceAddr(addrDocker),
			server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
			server.WithReadWriteTimeout(100*time.Second))

		kitexerr := svr.Run()

		if kitexerr != nil {
			log.Println(kitexerr.Error())
		}

	}
}
