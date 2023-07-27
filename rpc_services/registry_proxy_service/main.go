package main

import (
	"log"
	"net"
	registry_proxy_service "registry_proxy/kitex_gen/registry_proxy_service/registryproxy"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
)

func main() {

	settings, err := LoadSettings("serviceConfig.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	r, err := consul.NewConsulRegister(settings.ConsulAddress)
	if err != nil {
		log.Fatal(err)
	}

	//if the port is set to 0, get a random port.
	var addr *net.TCPAddr
	if settings.ServerPort == "0" {
		addr = getAddr() // Function in utils.go
	} else {
		addr, err = MakeAddress("127.0.0.1", settings.ServerPort) // Function in utils.go
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	svr := registry_proxy_service.NewServer(
		new(RegistryProxyImpl),
		server.WithRegistry(r),
		server.WithRegistryInfo(&registry.Info{
			ServiceName: "RegistryProxy",
			Weight:      1,
			Addr:        addr}),
		server.WithServiceAddr(addr),
		server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
		server.WithReadWriteTimeout(100*time.Second),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
