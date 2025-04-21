package http1

import (
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/remote"
	"net"
	"regexp"
)

type ServerTransHandler interface {
	Read(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error)
	Write(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error)
	OnRead(ctx context.Context, conn net.Conn) error
	OnInactive(ctx context.Context, conn net.Conn)
	OnError(ctx context.Context, err error, conn net.Conn)
	OnMessage(ctx context.Context, args, result remote.Message) (context.Context, error)
	SetPipeline(pipeline *remote.TransPipeline)
	SetInvokeHandleFunc(endpoint.Endpoint)
	OnActive(ctx context.Context, conn net.Conn) (context.Context, error)
}

var httpPattern = regexp.MustCompile(`^(GET|POST|PUT|HEAD|DELETE|OPTIONS|TRACE|CONNECT|PATCH)`)

type HTTP1SvrTransHandlerFactory struct{}

func (f *HTTP1SvrTransHandlerFactory) NewTransHandler(opt *remote.ServerOption) (remote.ServerTransHandler, error) {
	return &HTTP1Handler{}, nil
}

type HTTP1Handler struct{}

func (h *HTTP1Handler) ProtocolMatch(ctx context.Context, conn net.Conn) error {
	buf := make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if httpPattern.Match(buf[:n]) {
		return nil
	}
	return errors.New("not http")
}

func (h *HTTP1Handler) Read(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error) {
	// 下一阶段我们在这里解析 HTTP 请求
	return ctx, nil
}

func (h *HTTP1Handler) Write(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error) {
	// 下一阶段我们在这里构造 JSON 响应写回
	return ctx, nil
}

func (h *HTTP1Handler) OnRead(ctx context.Context, conn net.Conn) error {
	return nil
}

func (h *HTTP1Handler) OnInactive(ctx context.Context, conn net.Conn) {}

func (h *HTTP1Handler) OnError(ctx context.Context, err error, conn net.Conn) {}

func (h *HTTP1Handler) OnMessage(ctx context.Context, args, result remote.Message) (context.Context, error) {
	return ctx, nil
}

func (h *HTTP1Handler) SetPipeline(pipeline *remote.TransPipeline) {}

func (h *HTTP1Handler) SetInvokeHandleFunc(endpoint endpoint.Endpoint) {}

func (h *HTTP1Handler) OnActive(ctx context.Context, conn net.Conn) (context.Context, error) {
	return ctx, nil
}
