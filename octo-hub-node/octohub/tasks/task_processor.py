"""
任务处理器模块
"""

import asyncio
import time
from typing import Dict, Any, Callable, Awaitable

from ..utils.logger import setup_logger

logger = setup_logger(__name__)


class TaskProcessor:
    """任务处理器"""
    
    def __init__(self):
        self.task_handlers: Dict[str, Callable[[Dict[str, Any]], Awaitable[Any]]] = {}
        self._register_default_handlers()
    
    def _register_default_handlers(self):
        """注册默认的任务处理器"""
        self.register_handler("default", self._handle_default_task)
        self.register_handler("echo", self._handle_echo_task)
        self.register_handler("sleep", self._handle_sleep_task)
    
    def register_handler(self, 
                        task_type: str, 
                        handler: Callable[[Dict[str, Any]], Awaitable[Any]]):
        """注册任务处理器"""
        self.task_handlers[task_type] = handler
        logger.info(f"注册任务处理器: {task_type}")
    
    async def process_task(self, task_data: Dict[str, Any]) -> Dict[str, Any]:
        """处理任务"""
        task_id = task_data.get("task_id", "unknown")
        task_type = task_data.get("task_type", "default")
        
        logger.info(f"开始处理任务 - ID: {task_id}, 类型: {task_type}")
        
        try:
            # 获取对应的任务处理器
            handler = self.task_handlers.get(task_type, self.task_handlers["default"])
            
            # 执行任务处理
            result = await handler(task_data)
            
            # 返回成功结果
            response = {
                "type": "task_result",
                "task_id": task_id,
                "status": "completed",
                "result": result,
                "timestamp": int(time.time())
            }
            
            logger.info(f"任务 {task_id} 处理完成")
            return response
            
        except Exception as e:
            logger.error(f"任务 {task_id} 处理失败: {e}")
            
            # 返回失败结果
            return {
                "type": "task_result",
                "task_id": task_id,
                "status": "failed",
                "error": str(e),
                "timestamp": int(time.time())
            }
    
    async def _handle_default_task(self, task_data: Dict[str, Any]) -> str:
        """默认任务处理器"""
        task_id = task_data.get("task_id", "unknown")
        await asyncio.sleep(0.1)  # 模拟处理时间
        return f"任务 {task_id} 执行完成（默认处理器）"
    
    async def _handle_echo_task(self, task_data: Dict[str, Any]) -> Dict[str, Any]:
        """回声任务处理器"""
        return {
            "echo": task_data.get("data", {}),
            "processed_at": int(time.time())
        }
    
    async def _handle_sleep_task(self, task_data: Dict[str, Any]) -> str:
        """睡眠任务处理器"""
        sleep_time = task_data.get("data", {}).get("sleep_time", 1)
        sleep_time = max(0.1, min(10, float(sleep_time)))  # 限制在0.1-10秒之间
        
        logger.info(f"执行睡眠任务，时长: {sleep_time}秒")
        await asyncio.sleep(sleep_time)
        
        return f"睡眠任务完成，睡眠时长: {sleep_time}秒"
