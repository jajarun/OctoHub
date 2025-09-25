#!/bin/bash

# OctoHub 多服务启动脚本
# 同时启动前端、Java后端和Go WebSocket服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

# 检查必要的命令是否存在
check_dependencies() {
    log "检查依赖..."
    
    if ! command -v node &> /dev/null; then
        error "Node.js 未安装，请先安装 Node.js"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        error "npm 未安装，请先安装 npm"
        exit 1
    fi
    
    if ! command -v mvn &> /dev/null; then
        error "Maven 未安装，请先安装 Maven"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        error "Go 未安装，请先安装 Go"
        exit 1
    fi
    
    log "所有依赖检查通过"
}

# 创建日志目录
create_log_dir() {
    mkdir -p logs
    log "日志目录已创建: ./logs"
}

# 启动前端服务
start_frontend() {
    log "启动前端服务 (Next.js)..."
    cd octo-hub-front
    
    # 检查是否需要安装依赖
    if [ ! -d "node_modules" ]; then
        info "安装前端依赖..."
        npm install
    fi
    
    # 启动前端服务
    npm run dev > ../logs/frontend.log 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > ../logs/frontend.pid
    
    cd ..
    log "前端服务已启动 (PID: $FRONTEND_PID, 端口: 3000)"
}

# 启动Java后端服务
start_backend() {
    log "启动Java后端服务 (Spring Boot)..."
    cd octo-hub-server
    
    # 编译并启动Spring Boot应用
    mvn spring-boot:run > ../logs/backend.log 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > ../logs/backend.pid
    
    cd ..
    log "Java后端服务已启动 (PID: $BACKEND_PID, 端口: 8080)"
}

# 启动Go WebSocket服务
start_websocket() {
    log "启动Go WebSocket服务..."
    cd octo-hub-ws
    
    # 检查是否需要下载依赖
    if [ ! -f "go.sum" ]; then
        info "下载Go依赖..."
        go mod download
    fi
    
    # 启动Go服务
    go run main.go > ../logs/websocket.log 2>&1 &
    WEBSOCKET_PID=$!
    echo $WEBSOCKET_PID > ../logs/websocket.pid
    
    cd ..
    log "Go WebSocket服务已启动 (PID: $WEBSOCKET_PID, 端口: 8000)"
}

# 等待服务启动
wait_for_services() {
    log "等待所有服务启动完成..."
    
    # 等待前端服务
    info "等待前端服务启动 (http://localhost:3000)..."
    for i in {1..30}; do
        if curl -s http://localhost:3000 > /dev/null 2>&1; then
            log "前端服务已就绪"
            break
        fi
        if [ $i -eq 30 ]; then
            warn "前端服务启动超时，请检查日志: logs/frontend.log"
        fi
        sleep 2
    done
    
    # 等待后端服务
    info "等待Java后端服务启动 (http://localhost:8080)..."
    for i in {1..60}; do
        if curl -s http://localhost:8080/actuator/health > /dev/null 2>&1 || curl -s http://localhost:8080 > /dev/null 2>&1; then
            log "Java后端服务已就绪"
            break
        fi
        if [ $i -eq 60 ]; then
            warn "Java后端服务启动超时，请检查日志: logs/backend.log"
        fi
        sleep 2
    done
    
    # 等待WebSocket服务
    info "等待Go WebSocket服务启动 (ws://localhost:8000)..."
    for i in {1..30}; do
        if nc -z localhost 8000 > /dev/null 2>&1; then
            log "Go WebSocket服务已就绪"
            break
        fi
        if [ $i -eq 30 ]; then
            warn "Go WebSocket服务启动超时，请检查日志: logs/websocket.log"
        fi
        sleep 2
    done
}

# 显示服务状态
show_status() {
    log "所有服务启动完成！"
    echo ""
    echo -e "${GREEN}=== 服务状态 ===${NC}"
    echo -e "${BLUE}前端服务:${NC}     http://localhost:3000"
    echo -e "${BLUE}Java后端:${NC}     http://localhost:8080"
    echo -e "${BLUE}WebSocket:${NC}    ws://localhost:8000"
    echo ""
    echo -e "${GREEN}=== 日志文件 ===${NC}"
    echo -e "${BLUE}前端日志:${NC}     logs/frontend.log"
    echo -e "${BLUE}后端日志:${NC}     logs/backend.log"
    echo -e "${BLUE}WebSocket日志:${NC} logs/websocket.log"
    echo ""
    echo -e "${YELLOW}使用 './stop.sh' 停止所有服务${NC}"
    echo -e "${YELLOW}使用 'tail -f logs/*.log' 查看实时日志${NC}"
}

# 清理函数
cleanup() {
    error "收到中断信号，正在停止服务..."
    
    if [ -f "logs/frontend.pid" ]; then
        kill $(cat logs/frontend.pid) 2>/dev/null || true
        rm -f logs/frontend.pid
    fi
    
    if [ -f "logs/backend.pid" ]; then
        kill $(cat logs/backend.pid) 2>/dev/null || true
        rm -f logs/backend.pid
    fi
    
    if [ -f "logs/websocket.pid" ]; then
        kill $(cat logs/websocket.pid) 2>/dev/null || true
        rm -f logs/websocket.pid
    fi
    
    exit 1
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 主函数
main() {
    log "开始启动 OctoHub 服务..."
    
    check_dependencies
    create_log_dir
    
    # 并行启动所有服务
    start_frontend &
    start_backend &
    start_websocket &
    
    # 等待所有后台任务完成
    wait
    
    # 等待服务就绪
    wait_for_services
    
    # 显示状态
    show_status
    
    # 保持脚本运行
    log "服务运行中... 按 Ctrl+C 停止所有服务"
    while true; do
        sleep 1
    done
}

# 执行主函数
main "$@"
