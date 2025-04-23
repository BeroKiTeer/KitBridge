package service

import (
	"context"
	stability "github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability"
)

type TestSTReqService struct {
	ctx context.Context
} // NewTestSTReqService new TestSTReqService
func NewTestSTReqService(ctx context.Context) *TestSTReqService {
	return &TestSTReqService{ctx: ctx}
}

// Run create note info
func (s *TestSTReqService) Run(req *stability.STRequest) (resp *stability.STResponse, err error) {
	// Finish your business logic.
	resp = &stability.STResponse{
		Str:       req.Str,
		Mp:        req.StringMap,
		Name:      req.Name,
		Framework: req.Framework,
	}
	return resp, nil
}
