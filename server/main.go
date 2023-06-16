package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/generic"
	r "github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	consulApi "github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
)

const (
	serviceName         = "echo.server"
	timeToLift          = "5s"
	deregisterTimeout   = "12s"
	healthCheckInterval = "7s"
	consulAddr          = "127.0.0.1:8500"
)

func main() {
	// Parse IDL with Local Files
	// YOUR_IDL_PATH thrift file path,eg: ./idl/example.thrift
	p, err := generic.NewThriftFileProvider("../example_service.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	registry, err := consul.NewConsulRegister(consulAddr, consul.WithCheck(&consulApi.AgentServiceCheck{
		Interval:                       healthCheckInterval,
		Timeout:                        timeToLift,
		DeregisterCriticalServiceAfter: deregisterTimeout,
		TLSSkipVerify:                  true,
	}))

	if err != nil {
		log.Fatal(err)
	}

	//Fetch an available address
	serverAddr := getAddr()
	serverHostPortOpt := server.WithServiceAddr(serverAddr)

	//Weight param is later used for Load Balancing. Weight param must be > 0.
	svr := genericserver.NewServer(new(GenericServiceImpl), g, serverHostPortOpt, server.WithRegistry(registry), server.WithRegistryInfo(&r.Info{
		ServiceName: serviceName,
		Weight:      1,
	}))

	svr.Run()
	// resp is a JSON string
}

type GenericServiceImpl struct {
}

func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	// use jsoniter or other json parse sdk to assert request
	m := request.(string)
	fmt.Printf("Recv: %v\n", m)
	return "{\"Msg\": \"World!!\"}", nil
}

// Recursively call TCP resolver until an address is found. Possible infinite loop.
// Note: port 0 is a wilcard port which tells system to fetch a random available port.
func getAddr() *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Error occured." + err.Error() + "Retrying")
		return getAddr()
	}
	return addr
}
