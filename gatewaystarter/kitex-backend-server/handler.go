package main

import (
	"context"
)

// ServiceImpl implements the last service interface defined in the IDL.
type ServiceImpl struct{}

// DoSomething implements the ServiceImpl interface.
func (s *ServiceImpl) DoSomething(ctx context.Context, input string) (resp string, err error) {
	// TODO: Your code here...
	resp = "Response from first function: " + input
	return resp, nil
}

// DoSomethingMore implements the ServiceImpl interface.
func (s *ServiceImpl) DoSomethingMore(ctx context.Context, input string) (resp string, err error) {
	// TODO: Your code here...
	resp = "Response from second function: " + input
	return resp, nil
}
