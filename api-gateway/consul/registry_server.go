package consul

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/registry/consul"
)

var (
	wg      sync.WaitGroup
	localIP = "127.0.0.1"
)

func main() {
	config := consul.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	consulClient, err := consul.NewConsulClient(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		addr := net.JoinHostPort(localIP, "8888")
		r := consul.NewConsulRegister(consulClient)
		h := server.Default(
			server.WithHostPorts(addr),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "hertz.test.demo",
				Addr:        utils.NewNetAddr("tcp", addr),
				Weight:      10,
				Tags:        nil,
			}),
		)

		h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(200, utils.H{"ping": "pong1"})
		})
		h.Spin()
	}()

	wg.Wait()
}
