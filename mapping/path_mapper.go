package mapping

import (
	"errors"
	"strings"
)

var (
	ErrInvalidPath = errors.New("invalid API path format, should be /api/Service/Method")
)

// PathMapping 表示 HTTP 请求路径中提取出的 Service/Method
type PathMapping struct {
	Service string
	Method  string
}

// ParsePath 从 HTTP 路径中解析出 Service 与 Method
func ParsePath(path string) (*PathMapping, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 4 || parts[1] != "api" {
		return nil, ErrInvalidPath
	}

	return &PathMapping{
		Service: parts[2],
		Method:  parts[3],
	}, nil
}
