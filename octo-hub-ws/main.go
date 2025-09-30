package main

import (
	"log"
	"time"

	"OctoHub/Ws/internal/config"
	"OctoHub/Ws/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 检查是否使用默认密钥
	if cfg.IsDefaultSignatureKey() {
		log.Println("警告: 使用默认签名密钥，生产环境请设置OCTOHUB_SIGNATURE_KEY环境变量")
	}

	// 创建WebSocket服务器实例
	wsServer := server.NewWebSocketServer(cfg)

	// 创建Gin路由器
	r := gin.Default()

	// 设置路由
	r.GET("/ws/user", wsServer.HandleUserConnection)
	r.GET("/ws/node", wsServer.HandleNodeConnection)
	r.GET("/status", wsServer.GetDetailedConnectionStats)

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	log.Printf("WebSocket服务器启动在端口 %s", cfg.Server.Port)
	log.Printf("用户连接地址: ws://localhost:%s/ws/user?user_id=xxx&timestamp=xxx&signature=xxx", cfg.Server.Port)
	log.Printf("PC连接地址: ws://localhost:%s/ws/pc?pc_id=xxx&timestamp=xxx&signature=xxx", cfg.Server.Port)
	log.Printf("状态查询地址: http://localhost:%s/status", cfg.Server.Port)
	log.Printf("签名超时时间: %v", cfg.GetSignatureTimeout())
	log.Printf("最大连接数: %d", cfg.Connection.MaxConnections)

	// 启动服务器
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
