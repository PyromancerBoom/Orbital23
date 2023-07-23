package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	asset_management "rpc_services/asset_service/kitex_gen/asset_management/assetmanagement"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	server "github.com/cloudwego/kitex/server"
)

type Configuration struct {
	URL        string `json:"url"`
	Port       string `json:"port"`
	Env        string `json:"env"`
	ServiceURL string `json:"serviceurl"`
}

const (
	apikey = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	// gateway = "http://0.0.0.0:4200"
	gateway = "http://host.docker.internal:4200"
)

var addr = getAddr()

func init() {

	config, err := LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}

	//make a client
	gatewayClient := NewGatewayClient(apikey, "AssetManagement", gateway)

	//register the server to the system
	// id, err := gatewayClient.connectServer(addr.IP.String(), strconv.Itoa(addr.Port))
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	advertisedPort := os.Getenv("PORT")

	id, err := gatewayClient.connectServer(config.ServiceURL, advertisedPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	//enter a health loop
	gatewayClient.updateHealthLoop(id, 5)
}

func main() {
	// Make a client
	// gatewayClient := NewGatewayClient(apikey, "AssetManagement", gateway)

	// // Connect to the gateway server with retry
	// id, err := connectServerWithRetry(gatewayClient, addr.IP.String(), strconv.Itoa(addr.Port))
	// if err != nil {
	// 	log.Println("Error connecting to gateway:", err.Error())
	// 	return
	// }

	// // Enter a health loop
	// gatewayClient.updateHealthLoop(id, 5)
	// ip := os.Getenv("IP")

	config, err := LoadConfiguration("serverConfig.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.URL)
	fmt.Println(config.Port)
	fmt.Println(config.Env)

	url := config.URL
	port := config.Port
	// port := os.Getenv("PORT")
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

func LoadConfiguration(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}
	var c Configuration
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}
	return c, nil
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
