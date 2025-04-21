package http1

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/netpoll"
)

// 自定义错误类型
var (
	ErrInvalidRequestLine = errors.New("invalid HTTP request line")
	ErrInvalidPathFormat  = errors.New("invalid path format, expected /api/{Service}/{Method}")
	ErrInvalidContentLen  = errors.New("invalid Content-Length header")
)

// 读取一行（以 \n 结束）并去掉末尾的 \r\n
func readLine(reader netpoll.Reader) ([]byte, error) {
	var line []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		line = append(line, b)
		if b == '\n' {
			break
		}
	}
	// 去除 \r\n 或 \n
	line = bytes.TrimSuffix(line, []byte("\r\n"))
	line = bytes.TrimSuffix(line, []byte("\n"))
	return line, nil
}

// parseRequestLine 解析第一行请求行：如 POST /api/Service/Method HTTP/1.1
func parseRequestLine(reader netpoll.Reader) (method, serviceName, methodName string, err error) {
	line, err := readLine(reader)
	if err != nil {
		return "", "", "", err
	}

	parts := bytes.SplitN(line, []byte(" "), 3)
	if len(parts) < 3 {
		return "", "", "", ErrInvalidRequestLine
	}

	method = string(parts[0])
	path := string(parts[1])
	// version := string(parts[2]) // 可选保留

	pathParts := strings.Split(path, "/")
	if len(pathParts) < 4 || pathParts[1] != "api" {
		return "", "", "", ErrInvalidPathFormat
	}

	serviceName = pathParts[2]
	methodName = pathParts[3]
	return
}

// parseHeaders 解析 Header 字段，直到遇到空行 \r\n\r\n，返回 Header 映射和 Content-Length 值
func parseHeaders(reader netpoll.Reader) (map[string]string, int, error) {
	headers := make(map[string]string)
	var contentLength int

	for {
		line, err := readLine(reader)
		if err != nil {
			return nil, 0, err
		}
		if len(line) == 0 {
			break // 空行，header 结束
		}

		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(string(parts[0]))
		val := strings.TrimSpace(string(parts[1]))
		headers[key] = val

		if strings.EqualFold(key, "Content-Length") {
			cl, err := strconv.Atoi(val)
			if err != nil {
				return nil, 0, ErrInvalidContentLen
			}
			contentLength = cl
		}
	}

	return headers, contentLength, nil
}
