package main

import (
	asset_management "api-gateway/asset_service/kitex_gen/asset_management/assetmanagement"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	server "github.com/cloudwego/kitex/server"
)

const (
	apikey  = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	gateway = "http://127.0.0.1:4200"
)

var addr = getAddr()

func init() {
	id, err := connectServer(gateway, apikey, "AssetManagement", addr.IP.String(), strconv.Itoa(addr.Port))
	if err != nil {
		log.Fatal(err.Error())
	}

	go healthCheckLoop(gateway, apikey, id)
}

func healthCheckLoop(gateway string, api string, id string) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		updateHealth(gateway, api, id)
		<-ticker.C
	}
}

func main() {

	svr := asset_management.NewServer(new(AssetManagementImpl), server.WithServiceAddr(addr))
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
