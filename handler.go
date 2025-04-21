package main

import (
	"context"
	"github.com/BeroKiTeer/KitBridge/biz/service"
	"github.com/BeroKiTeer/KitBridge/kitex_gen/api"
	stability "github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability"
)

// STServiceImpl implements the last service interface defined in the IDL.
type STServiceImpl struct{}

// TestSTReq implements the STServiceImpl interface.
func (s *STServiceImpl) TestSTReq(ctx context.Context, req *stability.STRequest) (resp *stability.STResponse, err error) {
	// TODO: Your code here...
	return
}

type HelloImpl struct{}

// Echo implements the HelloImpl interface.
func (s *HelloImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	resp, err = service.NewEchoService(ctx).Run(req)

	return resp, err
}
