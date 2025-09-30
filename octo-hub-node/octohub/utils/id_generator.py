"""
PC ID生成器
"""

import hashlib
import platform
import time
import uuid


def generate_pc_id() -> str:
    """生成唯一的PC ID"""
    # 获取机器信息
    hostname = platform.node()
    mac_address = hex(uuid.getnode())[2:]  # 去掉'0x'前缀
    
    # 使用主机名和MAC地址生成唯一ID
    unique_string = f"{hostname}_{mac_address}_{int(time.time())}"
    pc_id = hashlib.md5(unique_string.encode()).hexdigest()[:16]
    
    return f"node_{pc_id}"
