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

// Creating the service registry which other files can use
func NewServiceRegistry() *ServiceRegistry {
	fmt.Println("Service registry instantiated!")
	return &ServiceRegistry{}
}

// Absolute path. Might not work for multiple gateway instances
// However, the service registry will be using a DB anyways, so no need to work on the path further.
var pathOfFile = "service_registry/services.csv"

func (sr *ServiceRegistry) GetHostPort(serviceID string) (string, error) {
	// getting the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Getting the relative file path of services.csv
	filePath := filepath.Join(workingDir, pathOfFile)

	// Opening up the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// getting ready to read
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	// Finding the host port for a particular unique service id
	for _, record := range records {
		if record[0] == serviceID {
			return record[2], nil
		}
	}

	return "", errors.New("service not found")
}
