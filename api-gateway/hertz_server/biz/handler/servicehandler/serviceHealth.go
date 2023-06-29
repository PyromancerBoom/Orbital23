package servicehandler

import (
	"context"
	"log"

	consul "github.com/hashicorp/consul/api"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type HealthRequest struct {
	ServerID string `json:"serverID"`
	APIKey   string `json:"api-key"`
}

var consulClient *consul.Client

func init() {
	var err error
	consulClient, err = consul.NewClient(&consul.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func HealthCheck(ctx context.Context, c *app.RequestContext) {
	var req HealthRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusExpectationFailed, "Invalid Request.")
		return
	}

	if !authoriseHealthCheckRight(req.APIKey, req.ServerID) {
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

// TODO:
// Implement auth
func authoriseHealthCheckRight(apiKey string, serverID string) bool {
	return true
}

func updateAsHealthy(checkID string) error {
	err2 := consulClient.Agent().UpdateTTL(checkID, "online", consul.HealthPassing)
	if err2 != nil {
		println(err2.Error())
		return err2
	}

	return nil
}