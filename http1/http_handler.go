package http1

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability"
	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/transmeta"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/netpoll"
	"net"
	"regexp"
)

// +------------------------------------------------------------+
// |                      TCP 连接进入                         |
// +------------------------------------------------------------+
//                             |
//                             v
// +------------------------------------------------------------+
// | detection.NewSvrTransHandlerFactory                        |  // ✅ 关键路径 1：协议嗅探工厂注册
// | - 注册 Thrift 默认处理器（netpoll）                         |
// | - 注册支持 ProtocolMatch 的 HTTP1HandlerFactory             |
// +------------------------------------------------------------+
//                             |
//                 +-----------+------------+
//                 |                        |
//       [Thrift 请求]             [HTTP 请求（如 POST /api/...）]
//                 |                        |
//                 v                        v
//     +--------------------+       +-------------------------+
//     | ThriftHandler      |       | HTTP1Handler            |   // ✅ 关键路径 2：ProtocolMatch 实现判断是否为 HTTP
//     +--------------------+       +-------------------------+
//                                         |
//                                         v
//   +-----------------------------------------------------------------+
//   | Read(ctx, conn, msg)                                            |  // ✅ 关键路径 3：解析 HTTP 请求体
//   | - 读取 HTTP 请求数据                                             |
//   | - 解析请求行 /api/Service/Method → 设置 msg.ServiceName/Method |
//   | - 解析 JSON Body → Thrift struct → 设置 msg.Args                |
//   | - TODO: 解析 Header/Query 参数 → 映射至字段（关键路径 5）         |
//   | - 设置 msg.MessageType = remote.Call                           |
//   | - 调用 h.handler(ctx, msg) 执行 Kitex 调用链（关键路径 7）       |
//   +-----------------------------------------------------------------+
//                                         |
//                                         v
//                        +----------------------------------+
//                        | Kitex RPC 服务逻辑（如 testSTReq）|
//                        +----------------------------------+
//                                         |
//                                         v
//   +-----------------------------------------------------------------+
//   | Write(ctx, conn, msg)                                           |  // ✅ 关键路径 4 & 6：统一响应封装
//   | - 判断是否为 BizError / 系统异常                                 |
//   | - 构造响应 JSON：{code, message, data}                         |
//   | - 设置 HTTP Header（Content-Type、Length）并写回 conn           |
//   +-----------------------------------------------------------------+

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

type JsonResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var httpPattern = regexp.MustCompile(`^(?:GET |POST|PUT|DELE|HEAD|OPTI|CONN|TRAC|PATC)$`)

type HTTP1SvrTransHandlerFactory struct{}

// NewTransHandler 是 Kitex 要求实现的工厂方法，用于创建一个 ServerTransHandler（即协议处理器）实例。
// opt 参数是框架在初始化阶段提供的服务上下文信息，包括服务结构、配置、结果工厂等。
func (f *HTTP1SvrTransHandlerFactory) NewTransHandler(opt *remote.ServerOption) (remote.ServerTransHandler, error) {
	return &HTTP1Handler{
		// 表示当前服务的元信息（如服务名、方法名、IDL 等）
		svcInfo: opt.TargetSvcInfo,
		// 服务元信息查找器
		svcSearcher: opt.SvcSearcher,
		// 保存整个服务配置上下文（包含 Payload 编解码器、错误处理器等）
		opt: opt,
	}, nil
}

type HTTP1Handler struct {
	// 当前服务的元信息（服务名、方法名、IDL 结构）
	svcInfo *serviceinfo.ServiceInfo
	// 按服务名查找 ServiceInfo 的查找器
	svcSearcher remote.ServiceSearcher
	// 完整的 ServerOption 配置上下文，用于读取 Codec、错误处理、ResultProvider 等
	opt *remote.ServerOption
	// 在 SetPipeline() 中注入，用于调度 Read → OnMessage → Write 的框架处理管道
	transPipe   *remote.TransPipeline
	handlerFunc endpoint.Endpoint
}

func (h *HTTP1Handler) ProtocolMatch(ctx context.Context, conn net.Conn) error {
	c, ok := conn.(netpoll.Connection)
	if ok {
		pre, _ := c.Reader().Peek(4)
		if httpPattern.Match(pre) {
			return nil
		}
	}
	return errors.New("error protocol not match")
}

