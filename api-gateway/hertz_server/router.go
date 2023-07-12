// Code generated by hertz generator.

package main

import (
	handler "api-gateway/hertz_server/biz/handler"

	servicehandler "api-gateway/hertz_server/biz/handler/servicehandler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {
	r.GET("/ping", handler.Ping)

	// your code ...
	// r.POST("/register", servicehandler.Register)
	r.PUT("/update", servicehandler.Update)
	r.POST("/connect", servicehandler.Connect)
	r.POST("/health", servicehandler.HealthCheck)

	// ---------------------------------------
	// Remove this endpoint before production
	// r.GET("/show", servicehandler.DisplayAll)

	r.GET("/testmongo", handler.InsertDataInMongo)
}
