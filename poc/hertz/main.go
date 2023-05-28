package main

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	h := server.Default()

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	// Test -> A method that should return a string passed as param in body
	h.POST("/echo", func(c context.Context) error {
		var input map[string]interface{}
		if err := c.BindJSON(&input); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, input)
	})

	h.Spin()
}
