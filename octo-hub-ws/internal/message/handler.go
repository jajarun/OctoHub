package message

import (
	"log"
	"time"
)

// Handler 消息处理器接口
type Handler interface {
	HandleMessage(msg *Message, senderID string, senderType string) *Message
}

// DefaultHandler 默认消息处理器
type DefaultHandler struct{}

// NewDefaultHandler 创建默认消息处理器
func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

// HandleMessage 处理消息
func (h *DefaultHandler) HandleMessage(msg *Message, senderID string, senderType string) *Message {
	log.Printf("收到消息 - Action: %s, From: %s (%s), Data: %v",
		msg.Action, senderID, senderType, msg.Data)

	switch msg.Action {
	case ActionPing:
		return h.handlePing(msg, senderID)
	case ActionEcho:
		return h.handleEcho(msg, senderID)
	case ActionStatus:
		return h.handleStatus(msg, senderID)
	case ActionCommand:
		return h.handleCommand(msg, senderID, senderType)
	default:
		return h.handleUnknownAction(msg, senderID)
	}
}

// handlePing 处理ping消息
func (h *DefaultHandler) handlePing(msg *Message, senderID string) *Message {
	return NewResponse(ActionPong, map[string]interface{}{
		"server_time": time.Now().Unix(),
		"sender_id":   senderID,
	}, msg.RequestID)
}

// handleEcho 处理回显消息
func (h *DefaultHandler) handleEcho(msg *Message, senderID string) *Message {
	return NewResponse(ActionEcho, map[string]interface{}{
		"original_message": msg.Data,
		"echoed_by":        "server",
		"sender_id":        senderID,
	}, msg.RequestID)
}

// handleStatus 处理状态查询
func (h *DefaultHandler) handleStatus(msg *Message, senderID string) *Message {
	statusData := StatusData{
		Online:     true,
		LastActive: time.Now().Unix(),
		ConnectionInfo: map[string]interface{}{
			"connected_at": time.Now().Unix(),
			"client_id":    senderID,
		},
	}

	return NewResponse(ActionStatus, statusData, msg.RequestID)
}

// handleCommand 处理命令消息
func (h *DefaultHandler) handleCommand(msg *Message, senderID string, senderType string) *Message {
	// 解析命令数据
	var cmdData CommandData
	if data, ok := msg.Data.(map[string]interface{}); ok {
		if cmd, exists := data["command"]; exists {
			cmdData.Command = cmd.(string)
		}
		if args, exists := data["args"]; exists {
			if argsMap, ok := args.(map[string]interface{}); ok {
				cmdData.Args = argsMap
			}
		}
	}

	log.Printf("执行命令: %s, 参数: %v, 来自: %s (%s)",
		cmdData.Command, cmdData.Args, senderID, senderType)

	// 这里可以添加具体的命令处理逻辑
	result := map[string]interface{}{
		"command": cmdData.Command,
		"status":  "executed",
		"result":  "命令执行成功",
	}

	return NewResponse(ActionCommand, result, msg.RequestID)
}

// handleUnknownAction 处理未知动作
func (h *DefaultHandler) handleUnknownAction(msg *Message, senderID string) *Message {
	errorData := ErrorData{
		Code:    ErrorCodeUnknownAction,
		Message: "未知的消息动作",
		Details: string(msg.Action),
	}

	return NewResponse(ActionError, errorData, msg.RequestID)
}

// SendSystemMessage 发送系统消息
func SendSystemMessage(action Action, data interface{}) *Message {
	msg := NewMessage(action, data)
	msg.From = "system"
	return msg
}

// SendDisconnectNotification 发送断开连接通知
func SendDisconnectNotification(reason string) *Message {
	data := DisconnectedData{
		Reason: reason,
		Code:   1000, // 正常关闭
	}
	return SendSystemMessage(ActionDisconnected, data)
}

// SendErrorMessage 发送错误消息
func SendErrorMessage(code int, message string, details string) *Message {
	data := ErrorData{
		Code:    code,
		Message: message,
		Details: details,
	}
	return SendSystemMessage(ActionError, data)
}
