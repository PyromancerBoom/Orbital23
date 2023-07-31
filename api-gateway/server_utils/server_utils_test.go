package server_utils

import (
	"testing"
)

func TestNewGatewayClient(t *testing.T) {
	// Test case 1: API Key, Service Name and Gateway Address are blank
	t.Run("Blank API Key, Service Name and Gateway Address", func(t *testing.T) {
		client := NewGatewayClient("", "", "")
		if client != nil {
			t.Error("Expected NewGatewayClient to return nil for blank API Key, Service Name and Gateway Address")
		}
	})

	// Test case 2: API Key is blank
	t.Run("Blank API Key", func(t *testing.T) {
		client := NewGatewayClient("", "serviceName", "gatewayAddress")
		if client != nil {
			t.Error("Expected NewGatewayClient to return nil for blank API Key")
		}
	})

	// Test case 3: Service Name is blank
	t.Run("Blank Service Name", func(t *testing.T) {
		client := NewGatewayClient("apiKey", "", "gatewayAddress")
		if client != nil {
			t.Error("Expected NewGatewayClient to return nil for blank Service Name")
		}
	})

	// Test case 4: Gateway Address is blank
	t.Run("Blank Gateway Address", func(t *testing.T) {
		client := NewGatewayClient("apiKey", "serviceName", "")
		if client != nil {
			t.Error("Expected NewGatewayClient to return nil for blank Gateway Address")
		}
	})

	// Test case 5: API Key, Service Name and Gateway Address are not blank
	t.Run("Valid API Key, Service Name and Gateway Address", func(t *testing.T) {
		apiKey := "validApiKey"
		serviceName := "validServiceName"
		gatewayAddress := "validGatewayAddress"
		client := NewGatewayClient(apiKey, serviceName, gatewayAddress)
		if client == nil {
			t.Error("Expected NewGatewayClient to return a valid GatewayClient instance")
		}
		if client.ApiKey != apiKey {
			t.Errorf("Expected ApiKey to be %s, but got %s", apiKey, client.ApiKey)
		}
		if client.ServiceName != serviceName {
			t.Errorf("Expected ServiceName to be %s, but got %s", serviceName, client.ServiceName)
		}
		if client.GatewayAddress != gatewayAddress {
			t.Errorf("Expected GatewayAddress to be %s, but got %s", gatewayAddress, client.GatewayAddress)
		}
	})
}
