package main

import (
	userservice "api-gateway/userService/kitex_gen/UserService/userservice"
	"log"
)

func main() {
	svr := userservice.NewServer(new(UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
