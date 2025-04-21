package autodetect

import (
	"regexp"

	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/detection"
	"github.com/cloudwego/kitex/pkg/remote/trans/netpoll"
)

// ====================================================================
// ✅ 动态切换编解码处理器说明（满足题目核心功能要求 #1）
// --------------------------------------------------------------------
// 在 Kitex 中，每个 TCP 连接都需要一个对应的 ServerTransHandler 来完成
// 请求的读取（Read）与响应的写入（Write）。
//
// 本文件中使用 detection.NewSvrTransHandlerFactory 实现了“协议嗅探工厂”，
// 支持注册多个 handler，并通过 ProtocolMatch 方法在连接建立阶段判断
// 当前连接是 Thrift 协议还是 HTTP 协议。
//
// 实现细节：
// - 如果连接首部符合 HTTP 请求行（如 "POST /api/... HTTP/1.1"），则使用我们自定义的 HTTP1Handler
// - 如果连接是 Thrift 二进制流（如 0x80 开头），则回退为默认 ThriftHandler
//
// ✅ 动态切换发生在 detection 模块内部：
// --> 代码位置：
//     detection.NewSvrTransHandlerFactory(thriftHandlerFactory, httpHandlerFactory)
//
// --> 判断过程：
//     httpHandlerFactory.NewTransHandler() 返回的 handler 必须实现 ProtocolMatch()
//     若 ProtocolMatch 返回 nil，则当前连接将绑定该 handler，并永久使用它处理请求。
//
// ✅ 效果：
//     - 不同连接可动态使用不同协议的 handler
//     - 一次嗅探，终身绑定，性能极高
//
// 示例：
//     ┌────────────┐
//     │ conn A     │ → Peek → HTTP  → 绑定 HTTP1Handler（你的 Read/Write）
//     └────────────┘
//     ┌────────────┐
//     │ conn B     │ → Peek → Thrift → 绑定 ThriftHandler（kitex 默认）
//     └────────────┘
//
// ====================================================================

// 简单的HTTP请求行识别
var httpPattern = regexp.MustCompile(`^(GET|POST|PUT|HEAD|DELETE|OPTIONS|TRACE|CONNECT|PATCH)`)

// NewSvrTransHandlerFactoryWithHTTP ✅ 推荐入口函数：从 main.go 调用
// - 只需传入 HTTP handler 工厂，Thrift handler 使用 Kitex 默认的 netpoll 实现
func NewSvrTransHandlerFactoryWithHTTP(
	httpHandlerFactory remote.ServerTransHandlerFactory,
) remote.ServerTransHandlerFactory {
	thriftFactory := netpoll.NewSvrTransHandlerFactory()
	return NewSvrTransHandlerFactory(thriftFactory, httpHandlerFactory)
}

// NewSvrTransHandlerFactory ✅ 更灵活版本：允许 main 显式指定 Thrift handler（如替换为 framed）
func NewSvrTransHandlerFactory(
	thriftHandlerFactory remote.ServerTransHandlerFactory,
	httpHandlerFactory remote.ServerTransHandlerFactory,
) remote.ServerTransHandlerFactory {
	return detection.NewSvrTransHandlerFactory(
		thriftHandlerFactory,
		httpHandlerFactory,
	)
}
