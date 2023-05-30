package main

import (
	service "gatewaystarter/kitex_gen/service/service"
	"log"
)

func main() {
	svr := service.NewServer(new(ServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
