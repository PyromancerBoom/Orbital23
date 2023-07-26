package server_utils

import (
	"testing"
)

func TestNewGatewayClient(t *testing.T) {
	t.Run("with valid input", func(t *testing.T) {
		apiKey := "master_api_key_uuid"
		serviceName := "test-service-name"
		client := NewGatewayClient(apiKey, serviceName)

		if client.ApiKey != apiKey {
			t.Errorf("NewGatewayClient() API key = %s, expected %s", client.ApiKey, apiKey)
		}

		if client.ServiceName != serviceName {
			t.Errorf("NewGatewayClient() service name = %s, expected %s", client.ServiceName, serviceName)
		}
	})

	t.Run("with invalid input", func(t *testing.T) {
		apiKey := "another-api-key"
		serviceName := "another-service"
		client := NewGatewayClient(apiKey, serviceName)

		if client.ApiKey == apiKey {
			t.Errorf("NewGatewayClient() API key = %s, expected %s", client.ApiKey, "different-api-key")
		}

		if client.ServiceName == serviceName {
			t.Errorf("NewGatewayClient() service name = %s, expected %s", client.ServiceName, "different-service-name")
		}
	})
}
