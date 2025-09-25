import { apiGet, ApiResponse, WebSocketConnectionResponse } from './api';

export interface WebSocketMessage {
  action: string;
  data?: any;
  timestamp?: number;
}

export interface WebSocketConnectionInfo {
  wsUrl: string;
}

export interface WebSocketManagerOptions {
  pingInterval?: number; // ping间隔，默认50秒
  reconnectInterval?: number; // 重连间隔，默认5秒
  maxReconnectAttempts?: number; // 最大重连次数，默认10次
  onMessage?: (message: WebSocketMessage) => void;
  onStatusChange?: (status: WebSocketStatus) => void;
  onError?: (error: Event) => void;
}

export enum WebSocketStatus {
  DISCONNECTED = 'disconnected',
  CONNECTING = 'connecting',
  CONNECTED = 'connected',
  RECONNECTING = 'reconnecting',
  ERROR = 'error'
}

export class WebSocketManager {
  private ws: WebSocket | null = null;
  private status: WebSocketStatus = WebSocketStatus.DISCONNECTED;
  private reconnectAttempts = 0;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private pingTimer: NodeJS.Timeout | null = null;
  private pongTimer: NodeJS.Timeout | null = null;
  
  private readonly options: Required<WebSocketManagerOptions>;

  constructor(options: WebSocketManagerOptions = {}) {
    this.options = {
      pingInterval: options.pingInterval ?? 10000, // 50秒 - 比服务端60秒超时短一些
      reconnectInterval: options.reconnectInterval ?? 5000, // 5秒
      maxReconnectAttempts: options.maxReconnectAttempts ?? 10,
      onMessage: options.onMessage ?? (() => {}),
      onStatusChange: options.onStatusChange ?? (() => {}),
      onError: options.onError ?? (() => {})
    };
  }

  /**
   * 获取WebSocket连接地址并建立连接
   */
  async connect(): Promise<void> {
    if (this.status === WebSocketStatus.CONNECTING || this.status === WebSocketStatus.CONNECTED) {
      console.log('WebSocket already connecting or connected');
      return;
    }

    try {
      this.setStatus(WebSocketStatus.CONNECTING);
      
      // 从后端获取WebSocket连接地址
      const response = await apiGet<ApiResponse>('/user/ws');
      if (response.errcode !== 0 || !response.data?.wsUrl) {
        throw new Error(response.errmsg || '获取WebSocket连接地址失败');
      }

      const wsUrl = response.data.wsUrl;
      console.log('Connecting to WebSocket:', wsUrl);
      
      // 建立WebSocket连接
      this.ws = new WebSocket(wsUrl);
      this.setupWebSocketEventHandlers();
      
    } catch (error) {
      console.error('Failed to connect to WebSocket:', error);
      this.setStatus(WebSocketStatus.ERROR);
      this.scheduleReconnect();
    }
  }

  /**
   * 断开WebSocket连接
   */
  disconnect(): void {
    this.clearTimers();
    this.reconnectAttempts = 0;
    
    if (this.ws) {
      this.ws.close(1000, 'User initiated disconnect');
      this.ws = null;
    }
    
    this.setStatus(WebSocketStatus.DISCONNECTED);
  }

  /**
   * 发送消息
   */
  send(message: WebSocketMessage): boolean {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      try {
        this.ws.send(JSON.stringify({
          ...message,
          timestamp: Date.now()
        }));
        return true;
      } catch (error) {
        console.error('Failed to send WebSocket message:', error);
        return false;
      }
    }
    
    console.warn('WebSocket is not connected, cannot send message');
    return false;
  }

  /**
   * 获取当前连接状态
   */
  getStatus(): WebSocketStatus {
    return this.status;
  }

  /**
   * 检查是否已连接
   */
  isConnected(): boolean {
    return this.status === WebSocketStatus.CONNECTED;
  }

  private setupWebSocketEventHandlers(): void {
    if (!this.ws) return;

    this.ws.onopen = () => {
      console.log('WebSocket connected successfully');
      this.setStatus(WebSocketStatus.CONNECTED);
      this.reconnectAttempts = 0;
      this.startPingTimer();
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as WebSocketMessage;
        
        // 处理pong消息
        if (message.action === 'pong') {
          console.log('Received pong response:', message);
          this.handlePong();
          return;
        }
        
        // 调用用户定义的消息处理器
        this.options.onMessage(message);
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    this.ws.onclose = (event) => {
      console.log('WebSocket connection closed:', event.code, event.reason);
      this.clearTimers();
      
      if (event.code !== 1000) { // 不是正常关闭
        this.setStatus(WebSocketStatus.DISCONNECTED);
        this.scheduleReconnect();
      } else {
        this.setStatus(WebSocketStatus.DISCONNECTED);
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.setStatus(WebSocketStatus.ERROR);
      this.options.onError(error);
    };
  }

  private setStatus(newStatus: WebSocketStatus): void {
    if (this.status !== newStatus) {
      this.status = newStatus;
      this.options.onStatusChange(newStatus);
    }
  }

  private startPingTimer(): void {
    this.clearPingTimer();
    
    this.pingTimer = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        // 发送JSON格式的ping消息，等待服务端返回pong消息
        console.log('Sending ping message');
        this.send({ action: 'ping' });
        
        // 设置pong超时检测
        this.pongTimer = setTimeout(() => {
          console.warn('Pong timeout after 10 seconds, closing WebSocket connection');
          this.ws?.close(1000, 'Pong timeout');
        }, 10000); // 10秒超时
      }
    }, this.options.pingInterval);
  }

  private handlePong(): void {
    console.log('Handling pong response, clearing timeout');
    if (this.pongTimer) {
      clearTimeout(this.pongTimer);
      this.pongTimer = null;
    }
  }

  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.options.maxReconnectAttempts) {
      console.error(`Max reconnect attempts (${this.options.maxReconnectAttempts}) reached, giving up`);
      this.setStatus(WebSocketStatus.ERROR);
      return;
    }

    this.setStatus(WebSocketStatus.RECONNECTING);
    this.reconnectAttempts++;
    
    console.log(`Scheduling reconnect attempt ${this.reconnectAttempts}/${this.options.maxReconnectAttempts} in ${this.options.reconnectInterval}ms`);
    
    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, this.options.reconnectInterval);
  }

  private clearTimers(): void {
    this.clearPingTimer();
    this.clearReconnectTimer();
  }

  private clearPingTimer(): void {
    if (this.pingTimer) {
      clearInterval(this.pingTimer);
      this.pingTimer = null;
    }
    
    if (this.pongTimer) {
      clearTimeout(this.pongTimer);
      this.pongTimer = null;
    }
  }

  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }
}
