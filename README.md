# KitBridge


## 描述

这是Kitext的插件

## 技术目标

扩展Kitex框架，使其服务端能同时原生支持Thrift RPC和HTTP RESTful协议，需完成以下核心功能：

### 核心功能要求

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

