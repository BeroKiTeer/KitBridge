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
├── go.mod
├── go.sum
├── idl/
├── protocol/
├── server/
└── client/
```