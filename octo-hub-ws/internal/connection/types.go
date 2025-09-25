package connection

import (
	"time"

	"github.com/gorilla/websocket"
)

// ConnectionType 连接类型枚举
type ConnectionType int

const (
	USER ConnectionType = iota
	PC
)

// String 返回连接类型的字符串表示
func (ct ConnectionType) String() string {
	switch ct {
	case USER:
		return "user"
	case PC:
		return "pc"
	default:
		return "unknown"
	}
}

// Connection 连接信息结构体
type Connection struct {
	ID         string          // user_id 或 pc_id
	Type       ConnectionType  // 连接类型
	Conn       *websocket.Conn // WebSocket连接
	SendChan   chan []byte     // 发送消息通道
	CloseChan  chan bool       // 关闭通道
	LastActive time.Time       // 最后活跃时间
}

// NewConnection 创建新连接
func NewConnection(id string, connType ConnectionType, wsConn *websocket.Conn, bufferSize int) *Connection {
	return &Connection{
		ID:         id,
		Type:       connType,
		Conn:       wsConn,
		SendChan:   make(chan []byte, bufferSize),
		CloseChan:  make(chan bool),
		LastActive: time.Now(),
	}
}

// UpdateActivity 更新最后活跃时间
func (c *Connection) UpdateActivity() {
	c.LastActive = time.Now()
}

// Close 关闭连接
func (c *Connection) Close() {
	close(c.CloseChan)
	c.Conn.Close()
}
