// Code generated by Kitex v0.6.0. DO NOT EDIT.

package userservice

import (
	userservice "rpc_services/userService/kitex_gen/UserService"
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return userServiceServiceInfo
}

var userServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "UserService"
	handlerType := (*userservice.UserService)(nil)
	methods := map[string]kitex.MethodInfo{
		"queryUser":  kitex.NewMethodInfo(queryUserHandler, newUserServiceQueryUserArgs, newUserServiceQueryUserResult, false),
		"insertUser": kitex.NewMethodInfo(insertUserHandler, newUserServiceInsertUserArgs, newUserServiceInsertUserResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "userservice",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.6.0",
		Extra:           extra,
	}
	return svcInfo
}

func queryUserHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*userservice.UserServiceQueryUserArgs)
	realResult := result.(*userservice.UserServiceQueryUserResult)
	success, err := handler.(userservice.UserService).QueryUser(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceQueryUserArgs() interface{} {
	return userservice.NewUserServiceQueryUserArgs()
}

func newUserServiceQueryUserResult() interface{} {
	return userservice.NewUserServiceQueryUserResult()
}

func insertUserHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*userservice.UserServiceInsertUserArgs)
	realResult := result.(*userservice.UserServiceInsertUserResult)
	success, err := handler.(userservice.UserService).InsertUser(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newUserServiceInsertUserArgs() interface{} {
	return userservice.NewUserServiceInsertUserArgs()
}

func newUserServiceInsertUserResult() interface{} {
	return userservice.NewUserServiceInsertUserResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) QueryUser(ctx context.Context, req *userservice.QueryUser) (r *userservice.QueryUserResponse, err error) {
	var _args userservice.UserServiceQueryUserArgs
	_args.Req = req
	var _result userservice.UserServiceQueryUserResult
	if err = p.c.Call(ctx, "queryUser", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) InsertUser(ctx context.Context, req *userservice.InsertUser) (r *userservice.InsertUserResponse, err error) {
	var _args userservice.UserServiceInsertUserArgs
	_args.Req = req
	var _result userservice.UserServiceInsertUserResult
	if err = p.c.Call(ctx, "insertUser", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
