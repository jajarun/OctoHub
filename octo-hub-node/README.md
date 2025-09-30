# OctoHub Node Client

OctoHub Node 客户端是一个基于Python的WebSocket客户端，用于连接OctoHub服务器并处理任务分发。

## 功能特性

- ✅ 自动生成唯一的PC ID标识符（基于主机名和MAC地址）
- ✅ 从服务器API获取WebSocket连接地址
- ✅ 支持签名认证（HMAC-SHA256）
- ✅ WebSocket连接和消息监听
- ✅ 自动重连机制
- ✅ 任务处理和响应
- ✅ 完整的日志记录
- ✅ 环境变量配置支持

## 安装依赖

```bash
# 安装项目依赖
pip install -e .

# 或者手动安装依赖
pip install websockets>=12.0 aiohttp>=3.9.0
```

## 配置

### 配置方式

支持多种配置方式，优先级：`.env文件` > `环境变量` > `默认值`

#### 1. 使用 .env 文件（推荐）

复制 `env.example` 为 `.env` 文件并修改配置：

```bash
cp env.example .env
```

编辑 `.env` 文件：
```bash
# 服务器配置
OCTOHUB_SERVER_HOST=localhost
OCTOHUB_SERVER_PORT=8080
OCTOHUB_WS_PORT=8000

# 签名认证配置（必须与服务端一致）
OCTOHUB_SIGNATURE_KEY=mySignatureKey123456789012345

# WebSocket配置
OCTOHUB_RECONNECT_INTERVAL=5
OCTOHUB_MAX_RECONNECT_ATTEMPTS=10
OCTOHUB_PING_INTERVAL=30
OCTOHUB_PING_TIMEOUT=10

# 日志配置
OCTOHUB_LOG_LEVEL=INFO
```

#### 2. 使用环境变量

```bash
export OCTOHUB_SERVER_HOST="localhost"
export OCTOHUB_SIGNATURE_KEY="mySignatureKey123456789012345"
# ... 其他配置
```

#### 3. 配置文件

配置逻辑在 `config.py` 文件中，会自动加载 `.env` 文件。

## 使用方法

### 基本使用

```bash
# 直接运行
python main.py
```

### 高级使用

```python
import asyncio
from octohub.client import NodeWebSocketClient

async def run_client():
    # 创建客户端实例
    client = NodeWebSocketClient(
        server_host="your-server-host",
        server_port=8080,
        signature_key="your-signature-key"
    )
    
    # 注册自定义任务处理器
    async def handle_custom_task(task_data):
        # 处理自定义任务逻辑
        return "任务处理完成"
    
    client.register_task_handler("custom_task", handle_custom_task)
    
    try:
        await client.start()
    except KeyboardInterrupt:
        await client.stop()

# 运行客户端
asyncio.run(run_client())
```

### 模块化结构

```
octohub/
├── __init__.py
├── client/
│   ├── __init__.py
│   └── websocket_client.py    # WebSocket客户端核心类
├── auth/
│   ├── __init__.py
│   └── signature_auth.py      # 签名认证模块
├── handlers/
│   ├── __init__.py
│   └── message_handler.py     # 消息处理器
├── tasks/
│   ├── __init__.py
│   └── task_processor.py      # 任务处理器
└── utils/
    ├── __init__.py
    ├── id_generator.py        # PC ID生成器
    └── logger.py              # 日志工具
```

## 工作流程

1. **初始化**: 生成唯一的PC ID，基于主机名和MAC地址
2. **获取连接地址**: 通过HTTP API从Spring Boot服务器获取WebSocket连接地址
3. **签名认证**: 使用HMAC-SHA256算法生成签名进行身份验证
4. **建立连接**: 连接到WebSocket服务器
5. **消息监听**: 监听服务器下发的任务和消息
6. **任务处理**: 处理接收到的任务并返回结果
7. **自动重连**: 连接断开时自动尝试重连

## 消息格式

### 发送消息格式

```json
{
  "type": "node_ready",
  "pc_id": "node_abc123def456",
  "timestamp": 1632123456
}
```

### 接收任务格式

```json
{
  "type": "task",
  "task_id": "task_123",
  "task_type": "example_task",
  "data": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

### 任务结果格式

```json
{
  "type": "task_result",
  "task_id": "task_123",
  "status": "completed",
  "result": "任务执行结果",
  "timestamp": 1632123456
}
```

## 日志示例

```
2024-01-01 12:00:00,000 - __main__ - INFO - 初始化 NodeWebSocketClient，PC ID: node_abc123def456
2024-01-01 12:00:01,000 - __main__ - INFO - 获取到WebSocket连接地址: ws://localhost:8000/ws/node?pc_id=node_abc123def456&timestamp=1632123456&signature=...
2024-01-01 12:00:02,000 - __main__ - INFO - WebSocket连接建立成功
2024-01-01 12:00:03,000 - __main__ - INFO - 收到消息类型: connected
2024-01-01 12:00:04,000 - __main__ - INFO - 连接成功确认
```

## 错误处理

- **连接失败**: 自动重连机制，可配置重连间隔和最大重连次数
- **签名验证失败**: 检查签名密钥是否与服务端一致
- **网络异常**: 自动检测连接状态并重连
- **任务处理异常**: 捕获异常并返回错误状态

## 开发说明

### 扩展任务处理

使用模块化的任务处理器，轻松扩展功能：

```python
from octohub.client import NodeWebSocketClient

# 创建客户端
client = NodeWebSocketClient()

# 注册自定义任务处理器
async def handle_file_task(task_data):
    file_path = task_data.get("data", {}).get("file_path")
    # 处理文件操作
    return f"文件 {file_path} 处理完成"

client.register_task_handler("file_operation", handle_file_task)
```

### 自定义消息处理

注册自定义消息处理器：

```python
async def handle_notification(message_data):
    notification = message_data.get("notification")
    print(f"收到通知: {notification}")

client.register_message_handler("notification", handle_notification)
```

### 使用示例

查看 `examples/` 目录中的示例：

- `examples/simple_client.py` - 简单使用示例
- `examples/custom_client.py` - 完整的自定义客户端示例

## 注意事项

1. **签名密钥**: 生产环境必须设置 `OCTOHUB_SIGNATURE_KEY` 环境变量
2. **网络配置**: 确保客户端能够访问Spring Boot服务器和WebSocket服务器
3. **防火墙**: 确保相关端口（默认8080和8000）已开放
4. **时间同步**: 签名验证依赖时间戳，确保客户端和服务端时间同步

## 故障排除

### 连接失败
- 检查服务器地址和端口配置
- 确认服务器正在运行
- 检查网络连接和防火墙设置

### 签名验证失败
- 确认签名密钥与服务端一致
- 检查时间同步
- 查看详细日志获取更多信息

### 重连失败
- 检查网络稳定性
- 调整重连间隔和最大重连次数
- 查看服务器日志确认连接限制

## 许可证

MIT License
