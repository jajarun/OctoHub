"""
OctoHub Node Client - 主入口文件
"""

import asyncio

from octohub.client import NodeWebSocketClient
from octohub.utils.logger import setup_logger
from config import config

# 设置日志
logger = setup_logger(__name__, config.LOG_LEVEL)


async def setup_custom_handlers(client: NodeWebSocketClient):
    """设置自定义处理器示例"""
    
    # 自定义任务处理器示例
    async def handle_custom_task(task_data):
        """自定义任务处理器"""
        task_id = task_data.get("task_id", "unknown")
        logger.info(f"处理自定义任务: {task_id}")
        
        # 在这里添加你的自定义任务逻辑
        await asyncio.sleep(0.5)  # 模拟处理时间
        
        return f"自定义任务 {task_id} 处理完成"
    
    # 自定义消息处理器示例
    async def handle_custom_message(message_data):
        """自定义消息处理器"""
        logger.info(f"收到自定义消息: {message_data}")
    
    # 注册处理器
    client.register_task_handler("custom_task", handle_custom_task)
    client.register_message_handler("custom_message", handle_custom_message)


async def main():
    """主函数"""
    # 检查是否使用默认签名密钥
    if config.is_default_signature_key():
        logger.warning("警告: 使用默认签名密钥，生产环境请设置OCTOHUB_SIGNATURE_KEY环境变量")
    
    # 创建WebSocket客户端实例
    client = NodeWebSocketClient()
    
    # 设置自定义处理器（可选）
    await setup_custom_handlers(client)
    
    try:
        # 启动客户端
        logger.info("启动 OctoHub Node 客户端...")
        await client.start()
    except KeyboardInterrupt:
        logger.info("收到中断信号，正在关闭客户端...")
        await client.stop()
    except Exception as e:
        logger.error(f"客户端运行时发生错误: {e}")
        await client.stop()


if __name__ == "__main__":
    # 运行异步主函数
    asyncio.run(main())
