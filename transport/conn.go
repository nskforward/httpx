package transport

import (
	"net"
	"time"
)

type Conn struct {
	IP string
	net.Conn
	Listener *Listener
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}

func (c *Conn) Close() error {
	c.Listener.mx.Lock()
	counter, ok := c.Listener.counters[c.IP]
	c.Listener.mx.Unlock()
	if ok {
		v := counter.Dec()
		if v < 1 {
			c.Listener.mx.Lock()
			delete(c.Listener.counters, c.IP)
			c.Listener.mx.Unlock()
		}
	}
	return c.Conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}
