package idlMap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Represents the relationship between services, method, and IDLs
type IdlMapping struct {
	Service string
	Method  string
	IDL     string
}

var IdlMap []IdlMapping

func init() {
	if err := loadMappings(); err != nil {
		fmt.Printf("Failed to init idl mappings: %v\n", err)
	}
}

func loadMappings() error {
	jsonPath := "idl_mapping.json"

	// Reading file
	idlData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(idlData, &IdlMap); err != nil {
		return err
	}

	return nil
}

func getIdlFile(service, method string) (string, error) {
	for _, value := range IdlMap {
		if value.Service == service && value.Method == method {
			return value.IDL, nil
		}
	}

	return "", fmt.Errorf("IDL not found")
}
