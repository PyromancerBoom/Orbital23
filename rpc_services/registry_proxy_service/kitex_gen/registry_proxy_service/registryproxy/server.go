// Code generated by Kitex v0.5.2. DO NOT EDIT.
package registryproxy

import (
	server "github.com/cloudwego/kitex/server"
	registry_proxy_service "registry_proxy/kitex_gen/registry_proxy_service"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler registry_proxy_service.RegistryProxy, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
