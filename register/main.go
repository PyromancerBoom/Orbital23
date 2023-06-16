package main

import (
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	consul "github.com/hashicorp/consul/api"
)

type Service struct {
	conssulClient *api.Client
}

func NewService() *Service {
	client, err := consul.NewClient(&consul.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &Service{
		conssulClient: client,
	}
}

// time to lift
const ttl = 5 * time.Second

// time to deregister
const ttd = ttl * 2

const id = "myservice"
const checkID = "6969-121"

// occurs once
func (s *Service) registerService() {
	check := &consul.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttd.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        checkID,
	}

	reg := &api.AgentServiceRegistration{
		ID:      id,
		Name:    "RPC cluster",
		Tags:    []string{"myservice"},
		Address: "127.0.0.1",
		Port:    3000,
		Check:   check,
	}

	err := s.conssulClient.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) healthCheckLoop() {
	ticker := time.NewTicker(time.Second * 3)
	for {
		err := s.conssulClient.Agent().UpdateTTL(checkID, "online", consul.HealthPassing)
		if err != nil {
			log.Fatal(err)
		}
		<-ticker.C
	}
}

func main() {
	s := NewService()

	s.registerService()
	s.healthCheckLoop()
}
