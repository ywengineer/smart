package pkg

import "github.com/panjf2000/gnet/v2"

type gnetReader struct {
	conn gnet.Conn
}

func (gr *gnetReader) Peek(n int) (buf []byte, err error) {
	return gr.conn.Peek(n)
}

func (gr *gnetReader) Skip(n int) (err error) {
	_, err = gr.conn.Discard(n)
	return err
}

func (gr *gnetReader) ReadBinary(n int) (p []byte, err error) {
	return gr.conn.Next(n)
}

func (gr *gnetReader) Release() (err error) {
	return nil
}

func (gr *gnetReader) Len() (length int) {
	return gr.conn.InboundBuffered()
}
