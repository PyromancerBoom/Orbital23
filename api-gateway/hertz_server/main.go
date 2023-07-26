package main

import (
	"api-gateway/hertz_server/biz/model/cache"
	repository "api-gateway/hertz_server/biz/model/repository"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/zap"
)

func main() {
	initLogger()

	hostURL := "0.0.0.0:4200"

	if err := repository.ConnectToMongoDB(); err != nil {
		panic(err)
	}
	// Perform health check
	go repository.MongoHealthCheck()

	cache.UpdateAllCache()

	h := server.Default(server.WithHostPorts(hostURL))

	zap.L().Info("Starting server", zap.String("hostURL: ", hostURL))

	register(h)
	h.Spin()
}
