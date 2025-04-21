package autodetect

import (
	"net"
	"time"
)

// PrependConn wraps a net.Conn and "prepends" a buffer before the actual conn.
// Used to restore bytes read during ProtocolMatch peek.
type PrependConn struct {
	net.Conn
	buf     []byte
	readPos int
}

// NewPrependConn 构造函数：传入原始 conn 和已读出的数据 buf
func NewPrependConn(conn net.Conn, buf []byte) net.Conn {
	return &PrependConn{
		Conn:    conn,
		buf:     buf,
		readPos: 0,
	}
}

// 重点：优先从 buf 中读，读完后再从真实连接读取
func (p *PrependConn) Read(b []byte) (int, error) {
	if p.readPos < len(p.buf) {
		n := copy(b, p.buf[p.readPos:])
		p.readPos += n
		return n, nil
	}
	return p.Conn.Read(b)
}

// 其余方法均委托给原始连接

func (p *PrependConn) Write(b []byte) (int, error) {
	return p.Conn.Write(b)
}

func (p *PrependConn) Close() error {
	return p.Conn.Close()
}

func (p *PrependConn) LocalAddr() net.Addr {
	return p.Conn.LocalAddr()
}

func (p *PrependConn) RemoteAddr() net.Addr {
	return p.Conn.RemoteAddr()
}

func (p *PrependConn) SetDeadline(t time.Time) error {
	return p.Conn.SetDeadline(t)
}

func (p *PrependConn) SetReadDeadline(t time.Time) error {
	return p.Conn.SetReadDeadline(t)
}

func (p *PrependConn) SetWriteDeadline(t time.Time) error {
	return p.Conn.SetWriteDeadline(t)
}
