package registry

import (
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	consul "github.com/hashicorp/consul/api"
)

type Service struct {
	consulClient *api.Client
}

func NewService() *Service {
	// Create a new Consul client.
	client, err := consul.NewClient(&consul.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &Service{
		consulClient: client,
	}
}

// Constants for time durations and service IDs.
const (
	timetolive       = 5 * time.Second
	timetoderegister = timetolive * 2
	serviceID        = "myservice"
	checkID          = "6969-121"
)

// Register the service
func (s *Service) registerService() {
	// Health check with the necessary parameters
	check := &consul.AgentServiceCheck{
		DeregisterCriticalServiceAfter: timetoderegister.String(),
		TLSSkipVerify:                  true,
		TTL:                            timetolive.String(),
		CheckID:                        checkID,
	}

	// Create a service registration with the service details.
	reg := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "RPC cluster",
		Tags:    []string{"myservice"},
		Address: "127.0.0.1",
		Port:    3000,
		Check:   check,
	}

	// Register the service with Consul.
	err := s.consulClient.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) healthCheckLoop() {
	ticker := time.NewTicker(time.Second * 3)
	for {
		err := s.consulClient.Agent().UpdateTTL(checkID, "online", consul.HealthPassing)
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
