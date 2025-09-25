'use client';

import React, { useState } from 'react';
import { useWebSocket, getWebSocketStatusText, getWebSocketStatusColor } from '@/hooks/useWebSocket';
import { WebSocketMessage } from '@/utils/websocket';

/**
 * WebSocket功能演示组件
 * 展示如何使用WebSocket连接、发送消息等功能
 */
export default function WebSocketDemo() {
  const [messageInput, setMessageInput] = useState('');
  const [messageHistory, setMessageHistory] = useState<WebSocketMessage[]>([]);

  const {
    status,
    isConnected,
    connect,
    disconnect,
    sendMessage,
    lastMessage,
    connectionCount
  } = useWebSocket({
    autoConnect: false
  });

  // 处理接收到的消息
  React.useEffect(() => {
    if (lastMessage && lastMessage.action === 'chat') {
      setMessageHistory(prev => [...prev.slice(-9), lastMessage]); // 保留最近10条消息
    }
  }, [lastMessage]);

  const handleSendMessage = () => {
    if (!messageInput.trim()) return;

    const message: WebSocketMessage = {
      action: 'chat',
      data: {
        text: messageInput.trim(),
        sender: 'user'
      }
    };

    const sent = sendMessage(message);
    if (sent) {
      setMessageInput('');
      // 添加到消息历史（发送的消息）
      setMessageHistory(prev => [...prev.slice(-9), {
        ...message,
        timestamp: Date.now()
      }]);
    }
  };

  const handleSendPing = () => {
    sendMessage({ action: 'ping' });
  };

  const handleClearHistory = () => {
    setMessageHistory([]);
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="mb-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-2">WebSocket 连接演示</h2>
        <p className="text-sm text-gray-600">
          这个组件演示了如何使用WebSocket进行实时通信，包括连接管理、消息收发和状态监控。
        </p>
      </div>

      {/* 连接状态 */}
      <div className="mb-6 p-4 bg-gray-50 rounded-lg">
        <div className="flex items-center justify-between mb-3">
          <h3 className="font-medium text-gray-900">连接状态</h3>
          <div className="flex items-center space-x-2">
            <span className={`text-sm font-medium ${getWebSocketStatusColor(status)}`}>
              {getWebSocketStatusText(status)}
            </span>
            {connectionCount > 0 && (
              <span className="text-xs text-gray-500">
                (连接次数: {connectionCount})
              </span>
            )}
          </div>
        </div>

        <div className="flex space-x-3">
          <button
            onClick={connect}
            disabled={isConnected}
            className="px-4 py-2 bg-blue-600 text-white text-sm rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            连接
          </button>
          <button
            onClick={disconnect}
            disabled={!isConnected}
            className="px-4 py-2 bg-red-600 text-white text-sm rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            断开
          </button>
          <button
            onClick={handleSendPing}
            disabled={!isConnected}
            className="px-4 py-2 bg-green-600 text-white text-sm rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            发送Ping
          </button>
        </div>
      </div>

      {/* 消息发送 */}
      <div className="mb-6">
        <h3 className="font-medium text-gray-900 mb-3">发送消息</h3>
        <div className="flex space-x-3">
          <input
            type="text"
            value={messageInput}
            onChange={(e) => setMessageInput(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
            placeholder="输入消息内容..."
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            disabled={!isConnected}
          />
          <button
            onClick={handleSendMessage}
            disabled={!isConnected || !messageInput.trim()}
            className="px-4 py-2 bg-blue-600 text-white text-sm rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            发送
          </button>
        </div>
      </div>

      {/* 消息历史 */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <h3 className="font-medium text-gray-900">消息历史</h3>
          <button
            onClick={handleClearHistory}
            disabled={messageHistory.length === 0}
            className="px-3 py-1 text-sm text-gray-600 hover:text-gray-800 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            清空
          </button>
        </div>

        <div className="bg-gray-50 rounded-lg p-4 h-64 overflow-y-auto">
          {messageHistory.length === 0 ? (
            <div className="text-center text-gray-500 text-sm py-8">
              暂无消息
            </div>
          ) : (
            <div className="space-y-2">
              {messageHistory.map((message, index) => (
                <div
                  key={index}
                  className={`p-2 rounded text-sm ${
                    message.data?.sender === 'user'
                      ? 'bg-blue-100 text-blue-900 ml-8'
                      : 'bg-white text-gray-900 mr-8'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-medium text-xs text-gray-500">
                      {message.action}
                    </span>
                    {message.timestamp && (
                      <span className="text-xs text-gray-400">
                        {new Date(message.timestamp).toLocaleTimeString()}
                      </span>
                    )}
                  </div>
                  <div className="mt-1">
                    {typeof message.data === 'object' 
                      ? JSON.stringify(message.data, null, 2)
                      : message.data || '(无数据)'
                    }
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* 最新消息提示 */}
      {lastMessage && (
        <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
          <div className="text-sm text-yellow-800">
            <strong>最新消息:</strong> {lastMessage.action}
            {lastMessage.timestamp && (
              <span className="ml-2 text-yellow-600">
                ({new Date(lastMessage.timestamp).toLocaleTimeString()})
              </span>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
