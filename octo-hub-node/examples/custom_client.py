"""
自定义客户端使用示例
"""

import asyncio
import time
from typing import Dict, Any

from octohub.client import NodeWebSocketClient
from octohub.utils.logger import setup_logger

logger = setup_logger(__name__)


class CustomNodeClient:
    """自定义Node客户端示例"""
    
    def __init__(self):
        self.client = NodeWebSocketClient()
        self._setup_handlers()
    
    def _setup_handlers(self):
        """设置自定义处理器"""
        
        # 注册各种任务处理器
        self.client.register_task_handler("file_operation", self._handle_file_operation)
        self.client.register_task_handler("system_info", self._handle_system_info)
        self.client.register_task_handler("calculation", self._handle_calculation)
        
        # 注册自定义消息处理器
        self.client.register_message_handler("heartbeat", self._handle_heartbeat)
        self.client.register_message_handler("config_update", self._handle_config_update)
    
    async def _handle_file_operation(self, task_data: Dict[str, Any]) -> Dict[str, Any]:
        """处理文件操作任务"""
        operation = task_data.get("data", {}).get("operation", "unknown")
        file_path = task_data.get("data", {}).get("file_path", "")
        
        logger.info(f"执行文件操作: {operation} -> {file_path}")
        
        # 模拟文件操作
        await asyncio.sleep(0.2)
        
        return {
            "operation": operation,
            "file_path": file_path,
            "status": "success",
            "message": f"文件操作 {operation} 完成"
        }
    
    async def _handle_system_info(self, task_data: Dict[str, Any]) -> Dict[str, Any]:
        """处理系统信息查询任务"""
        import platform
        import psutil
        
        logger.info("收集系统信息")
        
        return {
            "system": platform.system(),
            "platform": platform.platform(),
            "processor": platform.processor(),
            "memory": {
                "total": psutil.virtual_memory().total,
                "available": psutil.virtual_memory().available,
                "percent": psutil.virtual_memory().percent
            },
            "cpu_percent": psutil.cpu_percent(interval=1),
            "timestamp": int(time.time())
        }
    
    async def _handle_calculation(self, task_data: Dict[str, Any]) -> Dict[str, Any]:
        """处理计算任务"""
        data = task_data.get("data", {})
        operation = data.get("operation", "add")
        numbers = data.get("numbers", [])
        
        logger.info(f"执行计算: {operation} -> {numbers}")
        
        if operation == "add":
            result = sum(numbers)
        elif operation == "multiply":
            result = 1
            for num in numbers:
                result *= num
        elif operation == "average":
            result = sum(numbers) / len(numbers) if numbers else 0
        else:
            raise ValueError(f"不支持的操作: {operation}")
        
        return {
            "operation": operation,
            "numbers": numbers,
            "result": result
        }
    
    async def _handle_heartbeat(self, message_data: Dict[str, Any]):
        """处理心跳消息"""
        logger.info("收到心跳消息")
        # 这里可以更新客户端状态或执行其他操作
    
    async def _handle_config_update(self, message_data: Dict[str, Any]):
        """处理配置更新消息"""
        config_data = message_data.get("config", {})
        logger.info(f"收到配置更新: {config_data}")
        # 这里可以更新客户端配置
    
    async def start(self):
        """启动自定义客户端"""
        logger.info("启动自定义Node客户端")
        await self.client.start()
    
    async def stop(self):
        """停止自定义客户端"""
        logger.info("停止自定义Node客户端")
        await self.client.stop()
    
    def get_pc_id(self) -> str:
        """获取PC ID"""
        return self.client.get_pc_id()


async def main():
    """主函数"""
    client = CustomNodeClient()
    
    logger.info(f"自定义客户端PC ID: {client.get_pc_id()}")
    
    try:
        await client.start()
    except KeyboardInterrupt:
        logger.info("收到中断信号，正在关闭客户端...")
        await client.stop()
    except Exception as e:
        logger.error(f"客户端运行时发生错误: {e}")
        await client.stop()


if __name__ == "__main__":
    asyncio.run(main())
