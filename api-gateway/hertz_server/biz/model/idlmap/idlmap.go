package idlmap

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v3"
)

// Represents the relationship between services, method, and IDLs
type IdlMapping struct {
	Service string `yaml:"Service"`
	Path    string `yaml:"Path"`
	Method  string `yaml:"Method"`
	IDL     string `yaml:"Idl"`
}

var IdlHashMap map[string]IdlMapping

func init() {
	// Load IDL mappings from the YAML file
	loadIdlMappings()
}

func loadIdlMappings() {
	// Read the YAML file
	yamlFile, err := ioutil.ReadFile("biz/model/idlmap/idlmapping.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// Unmarshal YAML into a slice of IdlMapping
	var idlMappings []IdlMapping
	err = yaml.Unmarshal(yamlFile, &idlMappings)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Create the IDL hashmap
	IdlHashMap = make(map[string]IdlMapping)
	for _, mapping := range idlMappings {
		IdlHashMap[mapping.Service+"_"+mapping.Path] = mapping
	}
}

type GatewayService struct {
	Service string
	Path    string
	Method  string
	IDL     string
}

func GetService(service, path string) (GatewayService, error) {
	mapping, ok := IdlHashMap[service+"_"+path]
	if !ok {
		return GatewayService{}, fmt.Errorf("404: Service not found")
	}

	return GatewayService{
		Service: mapping.Service,
		Path:    mapping.Path,
		Method:  mapping.Method,
		IDL:     mapping.IDL,
	}, nil
}
