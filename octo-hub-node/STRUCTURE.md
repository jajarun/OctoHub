# 项目结构说明

## 模块化重构完成

原来的单文件 `main.py` (306行) 已经重构为模块化结构，代码更加清晰和易于维护。

## 目录结构

```
octo-hub-node/
├── main.py                    # 主入口文件 (简化至65行)
├── config.py                  # 配置管理
├── config.example             # 环境变量配置示例
├── start.sh                   # 启动脚本
├── README.md                  # 详细文档
├── STRUCTURE.md               # 本文件
├── pyproject.toml             # 项目依赖
│
├── octohub/                   # 核心包
│   ├── __init__.py
│   │
│   ├── client/                # 客户端模块
│   │   ├── __init__.py
│   │   └── websocket_client.py    # WebSocket客户端核心类
│   │
│   ├── auth/                  # 认证模块
│   │   ├── __init__.py
│   │   └── signature_auth.py      # HMAC-SHA256签名认证
│   │
│   ├── handlers/              # 处理器模块
│   │   ├── __init__.py
│   │   └── message_handler.py     # WebSocket消息处理器
│   │
│   ├── tasks/                 # 任务处理模块
│   │   ├── __init__.py
│   │   └── task_processor.py      # 任务处理器和调度
│   │
│   └── utils/                 # 工具模块
│       ├── __init__.py
│       ├── id_generator.py        # PC ID生成器
│       └── logger.py              # 日志工具
│
└── examples/                  # 使用示例
    ├── __init__.py
    ├── simple_client.py          # 简单使用示例
    └── custom_client.py          # 完整自定义客户端示例
```

## 模块功能说明

### 1. 客户端模块 (`octohub/client/`)
- **websocket_client.py**: WebSocket客户端核心类
  - 连接管理和自动重连
  - 消息发送和接收
  - 处理器注册接口

### 2. 认证模块 (`octohub/auth/`)
- **signature_auth.py**: 签名认证处理
  - HMAC-SHA256签名生成
  - API请求认证
  - WebSocket连接地址获取

### 3. 处理器模块 (`octohub/handlers/`)
- **message_handler.py**: 消息处理器
  - WebSocket消息解析和路由
  - 内置消息类型处理
  - 自定义处理器注册

### 4. 任务处理模块 (`octohub/tasks/`)
- **task_processor.py**: 任务处理器
  - 任务执行和结果返回
  - 内置任务类型（echo, sleep等）
  - 自定义任务处理器注册

### 5. 工具模块 (`octohub/utils/`)
- **id_generator.py**: PC ID生成器
- **logger.py**: 日志配置工具

## 使用方式

### 基本使用
```python
from octohub.client import NodeWebSocketClient

client = NodeWebSocketClient()
await client.start()
```

### 自定义任务处理
```python
async def my_task_handler(task_data):
    return "处理完成"

client.register_task_handler("my_task", my_task_handler)
```

### 自定义消息处理
```python
async def my_message_handler(message_data):
    print("收到自定义消息")

client.register_message_handler("my_message", my_message_handler)
```

## 优势

1. **模块化**: 代码按功能分离，易于维护和扩展
2. **可复用**: 各模块可独立使用和测试
3. **可扩展**: 支持自定义处理器注册
4. **清晰性**: 主文件简化至65行，逻辑清晰
5. **灵活性**: 支持多种使用方式和配置选项

## 向后兼容

原有的使用方式仍然有效，只需要更新导入路径：

```python
# 旧方式
from main import NodeWebSocketClient

# 新方式
from octohub.client import NodeWebSocketClient
```
