# OctoHub WebSocket 服务器

这是一个支持用户端和服务端节点连接的WebSocket服务器，具有签名验证和连接管理功能。

## 功能特性

- **双端连接支持**: 支持用户端(user_id)和服务端节点(pc_id)连接
- **签名验证**: 使用HMAC-SHA256算法验证连接的合法性
- **单连接限制**: 每个ID只允许一个活跃连接，新连接会断开旧连接
- **断开通知**: 当连接被新连接替换时会发送通知消息
- **连接监控**: 提供连接状态查询接口
- **心跳检测**: 客户端ping，服务端pong响应机制保持连接活跃
- **模块化设计**: 采用清晰的分层架构，易于维护和扩展

## 项目结构

```
octo-hub-ws/
├── main.go                          # 主入口文件
├── go.mod                           # Go模块文件
├── config.yaml                      # 配置文件
├── README.md                        # 项目文档
├── examples/
│   └── signature_example.go         # 签名生成示例
└── internal/                        # 内部模块
    ├── auth/
    │   └── signature.go             # 签名验证模块
    ├── config/
    │   └── config.go                # 配置管理模块
    ├── connection/
    │   ├── types.go                 # 连接类型定义
    │   └── manager.go               # 连接管理器
    ├── handler/
    │   └── websocket.go             # WebSocket消息处理
    └── server/
        └── server.go                # WebSocket服务器
```

## 安装依赖

```bash
go mod tidy
```

## 运行服务器

```bash
go run main.go
```

服务器将在端口8080上启动。

## API 接口

### WebSocket 连接

#### 用户端连接
```
ws://localhost:8080/ws/user?user_id={用户ID}&timestamp={时间戳}&signature={签名}
```

#### 服务端节点连接
```
ws://localhost:8080/ws/pc?pc_id={PC ID}&timestamp={时间戳}&signature={签名}
```

### HTTP 接口

#### 连接状态查询
```
GET /status
```
返回当前所有连接的状态信息。

#### 健康检查
```
GET /health
```
返回服务器健康状态。

## 签名算法

签名使用HMAC-SHA256算法生成：

1. **消息**: `{id}{timestamp}` (ID + 时间戳)
2. **密钥**: 服务器配置的签名密钥
3. **算法**: HMAC-SHA256
4. **编码**: 十六进制字符串

### 签名生成示例

```bash
# 运行签名生成示例
go run examples/signature_example.go
```

这将生成用户和PC连接的签名示例。

## 连接管理

### 单连接限制
- 每个`user_id`或`pc_id`只允许一个活跃连接
- 当同一ID发起新连接时，旧连接会被自动断开
- 断开时会向旧连接发送通知消息

### 断开通知格式
```json
{
  "type": "disconnect_notification",
  "reason": "新连接建立",
  "timestamp": 1632123456
}
```

### 消息格式

#### 回显消息示例
```json
{
  "type": "echo",
  "from": "user123",
  "message": "原始消息内容",
  "timestamp": 1632123456
}
```

## 配置

服务器支持多种配置方式，优先级从高到低：环境变量 > 配置文件 > 默认值

### 配置文件

服务器读取当前目录下的 `config.yaml` 配置文件：

```yaml
# 服务器配置
server:
  port: "8080"

# 签名验证配置
signature:
  key: "your-production-secret-key"
  timeout: 300                        # 签名有效期（秒）

# WebSocket 配置
websocket:
  read_timeout: 60                    # 读取超时时间（秒）
  write_timeout: 10                   # 写入超时时间（秒）
  max_message_size: 1048576           # 最大消息大小（字节）

# 连接管理配置
connection:
  max_connections: 10000              # 最大连接数
  buffer_size: 256                    # 消息缓冲区大小

# 日志配置
logging:
  level: "info"                       # 日志级别
  format: "text"                      # 日志格式
```

### 环境变量配置

环境变量使用 `OCTOHUB_` 前缀，配置项用下划线分隔：

```bash
# 基本配置
export OCTOHUB_SERVER_PORT="8080"
export OCTOHUB_SIGNATURE_KEY="your-secret"

# 高级配置
export OCTOHUB_SIGNATURE_TIMEOUT="300"
export OCTOHUB_WEBSOCKET_READ_TIMEOUT="60"
export OCTOHUB_WEBSOCKET_WRITE_TIMEOUT="10"
export OCTOHUB_CONNECTION_MAX_CONNECTIONS="10000"
export OCTOHUB_CONNECTION_BUFFER_SIZE="256"

# 日志配置
export OCTOHUB_LOGGING_LEVEL="debug"
export OCTOHUB_LOGGING_FORMAT="json"
```

### 配置说明

项目使用单一的 `config.yaml` 配置文件，简单易维护。如果配置文件不存在，服务器将使用默认配置启动。

### 签名验证
- 默认签名有效期：5分钟
- 超过有效期的连接请求将被拒绝
- 生产环境必须设置自定义的签名密钥

### 心跳机制
服务器采用**客户端主动ping**的策略：

- **客户端职责**: 定期发送ping消息保持连接活跃
- **服务端职责**: 自动响应pong消息，检测连接超时
- **优势**: 
  - 减少服务端资源消耗（特别是高并发场景）
  - 客户端更容易检测网络断开
  - 简化服务端逻辑，提高性能

**建议客户端ping间隔**: 30-50秒
**服务端读取超时**: 60秒（可配置）

## 测试

### 使用 WebSocket 客户端测试

1. 运行签名生成器获取有效的连接URL：
   ```bash
   go run signature_example.go
   ```

2. 使用WebSocket客户端（如wscat）连接：
   ```bash
   # 安装 wscat
   npm install -g wscat
   
   # 连接用户端
   wscat -c "ws://localhost:8080/ws/user?user_id=user123&timestamp=1632123456&signature=generated_signature"
   
   # 连接PC端
   wscat -c "ws://localhost:8080/ws/pc?pc_id=pc456&timestamp=1632123456&signature=generated_signature"
   ```

3. 查看连接状态：
   ```bash
   curl http://localhost:8080/status
   ```

## 日志输出

服务器会输出以下类型的日志：
- 连接建立日志
- 连接断开日志
- 消息接收日志
- 错误日志

## 扩展功能

当前实现提供了基础的连接管理功能，可以根据需要扩展：

1. **消息路由**: 在`handleMessage`函数中添加消息路由逻辑
2. **数据持久化**: 添加连接信息的数据库存储
3. **集群支持**: 使用Redis等实现多实例连接状态同步
4. **监控指标**: 添加Prometheus指标收集
5. **负载均衡**: 配置反向代理实现负载均衡

## 安全注意事项

1. **签名密钥**: 确保签名密钥的安全性，定期更换
2. **HTTPS**: 生产环境建议使用WSS(WebSocket Secure)
3. **速率限制**: 考虑添加连接频率限制
4. **输入验证**: 对所有输入参数进行严格验证
