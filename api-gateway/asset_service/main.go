package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	asset_management "api-gateway/asset_service/kitex_gen/asset_management/assetmanagement"
	"os"

	server "github.com/cloudwego/kitex/server"
)

func main() {

	port := "127.0.0.1:" + getPort()
	addr, _ := net.ResolveTCPAddr("tcp", port)
	svr := asset_management.NewServer(new(AssetManagementImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

// returns the port retrieved from cmd, if input is "", then from port.config file. If that fails, returns "8080" /default port.
func getPort() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please enter port: ")

	scanner.Scan()
	input := scanner.Text()

	//read the .config file if no port is mentioned
	if len(input) == 0 {
		fmt.Println("Reading port.config for port")
		data, err := ioutil.ReadFile("port.config")
		if err != nil {
			fmt.Println("File reading error. Using port 8080", err)
			return "8080" //default port
		}
		return string(data)
	}

	return string(input)
}
