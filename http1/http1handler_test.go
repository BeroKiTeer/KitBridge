package http1_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"

	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/transport"
	"github.com/stretchr/testify/assert"

	"github.com/BeroKiTeer/KitBridge/http1"
	"github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability"
)

// mockMsg 实现 remote.Message 接口，支持 SetData/GetData
type mockMsg struct {
	remote.Message
	args interface{}
}

func (m *mockMsg) Data() interface{}                         { return m.args }
func (m *mockMsg) SetData(v interface{})                     { m.args = v }
func (m *mockMsg) SetPayloadCodec(codec remote.PayloadCodec) {}
func (m *mockMsg) SetMessageType(t remote.MessageType)       {}
func (m *mockMsg) TransInfo() remote.TransInfo               { return nil }
func (m *mockMsg) TransportProtocol() transport.Protocol     { return transport.HTTP }

func TestHTTP1Handler_Read(t *testing.T) {
	// 模拟 net.Conn
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	// 构造 handler 和 message
	handler := &http1.HTTP1Handler{}
	req := &stability.STRequest{}
	msg := &mockMsg{}
	msg.SetData(req)

	// 构造 HTTP 请求内容（含 query 参数）
	httpRequest := `POST /api/STService/testSTReq?a=Alice HTTP/1.1
		Host: localhost
		Content-Length: 0
		Content-Type: application/x-www-form-urlencoded
		
		`

	// 异步写入请求内容到 clientConn（模拟客户端）
	go func() {
		_, _ = io.Copy(clientConn, bytes.NewReader([]byte(httpRequest)))
	}()

	// 调用 Read 方法
	ctx := context.Background()
	_, err := handler.Read(ctx, serverConn, msg)

	// 验证绑定是否成功
	assert.NoError(t, err)
	assert.Equal(t, "Alice", req.Name)
}
