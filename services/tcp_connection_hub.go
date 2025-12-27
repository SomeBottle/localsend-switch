package services

// TCP 连接管理模块

import (
	"errors"
	"net"
	"sync"
)

type TCPConnectionHub struct {
	// 控制对 conns 的并发访问
	mutex sync.Mutex
	conns map[string]*net.TCPConn
}

// NewTCPConnectionHub 创建一个新的 TCP 连接管理器
func NewTCPConnectionHub() *TCPConnectionHub {
	return &TCPConnectionHub{
		conns: make(map[string]*net.TCPConn),
	}
}

// AddConnection 添加一个新的 TCP 连接到管理器
func (hub *TCPConnectionHub) AddConnection(conn *net.TCPConn) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	// 使用连接发起地址 (含有端口) 作为键 (标记客户端)
	remoteAddrStr := conn.RemoteAddr().String()
	if _, exists := hub.conns[remoteAddrStr]; exists {
		return errors.New("Connection already exists")
	}
	hub.conns[remoteAddrStr] = conn
	return nil
}

// RemoveConnection 从管理器中移除一个 TCP 连接
func (hub *TCPConnectionHub) RemoveConnection(conn *net.TCPConn) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	remoteAddrStr := conn.RemoteAddr().String()
	delete(hub.conns, remoteAddrStr)
}

// NumConnections 返回当前管理的连接数
func (hub *TCPConnectionHub) NumConnections() int {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	return len(hub.conns)
}

// Close 关闭所有管理的 TCP 连接
func (hub *TCPConnectionHub) Close() {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	for _, conn := range hub.conns {
		// 连接关闭后，连接 handler 会自动从管理器中移除该连接
		conn.Close()
	}
}
