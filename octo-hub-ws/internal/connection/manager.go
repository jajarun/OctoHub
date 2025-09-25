package connection

import (
	"OctoHub/Ws/internal/message"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Manager 连接管理器
type Manager struct {
	// 用户连接映射 user_id -> Connection
	userConnections map[string]*Connection
	// PC连接映射 pc_id -> Connection
	pcConnections map[string]*Connection
	// 读写锁
	mutex sync.RWMutex
	// 清理定时器
	cleanupTicker *time.Ticker
	// 停止清理通道
	stopCleanup chan bool
	// 待清理连接队列 (避免长时间持锁)
	cleanupQueue chan *Connection
}

// NewManager 创建新的连接管理器
func NewManager(cleanupInterval, staleTimeout int) *Manager {
	m := &Manager{
		userConnections: make(map[string]*Connection),
		pcConnections:   make(map[string]*Connection),
		cleanupTicker:   time.NewTicker(time.Duration(cleanupInterval) * time.Second),
		stopCleanup:     make(chan bool),
		cleanupQueue:    make(chan *Connection, 1000), // 缓冲队列
	}

	// 启动定期清理协程
	go m.startCleanupRoutine(time.Duration(staleTimeout) * time.Second)
	// 启动异步清理处理协程
	go m.processCleanupQueue()

	return m
}

// AddConnection 添加连接
func (m *Manager) AddConnection(conn *Connection) *Connection {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var oldConn *Connection
	var exists bool

	switch conn.Type {
	case USER:
		oldConn, exists = m.userConnections[conn.ID]
		m.userConnections[conn.ID] = conn
	case PC:
		oldConn, exists = m.pcConnections[conn.ID]
		m.pcConnections[conn.ID] = conn
	}

	if exists && oldConn != nil {
		// 关闭旧连接并发送通知
		m.closeOldConnection(oldConn, "新连接建立")
	}

	log.Printf("%s连接建立: %s", conn.Type.String(), conn.ID)
	return oldConn
}

// RemoveConnection 移除连接
func (m *Manager) RemoveConnection(conn *Connection) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch conn.Type {
	case USER:
		delete(m.userConnections, conn.ID)
		log.Printf("用户连接断开: %s", conn.ID)
	case PC:
		delete(m.pcConnections, conn.ID)
		log.Printf("PC连接断开: %s", conn.ID)
	}
}

// GetConnection 获取连接
func (m *Manager) GetConnection(connType ConnectionType, id string) (*Connection, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	switch connType {
	case USER:
		conn, exists := m.userConnections[id]
		return conn, exists
	case PC:
		conn, exists := m.pcConnections[id]
		return conn, exists
	default:
		return nil, false
	}
}

// closeOldConnection 关闭旧连接并发送通知
func (m *Manager) closeOldConnection(oldConn *Connection, reason string) {
	if oldConn != nil {
		// 发送断开通知
		disconnectMsg := message.SendDisconnectNotification(reason)
		if data, err := disconnectMsg.ToJSON(); err == nil {
			select {
			case oldConn.SendChan <- data:
			default:
			}
		}

		// 关闭连接
		oldConn.Close()
		log.Printf("旧连接已断开: %s, 原因: %s", oldConn.ID, reason)
	}
}

// startCleanupRoutine 启动定期清理协程
func (m *Manager) startCleanupRoutine(staleThreshold time.Duration) {
	for {
		select {
		case <-m.cleanupTicker.C:
			m.cleanupStaleConnections(staleThreshold)
		case <-m.stopCleanup:
			m.cleanupTicker.Stop()
			return
		}
	}
}

// cleanupStaleConnections 清理过期连接 (优化版本)
func (m *Manager) cleanupStaleConnections(staleThreshold time.Duration) {
	now := time.Now()

	// 分批收集需要检查的连接，减少锁持有时间
	var candidatesForCleanup []*Connection

	// 快速扫描，只收集候选连接
	m.mutex.RLock()
	for _, conn := range m.userConnections {
		if now.Sub(conn.LastActive) > staleThreshold {
			candidatesForCleanup = append(candidatesForCleanup, conn)
		}
	}
	for _, conn := range m.pcConnections {
		if now.Sub(conn.LastActive) > staleThreshold {
			candidatesForCleanup = append(candidatesForCleanup, conn)
		}
	}
	m.mutex.RUnlock()

	// 异步处理候选连接，避免阻塞主流程
	for _, conn := range candidatesForCleanup {
		select {
		case m.cleanupQueue <- conn:
			// 成功加入清理队列
		default:
			// 队列满了，跳过这次清理
			log.Printf("清理队列已满，跳过连接: %s", conn.ID)
		}
	}
}

