#!/bin/bash

# OctoHub Node Client 启动脚本

echo "=== OctoHub Node Client 启动脚本 ==="

# 检查Python版本
python_version=$(python3 --version 2>/dev/null)
if [ $? -ne 0 ]; then
    echo "错误: 未找到 Python 3"
    echo "请安装 Python 3.12 或更高版本"
    exit 1
fi

echo "检测到 Python 版本: $python_version"

# 检查是否在虚拟环境中
if [[ "$VIRTUAL_ENV" != "" ]]; then
    echo "当前在虚拟环境中: $VIRTUAL_ENV"
else
    echo "建议创建虚拟环境："
    echo "python3 -m venv venv"
    echo "source venv/bin/activate"
fi

# 安装依赖
echo "安装项目依赖..."
pip install -e . || {
    echo "安装依赖失败，尝试手动安装..."
    pip install websockets>=12.0 aiohttp>=3.9.0 || {
        echo "错误: 依赖安装失败"
        exit 1
    }
}

echo "依赖安装完成"

# 检查环境变量
echo ""
echo "=== 环境变量检查 ==="
if [ -z "$OCTOHUB_SIGNATURE_KEY" ]; then
    echo "警告: 未设置 OCTOHUB_SIGNATURE_KEY 环境变量"
    echo "将使用默认密钥（仅适用于开发环境）"
else
    echo "✓ 签名密钥已设置"
fi

echo "服务器地址: ${OCTOHUB_SERVER_HOST:-localhost}"
echo "服务器端口: ${OCTOHUB_SERVER_PORT:-8080}"
echo "WebSocket端口: ${OCTOHUB_WS_PORT:-8000}"
echo "日志级别: ${OCTOHUB_LOG_LEVEL:-INFO}"

echo ""
echo "=== 启动客户端 ==="
echo "按 Ctrl+C 停止客户端"
echo ""

# 启动客户端
python3 main.py
