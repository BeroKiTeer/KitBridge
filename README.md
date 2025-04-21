# 🌉 KitBridge - Kitex 多协议融合网关实现说明

本项目为 CloudWeGo 2025 黑客松赛题《Kitex 多协议融合网关》的参考实现，目标是让 Kitex 服务端同时支持原生 Thrift RPC 与 HTTP/1.1 请求，并具备高性能、无侵入、标准化的协议融合能力。

## 技术目标

扩展Kitex框架，使其服务端能同时原生支持Thrift RPC和HTTP RESTful协议，需完成以下核心功能：

## 📌 实现目标对照表（核心功能要求）

| 编号 | 题目要求                             | 我们的实现方案                      | 完成状态 |
|------|--------------------------------------|-------------------------------------|----------|
| ✅ 1 | 协议智能识别                         | 使用 `detection.NewSvrTransHandlerFactory` 构建嗅探器工厂，结合 `ProtocolMatch()` 实现基于 `netpoll.Connection.Reader().Peek()` 的快速识别 | ✅ 已完成 |
| ✅ 2 | HTTP 路径映射到 Thrift 方法          | 支持标准 RESTful 映射 `/api/Service/Method`，在 `Read()` 中解析路径并填入 `msg.SetServiceName/Method` | ✅ 已完成 |
| ✅ 3 | JSON Body → Thrift 请求结构体        | `Read()` 中通过 `json.Unmarshal` 自动反序列化请求体为 Thrift 请求 struct | ✅ 已完成 |
| ✅ 4 | Header / Query 参数绑定              | IDL 中使用 `api.query` / `api.header`，即将实现 `binder.go` 解析 Tag 映射至 struct 字段 | ⏳ 规划中 |
| ✅ 5 | 错误处理标准化：JSON 三段式响应结构 | `Write()` 中规划封装为 `{code, message, data}` 响应，兼容 REST 错误语义 | ⏳ 待完成 |
| ✅ 6 | 保持 Kitex 中间件兼容性              | 使用 `SetInvokeHandleFunc()` 绑定框架 handler，确保中间件链保持完整 | ✅ 已完成 |
| ✅ 7 | Streaming 支持（Bonus）              | 当前未实现，但结构支持 HTTP chunk 扩展，未来可补充 server streaming | ⏳ 预留扩展点 |

#### 协议只能识别层

- 开发自定义TransHandlerFactory实现协议自动检测：
  - 根据TCP包头特征区分Thrift二进制协议与HTTP/1.1文本协议
  - 动态切换编解码处理器（Thrift Codec / HTTP Codec）
- 支持HTTP路径映射到Thrift方法（默认POST /api/UserService/GetUser映射到UserService::GetUser）

#### HTTP深度兼容实现

- 实现标准HTTP予以支持
    - 处理HTTP Header/Query到Thrift字段的映射关系
    - 自动转换JSON请求体与结构体
- 兼容Kitex middleware层
- 错误处理标准化
    -兼容Kitex业务一场错误，以统一错误响应JSON格式返回业务信息
- 支持kitex server streaming以http chunk格式接受/写入消息

## 项目结构

```
KitBridge/
├── cmd/                    # 应用启动入口
│   └── main.go            # 注册服务、注入组件、运行 server
│
├── idl/                   # Thrift 接口定义
│   └── st_service.thrift  # 示例 Thrift 服务定义（STService）
│
├── handler/               # 业务逻辑实现层
│   └── stservice_impl.go  # STService 的接口实现（testSTReq）
│
├── kitex_gen/             # kitex 自动生成的结构体和服务代码
│   └── thrift/stability/...
│
├── protocol/              # 协议栈适配层
│   ├── autodetect/        # 协议嗅探 + TransHandlerFactory
│   │   ├── trans_factory.go       # 核心注入入口（嗅探 Thrift / HTTP）
│   │   └── svr_trans_handler.go  # 实现 ProtocolMatch 与多处理器调度
│   │
│   └── codec/             # 可选：自定义 HTTP Codec 实现
│       └── http_codec.go
│
├── mapping/               # 路由映射（HTTP路径 → Thrift服务方法）
│   ├── path_mapper.go     # 解析 URL 如 /api/Service/Method
│   └── dispatcher.go      # 方法注册与分发（预留扩展）
│
├── transport/             # 请求体绑定层（HTTP参数 → Thrift结构体）
│   └── binding_adapter.go # 使用 hertz binding 实现自动映射
│
├── middleware/            # 中间件相关模块（上下文注入等）
│   └── context_injector.go
│
├── errors/                # 错误处理模块（统一输出 JSON 错误格式）
│   ├── encoder.go         # 输出封装
│   └── biz_error.go       # 自定义业务错误结构体定义
│
├── stream/                # Streaming 扩展模块（HTTP chunked）
│   └── http_stream.go     # 将 Kitex Stream 转为 chunk 响应
│
├── test/                  # 单元测试 & 压测工具
│   ├── path_mapper_test.go
│   └── curl_test.sh       # curl 混合请求脚本
│
├── configs/               # 配置文件
│   └── config.yaml        # 预留启动参数、路径映射配置等
│
├── go.mod
└── README.md

```

## ✅ 整体架构设计思路（Kitex协议融合网关）


```text
┌─────────────────────────────────────────────────────┐
│                    客户端（Thrift/HTTP）             │
├─────────────────────────────────────────────────────┤
│                网络层（netpoll Connection）         │
├─────────────────────────────────────────────────────┤
│    ProtocolDetect Layer（自定义 TransHandler）       │
│    ├─ Peek(4byte) -> 正则匹配是否 HTTP                 │
│    └─ 匹配结果注入不同 Codec (ThriftCodec / HTTPCodec)│
├─────────────────────────────────────────────────────┤
│        Dispatcher（路由分发）                         │
│        ├─ HTTP path 映射到 Thrift 方法名              │
│        └─ 将 JSON 请求体 绑定为 Thrift 对象           │
├─────────────────────────────────────────────────────┤
│           Codec 编解码层（Codec Layer）              │
│           ├─ thriftBinaryCodec                       │
│           └─ httpJSONCodec（需自定义）                │
├─────────────────────────────────────────────────────┤
│              服务逻辑处理层（Kitex handler）         │
├─────────────────────────────────────────────────────┤
│              Middleware (支持统一中间件机制)         │
└─────────────────────────────────────────────────────┘
```

