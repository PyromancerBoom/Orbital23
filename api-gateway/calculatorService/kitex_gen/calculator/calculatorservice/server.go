// Code generated by Kitex v0.6.0. DO NOT EDIT.
package calculatorservice

import (
	calculator "api-gateway/calculatorService/kitex_gen/calculator"
	server "github.com/cloudwego/kitex/server"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler calculator.CalculatorService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
