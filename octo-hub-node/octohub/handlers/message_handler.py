"""
消息处理器模块
"""

import json
import time
from typing import Dict, Any, Callable, Awaitable, Optional

from ..tasks.task_processor import TaskProcessor
from ..utils.logger import setup_logger

logger = setup_logger(__name__)


class MessageHandler:
    """WebSocket消息处理器"""
    
    def __init__(self, pc_id: str, send_message_callback: Callable[[Dict[str, Any]], Awaitable[None]]):
        self.pc_id = pc_id
        self.send_message = send_message_callback
        self.task_processor = TaskProcessor()
        self.custom_handlers: Dict[str, Callable[[Dict[str, Any]], Awaitable[None]]] = {}
    
    def register_message_handler(self, 
                                message_type: str, 
                                handler: Callable[[Dict[str, Any]], Awaitable[None]]):
        """注册自定义消息处理器"""
        self.custom_handlers[message_type] = handler
        logger.info(f"注册消息处理器: {message_type}")
    
    def register_task_handler(self, 
                             task_type: str, 
                             handler: Callable[[Dict[str, Any]], Awaitable[Any]]):
        """注册任务处理器"""
        self.task_processor.register_handler(task_type, handler)
    
    async def handle_message(self, message: str):
        """处理接收到的消息"""
        try:
            data = json.loads(message)
            message_type = data.get("type", "unknown")
            
            logger.info(f"收到消息类型: {message_type}")
            logger.debug(f"消息内容: {data}")
            
            # 检查是否有自定义处理器
            if message_type in self.custom_handlers:
                await self.custom_handlers[message_type](data)
                return
            
            # 内置消息处理
            if message_type == "connected":
                await self._handle_connected(data)
            
            elif message_type == "task":
                await self._handle_task(data)
            
            elif message_type == "ping":
                await self._handle_ping(data)
            
            elif message_type == "disconnect_notification":
                await self._handle_disconnect_notification(data)
            
            else:
                logger.info(f"收到未处理的消息类型: {message_type}")
                
        except json.JSONDecodeError as e:
            logger.error(f"消息JSON解析失败: {e}, 原始消息: {message}")
        except Exception as e:
            logger.error(f"处理消息时发生错误: {e}")
    
    async def _handle_connected(self, data: Dict[str, Any]):
        """处理连接成功消息"""
        logger.info("连接成功确认")
        await self.send_message({
            "type": "node_ready",
            "pc_id": self.pc_id,
            "timestamp": int(time.time())
        })
    
    async def _handle_task(self, data: Dict[str, Any]):
        """处理任务消息"""
        # 使用任务处理器处理任务
        result = await self.task_processor.process_task(data)
        
        # 发送处理结果
        await self.send_message(result)
    
    async def _handle_ping(self, data: Dict[str, Any]):
        """处理ping消息"""
        logger.debug("收到ping消息，发送pong响应")
        await self.send_message({"type": "pong"})
    
    async def _handle_disconnect_notification(self, data: Dict[str, Any]):
        """处理断开通知消息"""
        reason = data.get("reason", "未知原因")
        logger.warning(f"收到断开通知: {reason}")
    
    def get_task_processor(self) -> TaskProcessor:
        """获取任务处理器实例"""
        return self.task_processor
