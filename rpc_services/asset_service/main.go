package main

import (
	asset_management "rpc_services/asset_service/kitex_gen/asset_management/assetmanagement"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	server "github.com/cloudwego/kitex/server"
)

const (
	apikey  = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	gateway = "http://127.0.0.1:4200"
)

var addr = getAddr()

func init() {

	//make a client
	gatewayClient := NewGatewayClient(apikey, "AssetManagement", gateway)

	//register the server to the system
	id, err := gatewayClient.connectServer(addr.IP.String(), strconv.Itoa(addr.Port))
	if err != nil {
		log.Fatal(err.Error())
	}

	//enter a health loop
	gatewayClient.updateHealthLoop(id, 5)
}

func main() {

	svr := asset_management.NewServer(new(AssetManagementImpl), server.WithServiceAddr(addr),
		server.WithLimit(&limit.Option{MaxConnections: 10000, MaxQPS: 1000}),
		server.WithReadWriteTimeout(100*time.Second))
	err := svr.Run()

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