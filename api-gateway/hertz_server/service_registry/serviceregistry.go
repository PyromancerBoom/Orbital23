package serviceregistry

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ServiceRegistry struct{}

type Service struct {
	ServiceID   string
	ServiceName string
	HostPort    string
}

// Creating the service registry
func NewServiceRegistry() *ServiceRegistry {
	fmt.Println("Service registry instantiated!")
	return &ServiceRegistry{}
}

var pathOfFile = "service_registry/services.csv"

func (sr *ServiceRegistry) GetHostPort(serviceID string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(wd, pathOfFile)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	for _, record := range records {
		if record[0] == serviceID {
			return record[2], nil
		}
	}

	return "", errors.New("service not found")
}
