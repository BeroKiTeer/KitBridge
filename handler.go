package main

import (
	"context"
	stability "github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability"
)

// STServiceImpl implements the last service interface defined in the IDL.
type STServiceImpl struct{}

// TestSTReq implements the STServiceImpl interface.
func (s *STServiceImpl) TestSTReq(ctx context.Context, req *stability.STRequest) (resp *stability.STResponse, err error) {
	// TODO: Your code here...
	return
}