// processCleanupQueue 处理清理队列
func (m *Manager) processCleanupQueue() {
	for {
		select {
		case conn := <-m.cleanupQueue:
			m.processConnectionCleanup(conn)
		case <-m.stopCleanup:
			return
		}
	}
}

// processConnectionCleanup 处理单个连接的清理
func (m *Manager) processConnectionCleanup(conn *Connection) {
	// 先检查连接是否还活着（不持锁）
	if m.isConnectionAlive(conn) {
		// 连接还活着，更新活跃时间
		conn.UpdateActivity()
		return
	}

	// 连接已死，需要清理
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 二次检查连接是否还在映射中（可能已被其他地方清理）
	switch conn.Type {
	case USER:
		if existingConn, exists := m.userConnections[conn.ID]; exists && existingConn == conn {
			log.Printf("清理过期用户连接: %s, 最后活跃: %v", conn.ID, conn.LastActive)
			conn.Close()
			delete(m.userConnections, conn.ID)
		}
	case PC:
		if existingConn, exists := m.pcConnections[conn.ID]; exists && existingConn == conn {
			log.Printf("清理过期PC连接: %s, 最后活跃: %v", conn.ID, conn.LastActive)
			conn.Close()
			delete(m.pcConnections, conn.ID)
		}
	}
}

// isConnectionAlive 检查连接是否还活着 (轻量级检查)
func (m *Manager) isConnectionAlive(conn *Connection) bool {
	// 首先检查连接是否已经被标记为关闭
	select {
	case <-conn.CloseChan:
		return false
	default:
	}

	// 检查连接的底层状态，避免发送ping造成额外开销
	// 如果连接最近有活动（比如1分钟内），认为是活跃的
	if time.Since(conn.LastActive) < time.Minute {
		return true
	}

	// 对于长时间无活动的连接，进行轻量级检查
	// 尝试设置写入超时，如果连接已断开会立即返回错误
	conn.Conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	err := conn.Conn.WriteMessage(websocket.PingMessage, nil)

	if err != nil {
		// 连接已断开
		return false
	}

	// ping发送成功，连接可能还活着
	return true
}

// Stop 停止连接管理器
func (m *Manager) Stop() {
	close(m.stopCleanup)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 关闭所有连接
	for _, conn := range m.userConnections {
		conn.Close()
	}
	for _, conn := range m.pcConnections {
		conn.Close()
	}

	log.Println("连接管理器已停止")
}

// GetConnectionStats 获取连接统计信息 (优化版本)
func (m *Manager) GetConnectionStats() map[string]interface{} {
	m.mutex.RLock()

	// 快速获取基本统计
	userCount := len(m.userConnections)
	pcCount := len(m.pcConnections)

	// 只在连接数较少时返回详细信息，避免性能问题
	var userStats, pcStats []map[string]interface{}

	if userCount <= 100 { // 连接数少时返回详细信息
		userStats = make([]map[string]interface{}, 0, userCount)
		for id, conn := range m.userConnections {
			userStats = append(userStats, map[string]interface{}{
				"id":          id,
				"last_active": conn.LastActive.Unix(),
				"duration":    time.Since(conn.LastActive).Seconds(),
			})
		}
	}

	if pcCount <= 100 {
		pcStats = make([]map[string]interface{}, 0, pcCount)
		for id, conn := range m.pcConnections {
			pcStats = append(pcStats, map[string]interface{}{
				"id":          id,
				"last_active": conn.LastActive.Unix(),
				"duration":    time.Since(conn.LastActive).Seconds(),
			})
		}
	}

	m.mutex.RUnlock()

	result := map[string]interface{}{
		"user_connections":   userCount,
		"pc_connections":     pcCount,
		"total_connections":  userCount + pcCount,
		"cleanup_queue_size": len(m.cleanupQueue),
		"timestamp":          time.Now().Unix(),
	}

	// 只在连接数较少时包含详细信息
	if userStats != nil {
		result["user_details"] = userStats
	}
	if pcStats != nil {
		result["pc_details"] = pcStats
	}

	return result
}
