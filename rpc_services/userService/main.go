package main

import (
	"log"
	"net"
	"os"
	userservice "rpc_services/userService/kitex_gen/UserService/userservice"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/server"
)

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

		id, err := gatewayClient.connectServer(addr.IP.String(), strconv.Itoa(addr.Port))
		if err != nil {
			log.Fatal(err.Error())
		}

		go gatewayClient.updateHealthLoop(id, config.HealthCheckFrequency)

		svr := userservice.NewServer(new(UserServiceImpl),
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

		//still a bit confused about this. recheck with docker.
		//advertisedPort = GetFreePort()

		id, err := gatewayClient.connectServer(config.ServiceURL, advertisedPort)
		if err != nil {
			log.Fatal(err.Error())
		}

		go gatewayClient.updateHealthLoop(id, config.HealthCheckFrequency)

		url := config.DockerUrl
		port := config.DockerPort
		addrDocker, err := net.ResolveTCPAddr("tcp", url+":"+port)
		if err != nil {
			log.Fatal(err.Error())
		}

		svr := userservice.NewServer(new(UserServiceImpl),
			server.WithServiceAddr(addrDocker),
			server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
			server.WithReadWriteTimeout(100*time.Second))

		kitexerr := svr.Run()

		if kitexerr != nil {
			log.Println(kitexerr.Error())
		}
	}

}
