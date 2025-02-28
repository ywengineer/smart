package pkg

import (
	"github.com/panjf2000/gnet/v2"
)

func NetGNetConn(conn gnet.Conn) Conn {
	return &gnetConn{
		conn: conn,
		w:    &gnetWriter{conn: conn},
		r:    &gnetReader{conn: conn},
	}
}

type gnetConn struct {
	conn gnet.Conn
	w    Writer
	r    Reader
}

func (c *gnetConn) Fd() int {
	return c.conn.Fd()
}

func (c *gnetConn) Reader() Reader {
	return c.r
}

func (c *gnetConn) Writer() Writer {
	return c.w
}

func (c *gnetConn) Close() error {
	return c.conn.Close()
}
