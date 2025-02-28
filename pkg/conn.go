package pkg

type Conn interface {
	// Reader The recommended API for nocopy reading and writing.
	// Reader will return nocopy buffer data, or error after timeout which set by SetReadTimeout.
	Reader() Reader

	// Writer will write data to the connection by NIO mode,
	// so it will return an error only when the connection isn't Active.
	Writer() Writer

	Close() error

	Fd() int
}
