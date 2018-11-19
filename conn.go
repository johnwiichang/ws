package ws

import (
	"net"
	"net/http"
)

//RemoteAddr 远程计算机地址
func (conn *Connection) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

//LocalAddr 本地计算机地址
func (conn *Connection) LocalAddr() net.Addr {
	return conn.conn.LocalAddr()
}

//OriginalRequest 原始 HTTP 请求实例
func (conn *Connection) OriginalRequest() *http.Request {
	return conn.conn.Request()
}
