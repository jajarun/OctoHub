"""
签名认证模块
"""

import base64
import hashlib
import hmac
import time
from typing import Dict, Any
import aiohttp

from ..utils.logger import setup_logger

logger = setup_logger(__name__)


class SignatureAuth:
    """签名认证处理器"""
    
    def __init__(self, signature_key: str):
        self.signature_key = signature_key
    
    def generate_signature(self, method: str, uri: str, params: Dict[str, str], timestamp: str, nonce: str) -> str:
        """生成API签名 (Base64编码)"""
        # 构建签名字符串: METHOD&URI&PARAMS&TIMESTAMP&NONCE
        # 参数按key排序后拼接，格式与Java服务端一致
        if params:
            params_str = "&".join([f"{k}={v}" for k, v in sorted(params.items())])
        else:
            params_str = ""
        
        sign_string = f"{method.upper()}&{uri}&{params_str}&{timestamp}&{nonce}"
        
        logger.debug(f"签名字符串: {sign_string}")
        
        signature = hmac.new(
            self.signature_key.encode(),
            sign_string.encode(),
            hashlib.sha256
        )
        return base64.b64encode(signature.digest()).decode()
    
    def create_auth_headers(self, method: str, uri: str, params: Dict[str, str]) -> Dict[str, str]:
        """创建API认证请求头"""
        timestamp = str(int(time.time()))
        nonce = str(int(time.time() * 1000) + 999)
        signature = self.generate_signature(method, uri, params, timestamp, nonce)
        
        return {
            "X-Signature": signature,
            "X-Timestamp": timestamp,
            "X-Nonce": nonce,
            "Content-Type": "application/json"
        }
    
    async def get_websocket_url(self, 
                               pc_id: str, 
                               server_host: str, 
                               server_port: int) -> str:
        """从服务器获取WebSocket连接地址"""
        # 构造API请求参数
        method = "GET"
        uri = "/node/ws"
        params = {"pc_id": pc_id}
        
        # 构造API请求URL
        api_url = f"http://{server_host}:{server_port}{uri}"
        
        # 构造签名头部
        headers = self.create_auth_headers(method, uri, params)
        
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(api_url, params=params, headers=headers) as response:
                    if response.status == 200:
                        data = await response.json()
                        if data.get("errcode") == 0 and data.get("data", {}).get("wsUrl"):
                            ws_url = data["data"]["wsUrl"]
                            logger.info(f"获取到WebSocket连接地址: {ws_url}")
                            return ws_url
                        else:
                            raise Exception(f"API返回错误: {data.get('errmsg', '未知错误')}")
                    else:
                        error_text = await response.text()
                        raise Exception(f"HTTP错误 {response.status}: {error_text}")
        except Exception as e:
            logger.error(f"获取WebSocket连接地址失败: {e}")
            raise
