import { useState, useEffect, useRef, useCallback } from 'react';
import { WebSocketManager, WebSocketStatus, WebSocketMessage, WebSocketManagerOptions } from '@/utils/websocket';

export interface UseWebSocketReturn {
  status: WebSocketStatus;
  isConnected: boolean;
  connect: () => Promise<void>;
  disconnect: () => void;
  sendMessage: (message: WebSocketMessage) => boolean;
  lastMessage: WebSocketMessage | null;
  connectionCount: number;
}

export interface UseWebSocketOptions extends Omit<WebSocketManagerOptions, 'onMessage' | 'onStatusChange'> {
  autoConnect?: boolean; // 是否自动连接，默认true
}

/**
 * WebSocket Hook·
 * 提供WebSocket连接管理、消息收发等功能
 */
export function useWebSocket(options: UseWebSocketOptions = {}): UseWebSocketReturn {
  const [status, setStatus] = useState<WebSocketStatus>(WebSocketStatus.DISCONNECTED);
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);
  const [connectionCount, setConnectionCount] = useState(0);
  
  const wsManagerRef = useRef<WebSocketManager | null>(null);
  const { autoConnect = true, ...managerOptions } = options;

  // 初始化WebSocket管理器
  useEffect(() => {
    const wsManager = new WebSocketManager({
      ...managerOptions,
      onMessage: (message: WebSocketMessage) => {
        console.log('Received WebSocket message:', message);
        setLastMessage(message);
      },
      onStatusChange: (newStatus: WebSocketStatus) => {
        console.log('WebSocket status changed:', newStatus);
        setStatus(newStatus);
        
        // 统计连接次数
        if (newStatus === WebSocketStatus.CONNECTED) {
          setConnectionCount(prev => prev + 1);
        }
      },
      onError: (error: Event) => {
        console.error('WebSocket error in hook:', error);
      }
    });

    wsManagerRef.current = wsManager;

    // 自动连接
    if (autoConnect) {
      wsManager.connect().catch(error => {
        console.error('Auto connect failed:', error);
      });
    }

    // 清理函数
    return () => {
      wsManager.disconnect();
      wsManagerRef.current = null;
    };
  }, []); // 空依赖数组，只在组件挂载时初始化

  // 连接WebSocket
  const connect = useCallback(async () => {
    if (wsManagerRef.current) {
      await wsManagerRef.current.connect();
    }
  }, []);

  // 断开WebSocket连接
  const disconnect = useCallback(() => {
    if (wsManagerRef.current) {
      wsManagerRef.current.disconnect();
    }
  }, []);

  // 发送消息
  const sendMessage = useCallback((message: WebSocketMessage): boolean => {
    if (wsManagerRef.current) {
      return wsManagerRef.current.send(message);
    }
    return false;
  }, []);

  // 计算是否已连接
  const isConnected = status === WebSocketStatus.CONNECTED;

  return {
    status,
    isConnected,
    connect,
    disconnect,
    sendMessage,
    lastMessage,
    connectionCount
  };
}

/**
 * 获取WebSocket状态的显示文本
 */
export function getWebSocketStatusText(status: WebSocketStatus): string {
  switch (status) {
    case WebSocketStatus.DISCONNECTED:
      return '未连接';
    case WebSocketStatus.CONNECTING:
      return '连接中';
    case WebSocketStatus.CONNECTED:
      return '已连接';
    case WebSocketStatus.RECONNECTING:
      return '重连中';
    case WebSocketStatus.ERROR:
      return '连接错误';
    default:
      return '未知状态';
  }
}

/**
 * 获取WebSocket状态对应的颜色类名
 */
export function getWebSocketStatusColor(status: WebSocketStatus): string {
  switch (status) {
    case WebSocketStatus.DISCONNECTED:
      return 'text-gray-500';
    case WebSocketStatus.CONNECTING:
    case WebSocketStatus.RECONNECTING:
      return 'text-yellow-500';
    case WebSocketStatus.CONNECTED:
      return 'text-green-500';
    case WebSocketStatus.ERROR:
      return 'text-red-500';
    default:
      return 'text-gray-500';
  }
}
