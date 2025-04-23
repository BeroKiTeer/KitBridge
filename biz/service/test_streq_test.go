package service

import (
	"context"
	stability "github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability"
	"testing"
)

func TestTestSTReq_Run(t *testing.T) {
	ctx := context.Background()
	s := NewTestSTReqService(ctx)
	// init req and assert value

	req := &stability.STRequest{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