// 解析 HTTP 请求并转为 Kitex RPC 调用
func (h *HTTP1Handler) Read(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error) {
	// ---------------------------------------------------------
	// 1: 利用 parser.parseRequestLine() 解析请求行
	// - 获取 method / serviceName / methodName
	// - 校验 path 格式为 /api/{Service}/{Method}
	// ---------------------------------------------------------
	// 创建带缓冲的 Reader，并封装为 netpoll.Reader
	bufReader := bufio.NewReader(conn)
	reader := netpoll.NewReader(bufReader)

	// 使用 parser.go 中的 parseRequestLine 方法
	method, serviceName, methodName, err := parseRequestLine(reader)
	klog.Infof("HTTP1 request line parsed: method=%s, service=%s, method=%s", method, serviceName, methodName)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse request line: %w", err)
	}
	// ---------------------------------------------------------
	// 2: 利用 parser.parseHeaders() 获取 Header Map
	// - Content-Length 字段必须存在且合法
	// - 所有 header 存入 headers map[string]string
	// ---------------------------------------------------------
	_, contentLength, err := parseHeaders(reader)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse headers: %w", err)
	}
	// ---------------------------------------------------------
	// 3: 利用 reader.Next(n) 精准读取 JSON body
	// - Content-Length 决定 body 大小
	// - 返回值为 []byte 类型 JSON
	// ---------------------------------------------------------
	if contentLength <= 0 || contentLength > 10*1024*1024 { // 限制最大 10MB
		return ctx, fmt.Errorf("invalid content length: %d", contentLength)
	}
	bodyBytes, err := reader.Next(contentLength)
	if err != nil {
		return ctx, fmt.Errorf("failed to read body: %w", err)
	}
	// ---------------------------------------------------------
	// 4: 将 path/header 映射为 Kitex 元信息
	msg.SetMessageType(remote.Call)

	// 将 service 和 method 注入到 TransInfo 中
	msg.TransInfo().PutTransIntInfo(map[uint16]string{
		transmeta.ToService: serviceName,
		transmeta.ToMethod:  methodName,
	})

	// 将 header 中的一些字段透传（例如 X-Trace-ID）
	if metainfo.HasMetaInfo(ctx) {
		metaMap := make(map[string]string)
		metainfo.SaveMetaInfoToMap(ctx, metaMap)
		msg.TransInfo().PutTransStrInfo(metaMap)
	}

	// 5: JSON body → Thrift 请求 struct
	args := msg.Data().(*stability.STRequest)

	// JSON 反序列化
	if err := json.Unmarshal(bodyBytes, args); err != nil {
		return ctx, fmt.Errorf("failed to unmarshal body to thrift args: %w", err)
	}

	return ctx, nil
}

// 将 Kitex RPC 返回结果封装为标准 HTTP JSON 响应
func (h *HTTP1Handler) Write(ctx context.Context, conn net.Conn, msg remote.Message) (context.Context, error) {
	// 1: 判断调用结果是正常返回还是异常
	var (
		code    int32
		message string
		data    interface{}
	)
	rpcInfo := msg.RPCInfo()
	if bizErr := rpcInfo.Invocation().BizStatusErr(); bizErr != nil {
		code = bizErr.BizStatusCode()
		message = bizErr.BizMessage()
		data = nil
	} else if sysErr := rpcInfo.Stats().Error(); sysErr != nil {
		code = 500
		message = "internal error"
		data = nil
	} else {
		code = 200
		message = "success"
		data = msg.Data()
	}

	// 2: 构造标准 JSON 响应结构：
	resp := JsonResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}

	// 3: 根据错误类型构造不同 code/message：
	jsonBody, err := json.Marshal(resp)
	if err != nil {
		// 如果 JSON 编码失败（极少见，一般是结构体含非法类型）
		// 构造兜底 JSON 响应，防止崩溃
		jsonBody = []byte(`{"code":500,"message":"json encode error","data":null}`)
	}

	// 4: 构造 HTTP 响应头 响应头：HTTP/1.1 200 OK + Content-Type: application/json + Content-Length
	var buf bytes.Buffer

	// 写响应行和头部
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	buf.WriteString("Content-Type: application/json\r\n")
	buf.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(jsonBody)))
	buf.WriteString("Connection: keep-alive\r\n")

	// 空行分隔 header 和 body
	buf.WriteString("\r\n")
	// 5: 构造完整 HTTP 响应字符串
	// - 响应体：json 数据
	buf.Write(jsonBody)

	// 6: 写入 conn（用 conn.Write(...) 输出响应）
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	// 最终效果：HTTP 客户端收到标准 JSON 格式响应，与 REST 服务一致
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
