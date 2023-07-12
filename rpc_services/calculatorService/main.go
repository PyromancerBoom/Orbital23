package main

import (
	calculator "rpc_services/calculatorService/kitex_gen/calculator/calculatorservice"
	"log"
	"strconv"
)

const (
	apikey  = "36e991d3-646d-414a-ac66-0c0e8a310ced"
	gateway = "http://127.0.0.1:4200"
)

func init() {
	//make a client
	gatewayClient := NewGatewayClient(apikey, "CalculatorService", gateway)

	//register the server to the system
	id, err := gatewayClient.connectServer(addr.IP.String(), strconv.Itoa(addr.Port))
	if err != nil {
		log.Fatal(err.Error())
	}

	//enter a health loop
	gatewayClient.updateHealthLoop(id, 5)
}

func main() {
	svr := calculator.NewServer(new(CalculatorServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
