package http1

import (
	"bytes"
	"testing"

	"github.com/cloudwego/netpoll"
	"github.com/stretchr/testify/assert"
)

func TestParseRequestLine(t *testing.T) {
	// 模拟 HTTP 请求行
	input := []byte("POST /api/UserService/GetUser HTTP/1.1\r\n")

	// 构造 netpoll.Reader（使用 bytes.Reader 包装）
	r := netpoll.NewReader(bytes.NewReader(input))

	method, service, methodName, err := parseRequestLine(r)

	assert.NoError(t, err)
	assert.Equal(t, "POST", method)
	assert.Equal(t, "UserService", service)
	assert.Equal(t, "GetUser", methodName)
}

func TestParseRequestLine_InvalidPath(t *testing.T) {
	// 错误路径格式
	input := []byte("GET /wrongpath HTTP/1.1\r\n")
	r := netpoll.NewReader(bytes.NewReader(input))

	_, _, _, err := parseRequestLine(r)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPathFormat, err)
}

func TestReadLine(t *testing.T) {
	input := []byte("GET /test HTTP/1.1\r\n")
	r := netpoll.NewReader(bytes.NewReader(input))

	line, err := readLine(r)

	assert.NoError(t, err)
	assert.Equal(t, "GET /test HTTP/1.1", string(line))
}

func TestParseHeaders(t *testing.T) {
	input := []byte(
		"Host: localhost\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: 27\r\n" +
			"\r\n")
	reader := netpoll.NewReader(bytes.NewReader(input))

	headers, contentLength, err := parseHeaders(reader)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(headers))
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, 27, contentLength)
}

func TestParseHeaders_InvalidContentLength(t *testing.T) {
	input := []byte(
		"Host: localhost\r\n" +
			"Content-Length: notanumber\r\n" +
			"\r\n")
	reader := netpoll.NewReader(bytes.NewReader(input))

	_, _, err := parseHeaders(reader)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidContentLen, err)
}
