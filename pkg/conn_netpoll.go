package pkg

import "github.com/cloudwego/netpoll"

func NetNetpollConn(conn netpoll.Connection) Conn {
	return &netpollConn{conn}
}

type netpollConn struct {
	conn netpoll.Connection
}

func (c *netpollConn) Reader() Reader {
	return c.conn.Reader()
}

func (c *netpollConn) Writer() Writer {
	return c.conn.Writer()
}

func (c *netpollConn) Close() error {
	return c.conn.Close()
}
