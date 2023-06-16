package main

import (
	userservice "api-gateway/userService/kitex_gen/UserService/userservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	port := "127.0.0.1:8888"
	addr, _ := net.ResolveTCPAddr("tcp", port)
	svr := userservice.NewServer(new(UserServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
