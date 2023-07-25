// Code generated by Kitex v0.5.2. DO NOT EDIT.

package registryproxy

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	registry_proxy_service "registry_proxy/kitex_gen/registry_proxy_service"
)

func serviceInfo() *kitex.ServiceInfo {
	return registryProxyServiceInfo
}

var registryProxyServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "RegistryProxy"
	handlerType := (*registry_proxy_service.RegistryProxy)(nil)
	methods := map[string]kitex.MethodInfo{
		"connectServer":     kitex.NewMethodInfo(connectServerHandler, newRegistryProxyConnectServerArgs, newRegistryProxyConnectServerResult, false),
		"healthCheckServer": kitex.NewMethodInfo(healthCheckServerHandler, newRegistryProxyHealthCheckServerArgs, newRegistryProxyHealthCheckServerResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "registry_proxy_service",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.5.2",
		Extra:           extra,
	}
	return svcInfo
}

func connectServerHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*registry_proxy_service.RegistryProxyConnectServerArgs)
	realResult := result.(*registry_proxy_service.RegistryProxyConnectServerResult)
	success, err := handler.(registry_proxy_service.RegistryProxy).ConnectServer(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newRegistryProxyConnectServerArgs() interface{} {
	return registry_proxy_service.NewRegistryProxyConnectServerArgs()
}

func newRegistryProxyConnectServerResult() interface{} {
	return registry_proxy_service.NewRegistryProxyConnectServerResult()
}

func healthCheckServerHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*registry_proxy_service.RegistryProxyHealthCheckServerArgs)
	realResult := result.(*registry_proxy_service.RegistryProxyHealthCheckServerResult)
	success, err := handler.(registry_proxy_service.RegistryProxy).HealthCheckServer(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newRegistryProxyHealthCheckServerArgs() interface{} {
	return registry_proxy_service.NewRegistryProxyHealthCheckServerArgs()
}

func newRegistryProxyHealthCheckServerResult() interface{} {
	return registry_proxy_service.NewRegistryProxyHealthCheckServerResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) ConnectServer(ctx context.Context, req *registry_proxy_service.ConnectRequest) (r *registry_proxy_service.ConnectResponse, err error) {
	var _args registry_proxy_service.RegistryProxyConnectServerArgs
	_args.Req = req
	var _result registry_proxy_service.RegistryProxyConnectServerResult
	if err = p.c.Call(ctx, "connectServer", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) HealthCheckServer(ctx context.Context, req *registry_proxy_service.HealtRequest) (r *registry_proxy_service.HealthResponse, err error) {
	var _args registry_proxy_service.RegistryProxyHealthCheckServerArgs
	_args.Req = req
	var _result registry_proxy_service.RegistryProxyHealthCheckServerResult
	if err = p.c.Call(ctx, "healthCheckServer", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
