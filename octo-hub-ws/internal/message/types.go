package message

import (
	"encoding/json"
	"time"
)

// Action 消息动作枚举
type Action string

const (
	// 系统消息
	ActionConnected    Action = "connected"    // 连接成功
	ActionDisconnected Action = "disconnected" // 连接断开
	ActionPing         Action = "ping"         // 心跳检测
	ActionPong         Action = "pong"         // 心跳响应
	ActionError        Action = "error"        // 错误消息

	// 业务消息
	ActionEcho         Action = "echo"         // 回显消息
	ActionBroadcast    Action = "broadcast"    // 广播消息
	ActionPrivate      Action = "private"      // 私聊消息
	ActionNotification Action = "notification" // 通知消息

	// 控制消息
	ActionCommand     Action = "command"     // 执行命令
	ActionStatus      Action = "status"      // 状态查询
	ActionSubscribe   Action = "subscribe"   // 订阅
	ActionUnsubscribe Action = "unsubscribe" // 取消订阅
)

// Message WebSocket消息标准格式
type Message struct {
	Action    Action      `json:"action"`               // 消息动作
	Data      interface{} `json:"data,omitempty"`       // 消息数据，可选
	Timestamp int64       `json:"timestamp"`            // 时间戳
	RequestID string      `json:"request_id,omitempty"` // 请求ID，用于响应匹配
	From      string      `json:"from,omitempty"`       // 发送者ID
	To        string      `json:"to,omitempty"`         // 接收者ID，用于私聊
}

// NewMessage 创建新消息
func NewMessage(action Action, data interface{}) *Message {
	return &Message{
		Action:    action,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewMessageWithFrom 创建带发送者的消息
func NewMessageWithFrom(action Action, data interface{}, from string) *Message {
	return &Message{
		Action:    action,
		Data:      data,
		Timestamp: time.Now().Unix(),
		From:      from,
	}
}

// NewResponse 创建响应消息
func NewResponse(action Action, data interface{}, requestID string) *Message {
	return &Message{
		Action:    action,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: requestID,
	}
}

// ToJSON 转换为JSON字符串
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON 从JSON字符串解析消息
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// 预定义的常用消息数据结构

// ConnectedData 连接成功数据
type ConnectedData struct {
	UserID     string `json:"user_id,omitempty"`
	PCID       string `json:"pc_id,omitempty"`
	SessionID  string `json:"session_id"`
	ServerTime int64  `json:"server_time"`
}

// DisconnectedData 断开连接数据
type DisconnectedData struct {
	Reason string `json:"reason"`
	Code   int    `json:"code"`
}

// ErrorData 错误数据
type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// CommandData 命令数据
type CommandData struct {
	Command string                 `json:"command"`
	Args    map[string]interface{} `json:"args,omitempty"`
}

// StatusData 状态数据
type StatusData struct {
	Online         bool                   `json:"online"`
	LastActive     int64                  `json:"last_active"`
	ConnectionInfo map[string]interface{} `json:"connection_info,omitempty"`
}

// NotificationData 通知数据
type NotificationData struct {
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Type    string                 `json:"type"` // info, warning, error, success
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

// 错误代码常量
const (
	ErrorCodeInvalidMessage   = 1001
	ErrorCodeUnknownAction    = 1002
	ErrorCodePermissionDenied = 1003
	ErrorCodeInternalError    = 1004
	ErrorCodeRateLimited      = 1005
)
