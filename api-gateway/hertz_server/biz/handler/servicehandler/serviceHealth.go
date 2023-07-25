package servicehandler

import (
	"context"
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type HealthRequest struct {
	ServerID string `json:"ServerID"`
	ApiKey   string `json:"ApiKey"`
}

var consulClient *consul.Client

func init() {
	var err error
	consulClient, err = consul.NewClient(&consul.Config{})
	if err != nil {
		zap.L().Error("Failed to create consul client", zap.Error(err))
	}
}

// HealthCheck is the handler for health check requests
// @Route = /health
func HealthCheck(ctx context.Context, c *app.RequestContext) {
	var req HealthRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusExpectationFailed, "Invalid Request.")
		return
	}

	// Not required. Health checks need only be done for authorised connected servers
	if !authoriseHealthCheckRight(req.ApiKey, req.ServerID) {
		c.String(consts.StatusUnauthorized, "Unauthorized.")
		return
	}

	err2 := updateAsHealthy(req.ServerID)
	if err2 != nil {
		c.String(consts.StatusInternalServerError, "Unable to process health update request.")
		return
	}

	res := make(map[string]string)
	res["status"] = "status OK"
	res["message"] = "Successfully Updated the healh of server"

	c.JSON(consts.StatusOK, res)
}

// Updates a service as Healthy
// @Params
// checkID: The ID of the service to be updated
// @Returns
// error: If any error occurs while updating the service
func updateAsHealthy(checkID string) error {
	maxRetry := 10
	for retry := 0; retry < maxRetry; retry++ {
		err := consulClient.Agent().UpdateTTL(checkID, "online", consul.HealthPassing)
		if err == nil {
			return nil // Health update successful
		}

		zap.L().Error("Failed to update health. Retrying health update in 5 seconds...", zap.Error(err))
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("failed to update health after %d retries", maxRetry)
}

// Mocked
func authoriseHealthCheckRight(apiKey string, serverID string) bool {
	return true
}
