'use client';

import React, { createContext, useContext, ReactNode } from 'react';
import { UseWebSocketReturn, useWebSocket, UseWebSocketOptions } from '@/hooks/useWebSocket';

interface WebSocketContextType extends UseWebSocketReturn {
  // 可以添加额外的上下文方法
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

interface WebSocketProviderProps {
  children: ReactNode;
  options?: UseWebSocketOptions;
}

/**
 * WebSocket上下文提供者
 * 在应用顶层提供WebSocket连接管理
 */
export function WebSocketProvider({ children, options = {} }: WebSocketProviderProps) {
  const webSocketState = useWebSocket({
    autoConnect: false, // 由使用者决定何时连接
    pingInterval: 10000, // 10秒心跳
    reconnectInterval: 5000, // 5秒重连间隔
    maxReconnectAttempts: 10, // 最大重连10次
    ...options
  });

  return (
    <WebSocketContext.Provider value={webSocketState}>
      {children}
    </WebSocketContext.Provider>
  );
}

/**
 * 使用WebSocket上下文的Hook
 */
export function useWebSocketContext(): WebSocketContextType {
  const context = useContext(WebSocketContext);
  
  if (!context) {
    throw new Error('useWebSocketContext must be used within a WebSocketProvider');
  }
  
  return context;
}

/**
 * 检查是否在WebSocket上下文中
 */
export function useWebSocketContextOptional(): WebSocketContextType | null {
  return useContext(WebSocketContext);
}
