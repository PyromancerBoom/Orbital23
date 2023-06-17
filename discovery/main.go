package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func main() {

	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	discoverClient, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	handle := discoverClient.Catalog()

	query := &api.QueryOptions{}

	s, _, _ := handle.Service("echo.server", "", query)

	for i := 0; i < len(s); i++ {
		fmt.Println(s[i].ServiceAddress, ":", s[i].ServicePort)
	}

}
