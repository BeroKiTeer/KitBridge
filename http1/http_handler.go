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

// 解析 HTTP 请求并转为 Kitex RPC 调用
func (h *HTTP1Handler) Read(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error) {
	// TODO 1: 读取完整 HTTP 请求数据（包含 请求行、Header、Body）
	// - 从 conn 读取
	// - 可设置最大长度，避免内存攻击

	// TODO 2: 解析 HTTP 请求行（例如：POST /api/STService/testSTReq HTTP/1.1）
	// - 拆出 path，校验是否符合 REST 映射规则
	// - 例如 /api/{Service}/{Method}

	// TODO 3: 映射为 Kitex 调用目标
	// - path → serviceName + methodName
	// - 填入 msg.SetServiceName() 和 msg.SetMethod()

	// TODO 4: 从 Header 中提取参数（支持 query/header 映射）
	// - 可选：根据 IDL 的 tag 配置自动映射 Header → struct 字段

	// TODO 5: 读取并解析 JSON Body → Thrift Request Struct
	// - body -> json.Unmarshal() -> STRequest（通过 methodName 映射 struct）

	// TODO 6: 将请求 struct 设置为 msg.Args
	// - msg.SetArgs(&req)
	// - msg.SetMessageType(remote.Call)

	// ✅ 最终效果：Kitex 调用链接收到服务名、方法名、请求 struct，可继续往下执行

	return ctx, nil
}

// 将 Kitex RPC 返回结果封装为标准 HTTP JSON 响应
func (h *HTTP1Handler) Write(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error) {
	// TODO 1: 判断调用结果是正常返回还是异常
	// - msg.RPCInfo().Invocation().BizStatusError() != nil → 业务异常
	// - msg.RPCInfo().Stats().Error() != nil → 框架异常

	// TODO 2: 构造标准 JSON 响应结构：
	// {
	//   "code": 0,
	//   "message": "success",
	//   "data": { ... } // Thrift 返回 struct 的 JSON 表达
	// }

	// TODO 3: 根据错误类型构造不同 code/message：
	// - 正常调用：code = 0, message = "success"
	// - BizError：code = 自定义错误码, message = 错误提示
	// - 系统异常：code = 500, message = "internal error"

	// TODO 4: 通过 json.Marshal(...) 将响应 struct 编码为 []byte

	// TODO 5: 构造完整 HTTP Response：
	// - 响应头：HTTP/1.1 200 OK + Content-Type: application/json + Content-Length
	// - 响应体：json 数据

	// TODO 6: 写入 conn（用 conn.Write(...) 输出响应）

	// ✅ 最终效果：HTTP 客户端收到标准 JSON 格式响应，与 REST 服务一致
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
