#!/bin/bash

# OctoHub 多服务停止脚本
# 优雅地停止所有服务

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

# 停止服务函数
stop_service() {
    local service_name=$1
    local pid_file=$2
    local service_port=$3
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            info "正在停止 $service_name (PID: $pid)..."
            
            # 首先尝试优雅停止
            kill -TERM "$pid" 2>/dev/null || true
            
            # 等待进程结束
            local count=0
            while kill -0 "$pid" 2>/dev/null && [ $count -lt 10 ]; do
                sleep 1
                count=$((count + 1))
            done
            
            # 如果还没结束，强制杀死
            if kill -0 "$pid" 2>/dev/null; then
                warn "$service_name 未能优雅停止，强制终止..."
                kill -KILL "$pid" 2>/dev/null || true
                sleep 1
            fi
            
            if ! kill -0 "$pid" 2>/dev/null; then
                log "$service_name 已停止"
            else
                error "无法停止 $service_name"
            fi
        else
            warn "$service_name 进程不存在 (PID: $pid)"
        fi
        rm -f "$pid_file"
    else
        info "$service_name 没有运行 (PID文件不存在)"
    fi
    
    # 检查端口是否还在使用
    if [ -n "$service_port" ]; then
        local port_pid=$(lsof -ti:$service_port 2>/dev/null || true)
        if [ -n "$port_pid" ]; then
            warn "端口 $service_port 仍被占用 (PID: $port_pid)，尝试释放..."
            kill -TERM "$port_pid" 2>/dev/null || true
            sleep 2
            port_pid=$(lsof -ti:$service_port 2>/dev/null || true)
            if [ -n "$port_pid" ]; then
                kill -KILL "$port_pid" 2>/dev/null || true
            fi
        fi
    fi
}

# 停止所有Node.js进程（前端相关）
stop_nodejs_processes() {
    info "检查并停止相关的 Node.js 进程..."
    
    # 查找运行在3000端口的进程
    local node_pids=$(lsof -ti:3000 2>/dev/null || true)
    if [ -n "$node_pids" ]; then
        for pid in $node_pids; do
            local cmd=$(ps -p "$pid" -o comm= 2>/dev/null || true)
            if [[ "$cmd" == *"node"* ]] || [[ "$cmd" == *"next"* ]]; then
                info "停止 Node.js 进程 (PID: $pid)..."
                kill -TERM "$pid" 2>/dev/null || true
                sleep 1
                if kill -0 "$pid" 2>/dev/null; then
                    kill -KILL "$pid" 2>/dev/null || true
                fi
            fi
        done
    fi
}

# 停止Java进程
stop_java_processes() {
    info "检查并停止相关的 Java 进程..."
    
    # 查找运行在8080端口的进程
    local java_pids=$(lsof -ti:8080 2>/dev/null || true)
    if [ -n "$java_pids" ]; then
        for pid in $java_pids; do
            local cmd=$(ps -p "$pid" -o comm= 2>/dev/null || true)
            if [[ "$cmd" == *"java"* ]]; then
                info "停止 Java 进程 (PID: $pid)..."
                kill -TERM "$pid" 2>/dev/null || true
                sleep 2
                if kill -0 "$pid" 2>/dev/null; then
                    kill -KILL "$pid" 2>/dev/null || true
                fi
            fi
        done
    fi
    
    # 也检查Maven进程
    local maven_pids=$(pgrep -f "maven" 2>/dev/null || true)
    if [ -n "$maven_pids" ]; then
        for pid in $maven_pids; do
            info "停止 Maven 进程 (PID: $pid)..."
            kill -TERM "$pid" 2>/dev/null || true
            sleep 1
            if kill -0 "$pid" 2>/dev/null; then
                kill -KILL "$pid" 2>/dev/null || true
            fi
        done
    fi
}

# 停止Go进程
stop_go_processes() {
    info "检查并停止相关的 Go 进程..."
    
    # 查找运行在8000端口的进程
    local go_pids=$(lsof -ti:8000 2>/dev/null || true)
    if [ -n "$go_pids" ]; then
        for pid in $go_pids; do
            local cmd=$(ps -p "$pid" -o args= 2>/dev/null || true)
            if [[ "$cmd" == *"main.go"* ]] || [[ "$cmd" == *"websocket-server"* ]]; then
                info "停止 Go 进程 (PID: $pid)..."
                kill -TERM "$pid" 2>/dev/null || true
                sleep 1
                if kill -0 "$pid" 2>/dev/null; then
                    kill -KILL "$pid" 2>/dev/null || true
                fi
            fi
        done
    fi
}

# 清理日志和PID文件
cleanup_files() {
    info "清理临时文件..."
    
    # 删除PID文件
    rm -f logs/frontend.pid logs/backend.pid logs/websocket.pid
    
    # 可选：清理日志文件（注释掉以保留日志）
    # rm -f logs/frontend.log logs/backend.log logs/websocket.log
    
    log "临时文件已清理"
}

# 显示停止状态
show_stop_status() {
    log "所有服务已停止！"
    echo ""
    echo -e "${GREEN}=== 端口状态检查 ===${NC}"
    
    # 检查端口是否已释放
    local ports=(3000 8080 8000)
    local port_names=("前端" "Java后端" "WebSocket")
    
    for i in "${!ports[@]}"; do
        local port=${ports[$i]}
        local name=${port_names[$i]}
        if lsof -ti:$port >/dev/null 2>&1; then
            error "$name 端口 $port 仍被占用"
        else
            log "$name 端口 $port 已释放"
        fi
    done
    
    echo ""
    if [ -d "logs" ] && [ "$(ls -A logs 2>/dev/null)" ]; then
        echo -e "${BLUE}日志文件保留在 logs/ 目录中${NC}"
    fi
}

# 主函数
main() {
    log "开始停止 OctoHub 服务..."
    
    # 检查logs目录是否存在
    if [ ! -d "logs" ]; then
        info "logs 目录不存在，可能没有服务在运行"
    fi
    
    # 按PID文件停止服务
    stop_service "前端服务" "logs/frontend.pid" "3000"
    stop_service "Java后端服务" "logs/backend.pid" "8080"
    stop_service "Go WebSocket服务" "logs/websocket.pid" "8000"
    
    # 额外检查并停止可能遗留的进程
    stop_nodejs_processes
    stop_java_processes
    stop_go_processes
    
    # 清理文件
    cleanup_files
    
    # 显示状态
    show_stop_status
}

# 如果脚本被直接运行
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
