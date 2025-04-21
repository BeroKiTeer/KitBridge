package service

import (
	"context"
	api "github.com/BeroKiTeer/KitBridge/kitex_gen/api"
)

type EchoService struct {
	ctx context.Context
} // NewEchoService new EchoService
func NewEchoService(ctx context.Context) *EchoService {
	return &EchoService{ctx: ctx}
}

// Run create note info
func (s *EchoService) Run(req *api.Request) (resp *api.Response, err error) {
	// Finish your business logic.
	resp = &api.Response{
		Message: req.Message,
	}
	return
}
