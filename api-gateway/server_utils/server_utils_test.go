package tests

import (
	"testing"
)

func TestNewGatewayClient(t *testing.T) {
	// Test case 1
	apiKey := "master_api_key_uuid"
	serviceName := "test-service-name"
	client := NewGatewayClient(apiKey, serviceName)

	// Assert that the GatewayClient fields are set correctly
	if client.ApiKey != apiKey {
		t.Errorf("NewGatewayClient() API key = %s, expected %s", client.ApiKey, apiKey)
	}

	if client.ServiceName != serviceName {
		t.Errorf("NewGatewayClient() service name = %s, expected %s", client.ServiceName, serviceName)
	}

	// Assert that the GatewayAddress is set to the default value
	if client.GatewayAddress != gatewayAddress {
		t.Errorf("NewGatewayClient() GatewayAddress = %s, expected %s", client.GatewayAddress, gatewayAddress)
	}

	// Test case 2 with different values
	apiKey = "wrong-api-key"
	serviceName = "another-service"
	client = NewGatewayClient(apiKey, serviceName)

	// Assert that the GatewayClient fields are set correctly for the new values
	if client.ApiKey != apiKey {
		t.Errorf("NewGatewayClient() API key = %s, expected %s", client.ApiKey, apiKey)
	}

	if client.ServiceName != serviceName {
		t.Errorf("NewGatewayClient() service name = %s, expected %s", client.ServiceName, serviceName)
	}
}
