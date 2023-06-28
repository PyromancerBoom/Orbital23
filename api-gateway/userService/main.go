package main

import (
	userservice "api-gateway/userService/kitex_gen/UserService/userservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func init() {
	//MAKE curl request. get
}

func main() {
	//0 port is a wildercard port. Use it after registration.
	port := "127.0.0.1:8888"
	addr, _ := net.ResolveTCPAddr("tcp", port)

	//add code here to make register itself to the ecosystem. Bascially send request

	svr := userservice.NewServer(new(UserServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
