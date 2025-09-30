"""
OctoHub Node 配置文件
"""

import os
from pathlib import Path
from typing import Dict, Any

try:
    from dotenv import load_dotenv
    # 加载 .env 文件 - 从当前工作目录查找
    env_path = Path.cwd() / '.env'
    if env_path.exists():
        load_dotenv(env_path)
        print(f"已加载环境配置文件: {env_path}")
    else:
        # 尝试从配置文件所在目录查找
        env_path = Path(__file__).parent / '.env'
        if env_path.exists():
            load_dotenv(env_path)
            print(f"已加载环境配置文件: {env_path}")
        else:
            print("未找到 .env 文件，使用系统环境变量或默认配置")
except ImportError:
    print("python-dotenv 未安装，跳过 .env 文件加载")


class Config:
    """配置类"""
    
    def __init__(self):
        # 服务器配置
        self.SERVER_HOST = os.getenv("OCTOHUB_SERVER_HOST", "localhost")
        self.SERVER_PORT = int(os.getenv("OCTOHUB_SERVER_PORT", "8080"))
        self.WS_PORT = int(os.getenv("OCTOHUB_WS_PORT", "8000"))
        
        # 签名认证配置
        self.SIGNATURE_KEY = os.getenv("OCTOHUB_SIGNATURE_KEY", "mySignatureKey123456789012345")
        
        # WebSocket配置
        self.RECONNECT_INTERVAL = int(os.getenv("OCTOHUB_RECONNECT_INTERVAL", "5"))
        self.MAX_RECONNECT_ATTEMPTS = int(os.getenv("OCTOHUB_MAX_RECONNECT_ATTEMPTS", "10"))
        self.PING_INTERVAL = int(os.getenv("OCTOHUB_PING_INTERVAL", "30"))
        self.PING_TIMEOUT = int(os.getenv("OCTOHUB_PING_TIMEOUT", "10"))
        
        # 日志配置
        self.LOG_LEVEL = os.getenv("OCTOHUB_LOG_LEVEL", "INFO")
        
    def to_dict(self) -> Dict[str, Any]:
        """转换为字典"""
        return {
            "server_host": self.SERVER_HOST,
            "server_port": self.SERVER_PORT,
            "ws_port": self.WS_PORT,
            "signature_key": self.SIGNATURE_KEY,
            "reconnect_interval": self.RECONNECT_INTERVAL,
            "max_reconnect_attempts": self.MAX_RECONNECT_ATTEMPTS,
            "ping_interval": self.PING_INTERVAL,
            "ping_timeout": self.PING_TIMEOUT,
            "log_level": self.LOG_LEVEL
        }
    
    def is_default_signature_key(self) -> bool:
        """检查是否使用默认签名密钥"""
        return self.SIGNATURE_KEY == "mySignatureKey123456789012345"


# 全局配置实例
config = Config()
