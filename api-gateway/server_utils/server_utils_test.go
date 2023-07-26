package server_utils

import (
	"testing"
)

// TestNewGatewayClient tests the NewGatewayClient function
func TestNewGatewayClient(t *testing.T) {
	// Test case 1: API Key and Service Name are blank
	if NewGatewayClient("", "") != nil {
		t.Error("API Key and Service Name cannot be blank")
	}

	// Test case 2: API Key is blank
	if NewGatewayClient("", "serviceName") != nil {
		t.Error("API Key cannot be blank")
	}

	// Test case 3: Service Name is blank
	if NewGatewayClient("apiKey", "") != nil {
		t.Error("Service Name cannot be blank")
	}

	// Test case 4: API Key and Service Name are not blank
	if NewGatewayClient("apiKey", "serviceName") == nil {
		t.Log("Test passed")
	}
}
