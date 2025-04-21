package autodetect

import (
	"regexp"

	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/detection"
	"github.com/cloudwego/kitex/pkg/remote/trans/netpoll"
)

// 简单的HTTP请求行识别
var httpPattern = regexp.MustCompile(`^(GET|POST|PUT|HEAD|DELETE|OPTIONS|TRACE|CONNECT|PATCH)`)

func NewSvrTransHandlerFactoryWithHTTP(httpHandlerFactory remote.ServerTransHandlerFactory) remote.ServerTransHandlerFactory {
	thriftFactory := netpoll.NewSvrTransHandlerFactory()
	return detection.NewSvrTransHandlerFactory(
		thriftFactory,
		httpHandlerFactory, // 👈 放在这里，框架内部通过 ProtocolMatch 自动识别
	)
}
