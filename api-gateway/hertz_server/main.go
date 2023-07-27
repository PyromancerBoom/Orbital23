package main

import (
	"api-gateway/hertz_server/biz/model/cache"
	repository "api-gateway/hertz_server/biz/model/repository"
	"api-gateway/hertz_server/biz/model/settings"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/zap"
)

func init() {
	err := settings.InitialiseSettings("serverconfig.json")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	initLogger()

	serverSettings := settings.GetSettings()
	// hostURL := "127.0.0.1:4200"
	hostURL := "0.0.0.0:" + serverSettings.ServerPort

	if err := repository.ConnectToMongoDB(); err != nil {
		panic(err)
	}
	defer repository.CloseMongoDB()

	// Perform health check
	// Commented out to test performance improvement
	go repository.MongoHealthCheck()

	go cache.UpdateCache()

	h := server.Default(server.WithHostPorts(hostURL))

	zap.L().Info("Starting server", zap.String("hostURL: ", hostURL))

	register(h)
	h.Spin()
}
