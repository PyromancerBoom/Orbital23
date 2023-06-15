package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

func main() {
	// Parse IDL with Local Files
	// YOUR_IDL_PATH thrift file path, eg:./idl/example.thrift
	p, err := generic.NewThriftFileProvider("../example_service.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	connTimeout := client.WithRPCTimeout(3 * time.Second)
	cli, err := genericclient.NewClient("psm", g, client.WithHostPorts("0.0.0.0:8888"), connTimeout)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 'ExampleMethod' method name must be passed as param
	resp, err := cli.GenericCall(ctx, "ExampleMethod", "{\"Msg\": \"hello\", \"Base\":{\"Addr\":\"0.0.0.0:1224\"}}")
	if err != nil {
		panic(err)
	}
	// resp is a JSON string
	fmt.Print(resp)
}
