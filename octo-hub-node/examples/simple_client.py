"""
简单客户端使用示例
"""

import asyncio
from octohub.client import NodeWebSocketClient
from octohub.utils.logger import setup_logger

logger = setup_logger(__name__)


async def main():
    """简单使用示例"""
    
    # 创建客户端
    client = NodeWebSocketClient()
    
    # 添加一个简单的自定义任务处理器
    async def handle_greeting_task(task_data):
        name = task_data.get("data", {}).get("name", "World")
        return f"Hello, {name}!"
    
    # 注册任务处理器
    client.register_task_handler("greeting", handle_greeting_task)
    
    logger.info(f"客户端PC ID: {client.get_pc_id()}")
    
    try:
        # 启动客户端
        await client.start()
    except KeyboardInterrupt:
        logger.info("客户端已停止")
        await client.stop()


if __name__ == "__main__":
    asyncio.run(main())
