package server

import (
	"log"
	"net/http"

	"OctoHub/Ws/internal/auth"
	"OctoHub/Ws/internal/config"
	"OctoHub/Ws/internal/connection"
	"OctoHub/Ws/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketServer WebSocket服务器
type WebSocketServer struct {
	config        *config.Config
	connManager   *connection.Manager
	wsHandler     *handler.WebSocketHandler
	authValidator *auth.SignatureValidator
	upgrader      websocket.Upgrader
}

// NewWebSocketServer 创建新的WebSocket服务器
func NewWebSocketServer(cfg *config.Config) *WebSocketServer {
	connManager := connection.NewManager(cfg.Connection.CleanupInterval, cfg.Connection.StaleTimeout)
	wsHandler := handler.NewWebSocketHandler(
		connManager,
		cfg.GetReadTimeout(),
		cfg.GetWriteTimeout(),
	)
	authValidator := auth.NewSignatureValidator(cfg.Signature.Key)

	return &WebSocketServer{
		config:        cfg,
		connManager:   connManager,
		wsHandler:     wsHandler,
		authValidator: authValidator,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许跨域连接
			},
		},
	}
}

// HandleUserConnection 处理用户连接
func (s *WebSocketServer) HandleUserConnection(c *gin.Context) {
	userID := c.Query("user_id")
	timestamp := c.Query("timestamp")
	signature := c.Query("signature")

	if userID == "" || timestamp == "" || signature == "" {
		c.JSON(400, gin.H{"error": "缺少必要参数"})
		return
	}

	// 验证签名
	if !s.authValidator.ValidateSignature(userID, timestamp, signature, s.config.Signature.Timeout) {
		c.JSON(401, gin.H{"error": "签名验证失败"})
		return
	}

	// 升级为WebSocket连接
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 创建连接对象
	wsConn := connection.NewConnection(userID, connection.USER, conn, s.config.Connection.BufferSize)

	// 添加到连接管理器
	s.connManager.AddConnection(wsConn)

	// 发送连接成功消息
	s.wsHandler.SendConnectedMessage(wsConn)

	// 启动连接处理
	go s.wsHandler.HandleConnection(wsConn)
}

// HandleNodeConnection 处理Node节点连接
func (s *WebSocketServer) HandleNodeConnection(c *gin.Context) {
	pcID := c.Query("pc_id")
	timestamp := c.Query("timestamp")
	signature := c.Query("signature")

	if pcID == "" || timestamp == "" || signature == "" {
		c.JSON(400, gin.H{"error": "缺少必要参数"})
		return
	}

	// 验证签名
	if !s.authValidator.ValidateSignature(pcID, timestamp, signature, s.config.Signature.Timeout) {
		c.JSON(401, gin.H{"error": "签名验证失败"})
		return
	}

	// 升级为WebSocket连接
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 创建连接对象
	wsConn := connection.NewConnection(pcID, connection.PC, conn, s.config.Connection.BufferSize)

	// 添加到连接管理器
	s.connManager.AddConnection(wsConn)

	// 发送连接成功消息
	s.wsHandler.SendConnectedMessage(wsConn)

	// 启动连接处理
	go s.wsHandler.HandleConnection(wsConn)
}

// GetDetailedConnectionStats 获取详细连接统计信息
func (s *WebSocketServer) GetDetailedConnectionStats(c *gin.Context) {
	stats := s.connManager.GetConnectionStats()
	c.JSON(200, stats)
}
