package handler

import (
	"log"
	"time"

	"OctoHub/Ws/internal/connection"
	"OctoHub/Ws/internal/message"

	"github.com/gorilla/websocket"
)

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	connManager    *connection.Manager
	readTimeout    time.Duration
	writeTimeout   time.Duration
	messageHandler message.Handler
}

// NewWebSocketHandler 创建新的WebSocket处理器
func NewWebSocketHandler(connManager *connection.Manager, readTimeout, writeTimeout time.Duration) *WebSocketHandler {
	return &WebSocketHandler{
		connManager:    connManager,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		messageHandler: message.NewDefaultHandler(),
	}
}

// HandleConnection 处理单个连接的消息收发
func (h *WebSocketHandler) HandleConnection(conn *connection.Connection) {
	defer func() {
		h.connManager.RemoveConnection(conn)
		conn.Conn.Close()
	}()

	// 启动发送协程
	go h.writePump(conn)

	// 设置读取超时
	conn.Conn.SetReadDeadline(time.Now().Add(h.readTimeout))

	// 设置pong处理器，当收到客户端的ping时自动回复pong
	conn.Conn.SetPingHandler(func(message string) error {
		conn.UpdateActivity()
		conn.Conn.SetReadDeadline(time.Now().Add(h.readTimeout))

		// 回复pong消息
		conn.Conn.SetWriteDeadline(time.Now().Add(h.writeTimeout))
		return conn.Conn.WriteMessage(websocket.PongMessage, []byte(message))
	})

	// 设置pong处理器，用于处理客户端对服务端ping的响应（如果有的话）
	conn.Conn.SetPongHandler(func(string) error {
		conn.UpdateActivity()
		conn.Conn.SetReadDeadline(time.Now().Add(h.readTimeout))
		return nil
	})

	// 读取消息循环
	for {
		select {
		case <-conn.CloseChan:
			return
		default:
			_, message, err := conn.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket错误: %v", err)
				}
				return
			}

			conn.UpdateActivity()

			// 处理接收到的消息
			h.processMessage(conn, message)
		}
	}
}

// writePump 发送消息协程
func (h *WebSocketHandler) writePump(conn *connection.Connection) {
	for {
		select {
		case message, ok := <-conn.SendChan:
			conn.Conn.SetWriteDeadline(time.Now().Add(h.writeTimeout))
			if !ok {
				conn.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-conn.CloseChan:
			return
		}
	}
}

// processMessage 处理接收到的消息
func (h *WebSocketHandler) processMessage(conn *connection.Connection, rawMessage []byte) {
	// 解析消息
	msg, err := message.FromJSON(rawMessage)
	if err != nil {
		log.Printf("解析消息失败: %v, 原始消息: %s", err, string(rawMessage))
		h.sendErrorMessage(conn, message.ErrorCodeInvalidMessage, "消息格式错误", err.Error())
		return
	}

	// 确定发送者类型
	senderType := "user"
	if conn.Type == connection.PC {
		senderType = "pc"
	}

	// 处理消息
	response := h.messageHandler.HandleMessage(msg, conn.ID, senderType)
	if response != nil {
		h.sendMessage(conn, response)
	}
}

// sendMessage 发送消息
func (h *WebSocketHandler) sendMessage(conn *connection.Connection, msg *message.Message) {
	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}

	select {
	case conn.SendChan <- data:
		// 消息发送成功
	default:
		// 发送通道已满或已关闭
		log.Printf("发送消息失败，连接可能已断开: %s", conn.ID)
		close(conn.SendChan)
	}
}

// sendErrorMessage 发送错误消息
func (h *WebSocketHandler) sendErrorMessage(conn *connection.Connection, code int, msg string, details string) {
	errorMsg := message.SendErrorMessage(code, msg, details)
	h.sendMessage(conn, errorMsg)
}

// SendConnectedMessage 发送连接成功消息
func (h *WebSocketHandler) SendConnectedMessage(conn *connection.Connection) {
	var data message.ConnectedData
	data.SessionID = conn.ID
	data.ServerTime = time.Now().Unix()

	if conn.Type == connection.USER {
		data.UserID = conn.ID
	} else {
		data.PCID = conn.ID
	}

	connectedMsg := message.NewMessage(message.ActionConnected, data)
	h.sendMessage(conn, connectedMsg)
}
