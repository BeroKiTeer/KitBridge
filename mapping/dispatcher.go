package mapping

import (
	"fmt"
)

// 暂存注册的方法表（服务名 + 方法名）
var registeredMethods = make(map[string]func())

// RegisterMock 注册 mock handler（可用于测试）
func RegisterMock(service string, method string, handler func()) {
	key := fmt.Sprintf("%s::%s", service, method)
	registeredMethods[key] = handler
}

// Dispatch 调用绑定方法
func Dispatch(service string, method string) {
	key := fmt.Sprintf("%s::%s", service, method)
	if handler, ok := registeredMethods[key]; ok {
		handler()
	} else {
		fmt.Println("404 Not Found: ", key)
	}
}
