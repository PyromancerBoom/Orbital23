package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func main() {
	// Create a new Consul client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Create a new agent instance
	agent := client.Agent()

	// Define the service registration parameters
	registration := &api.AgentServiceRegistration{
		ID:      "my-service",
		Name:    "My Service",
		Port:    8080,
		Address: "localhost",
	}

	// Register the service with Consul
	err = agent.ServiceRegister(registration)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Consul server started")

	// Wait indefinitely
	select {}
}
