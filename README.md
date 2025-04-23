# KitBridge：Kitex 多协议融合网关

🚀 **CloudWeGo 2025 黑客松一等奖项目**
 兼容 RESTful 与 Thrift 的高性能微服务网关扩展

## 🔍 背景与动机

在微服务架构中，服务间通信多采用高性能的 Thrift / gRPC 等 RPC 协议，但面向公网的 API 通常仍需提供 RESTful 接口，造成：

- 重复维护两套接口逻辑
- REST 网关接入链路复杂，转发性能开销高
- 协议迁移难度大，阻碍架构演进

本项目旨在打造一套轻量级、高性能的 **协议融合解决方案**，帮助开发者在 **不改动现有业务代码的前提下**，实现 REST API 到 Thrift 的无缝映射与调度。

------

## 🌟 项目亮点

### ✅ 协议智能识别与动态编解码切换

- 基于 TCP 首部特征实现 **HTTP1.1 / Thrift 协议自动识别**
- 使用 `detection.NewSvrTransHandlerFactory` 注入协议嗅探层
- 首次请求自动绑定合适的 handler，**一连接一协议，终身绑定，零损耗切换**

### ✅ HTTP 深度兼容（REST → Thrift）

- 支持标准 `POST /api/{Service}/{Method}` 路径映射到 Thrift 方法
- 实现 Header / Query 参数 → Thrift 字段映射
- JSON Body 自动反序列化为 Thrift 请求结构体
- 自定义 `TransHandler`，兼容 Kitex 中间件与服务注册机制
- 返回统一格式 JSON 响应 `{ code, message, data }`

### ✅ 插件式集成，零侵入

- 遵循 Kitex v0.13+ 插件扩展机制
- 可与现有 Thrift 客户端 100% 兼容，无需调整
- 仅需一行 Option 注入即可启用 HTTP 支持

------

## 🚀 适用场景

**适用：**

- 希望从 REST 平滑迁移到 RPC 的业务系统
- 面向公网同时暴露 REST 与 RPC 的 API 网关
- 快速上线的中小型后端服务，追求统一服务定义、减少重复开发

**不适用：**

- 高度复杂的 API 设计（如 `/users/{id}/posts` 风格）
- 对于安全性和协议隔离要求极高的生产环境（建议加上 API 网关如 APISIX 做隔离）
- 已有成熟 BFF 层的架构 

------

## 🔧 使用方法

```go
import (
    "github.com/BeroKiTeer/KitBridge/http1"
    "github.com/BeroKiTeer/KitBridge/autodetect"
    "github.com/cloudwego/kitex/server"
)

func main() {
    opts := []server.Option{
        server.WithTransHandlerFactory(
            autodetect.NewSvrTransHandlerFactoryWithHTTP(&http1.HTTP1SvrTransHandlerFactory{}),
        ),
    }

    svr := myService.NewServer(new(MyServiceImpl), opts...)
    svr.Run()
}
```

------

## 📊 未来展望

- 支持 HTTP2 / gRPC 协议融合
- 方法参数自动注册与反射机制，去除硬编码解耦
- 结合 Jaeger / Prometheus 监控链路打通

------

## 👏 致谢

本项目为 CloudWeGo 2025 黑客松一等奖作品，感谢 CloudWeGo 社区提供的开放平台与指导！