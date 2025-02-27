package pkg

import "github.com/panjf2000/gnet/v2"

type gnetWriter struct {
	conn gnet.Conn
}

func (g *gnetWriter) WriteBinary(b []byte) (n int, err error) {
	return g.conn.Write(b)
}

func (g *gnetWriter) Flush() (err error) {
	return g.conn.Flush()
}
