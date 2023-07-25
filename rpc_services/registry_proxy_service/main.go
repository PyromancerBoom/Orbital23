package main

import (
	"fmt"
	"log"
	"net"
	registry_proxy_service "registry_proxy/kitex_gen/registry_proxy_service/registryproxy"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
)

func main() {

	r, err := consul.NewConsulRegister("127.0.0.1:8500")
	if err != nil {
		log.Fatal(err)
	}

	addr := getAddr()

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

func getAddr() *net.TCPAddr {

	port, _ := GetFreePort()

	a := "127.0.0.1:" + strconv.Itoa(port)

	addr, err := net.ResolveTCPAddr("tcp", a)
	if err != nil {
		fmt.Println("Error occured." + err.Error() + "Retrying")
		return getAddr()
	}
	return addr
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
