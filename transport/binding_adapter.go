package transport

import (
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

var ErrBinding = errors.New("BindToStruct failed")

// BindToStruct 从 Header + Query + JSON 三个来源解析 HTTP 请求，填充到 Thrift 结构体中
func BindToStruct(c *app.RequestContext, obj interface{}) error {
	// 1. Query 参数绑定（支持 go.tag = "query:\"xxx\""）
	if err := c.BindQuery(obj); err != nil {
		debugLog("Query bind failed: %v", err)
	}

	// 2. Header 参数绑定（支持 go.tag = "header:\"xxx\""）
	if err := c.BindHeader(obj); err != nil {
		debugLog("Header bind failed: %v", err)
	}

	// 3. JSON Body 参数绑定（POST请求 + Content-Type: application/json）
	if c.Request.Header.ContentLength() > 0 &&
		string(c.Request.Header.ContentType()) == "application/json" {

		if err := c.BindJSON(obj); err != nil {
			return fmt.Errorf("%w: JSON bind failed: %v", ErrBinding, err)
		}
	}

	return nil
}

func debugLog(format string, args ...interface{}) {
	// 如果你接了 zap / logrus / slog，可以换成你自己的日志系统
	fmt.Printf("[KitBridge:BindDebug] "+format+"\n", args...)
}
