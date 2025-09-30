"""
WebSocket客户端核心模块
"""

import asyncio
import json
from typing import Optional, Dict, Any, Callable, Awaitable
try:
    import websockets
    from websockets.exceptions import ConnectionClosed, InvalidStatusCode
except ImportError:
    # websockets包未安装时的占位符
    websockets = None
    ConnectionClosed = Exception
    InvalidStatusCode = Exception

from ..auth.signature_auth import SignatureAuth
from ..handlers.message_handler import MessageHandler
from ..utils.id_generator import generate_pc_id
from ..utils.logger import setup_logger
try:
    from config import config
except ImportError:
    # 如果作为模块导入时找不到config，使用相对导入
    import sys
    import os
    sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(__file__))))
    from config import config

logger = setup_logger(__name__)


class NodeWebSocketClient:
    """OctoHub Node WebSocket 客户端"""
    
    def __init__(self, 
                 server_host: str = None,
                 server_port: int = None,
                 ws_port: int = None,
                 signature_key: str = None,
                 reconnect_interval: int = None,
                 max_reconnect_attempts: int = None,
                 pc_id: str = None):
        
        # 使用配置文件或传入的参数
        self.server_host = server_host or config.SERVER_HOST
        self.server_port = server_port or config.SERVER_PORT
        self.ws_port = ws_port or config.WS_PORT
        self.signature_key = signature_key or config.SIGNATURE_KEY
        self.reconnect_interval = reconnect_interval or config.RECONNECT_INTERVAL
        self.max_reconnect_attempts = max_reconnect_attempts or config.MAX_RECONNECT_ATTEMPTS
        
        # 生成或使用指定的PC ID
        self.pc_id = pc_id or generate_pc_id()
        
        # 初始化组件
        self.auth = SignatureAuth(self.signature_key)
        self.message_handler = MessageHandler(self.pc_id, self._send_message)
        
        # WebSocket连接相关
        self.websocket: Optional[websockets.WebSocketServerProtocol] = None
        self.is_running = False
        self.reconnect_attempts = 0
        
        logger.info(f"初始化 NodeWebSocketClient，PC ID: {self.pc_id}")
    
    async def _connect_websocket(self) -> bool:
        """建立WebSocket连接"""
        try:
            ws_url = await self.auth.get_websocket_url(
                self.pc_id, 
                self.server_host, 
                self.server_port
            )
            logger.info(f"尝试连接WebSocket: {ws_url}")
            
            # 建立WebSocket连接
            self.websocket = await websockets.connect(
                ws_url,
                ping_interval=config.PING_INTERVAL,  # ping间隔
                ping_timeout=config.PING_TIMEOUT,   # ping超时
                close_timeout=10   # 10秒关闭超时
            )
            
            logger.info("WebSocket连接建立成功")
            self.reconnect_attempts = 0
            return True
            
        except Exception as e:
            logger.error(f"WebSocket连接失败: {e}")
            return False
    
    async def _send_message(self, message: Dict[str, Any]):
        """发送消息"""
        if self.websocket and not self.websocket.closed:
            try:
                message_str = json.dumps(message, ensure_ascii=False)
                await self.websocket.send(message_str)
                logger.debug(f"发送消息: {message_str}")
            except Exception as e:
                logger.error(f"发送消息失败: {e}")
        else:
            logger.warning("WebSocket未连接，无法发送消息")
    
    async def _listen_messages(self):
        """监听WebSocket消息"""
        try:
            async for message in self.websocket:
                await self.message_handler.handle_message(message)
        except ConnectionClosed:
            logger.warning("WebSocket连接已关闭")
        except Exception as e:
            logger.error(f"监听消息时发生错误: {e}")
    
    async def _reconnect(self):
        """自动重连逻辑"""
        while self.is_running and self.reconnect_attempts < self.max_reconnect_attempts:
            self.reconnect_attempts += 1
            
            logger.info(f"尝试重连 ({self.reconnect_attempts}/{self.max_reconnect_attempts})")
            
            # 等待重连间隔
            await asyncio.sleep(self.reconnect_interval)
            
            # 尝试重新连接
            if await self._connect_websocket():
                logger.info("重连成功")
                # 重新开始监听消息
                await self._listen_messages()
                break
            else:
                logger.warning(f"重连失败，{self.reconnect_interval}秒后重试")
        
        if self.reconnect_attempts >= self.max_reconnect_attempts:
            logger.error(f"达到最大重连次数 ({self.max_reconnect_attempts})，停止重连")
            self.is_running = False
    
    async def start(self):
        """启动WebSocket客户端"""
        logger.info(f"启动 OctoHub Node 客户端，PC ID: {self.pc_id}")
        self.is_running = True
        
        # 初始连接
        if not await self._connect_websocket():
            logger.error("初始连接失败，开始重连流程")
            await self._reconnect()
            return
        
        try:
            # 开始监听消息
            await self._listen_messages()
        except Exception as e:
            logger.error(f"消息监听过程中发生错误: {e}")
        
        # 如果连接断开且客户端仍在运行，尝试重连
        if self.is_running:
            logger.info("连接断开，开始重连流程")
            await self._reconnect()
    
    async def stop(self):
        """停止WebSocket客户端"""
        logger.info("停止 OctoHub Node 客户端")
        self.is_running = False
        
        if self.websocket and not self.websocket.closed:
            await self.websocket.close()
    
    def register_message_handler(self, 
                                message_type: str, 
                                handler: Callable[[Dict[str, Any]], Awaitable[None]]):
        """注册自定义消息处理器"""
        self.message_handler.register_message_handler(message_type, handler)
    
    def register_task_handler(self, 
                             task_type: str, 
                             handler: Callable[[Dict[str, Any]], Awaitable[Any]]):
        """注册自定义任务处理器"""
        self.message_handler.register_task_handler(task_type, handler)
    
    def get_pc_id(self) -> str:
        """获取PC ID"""
        return self.pc_id
    
    def is_connected(self) -> bool:
        """检查是否已连接"""
        return self.websocket is not None and not self.websocket.closed
    
    def get_task_processor(self):
        """获取任务处理器"""
        return self.message_handler.get_task_processor()
