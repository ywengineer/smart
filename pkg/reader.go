package pkg

type Reader interface {
	// Peek returns the next n bytes without advancing the reader.
	// Other behavior is the same as Next.
	Peek(n int) (buf []byte, err error)

	// Skip the next n bytes and advance the reader, which is
	// a faster implementation of Next when the next data is not used.
	Skip(n int) (err error)

	// ReadBinary is a faster implementation of Next when it needs to
	// return a copy of the slice that is not shared with the underlying layer.
	// It replaces:
	//
	//  var p, err = Next(n)
	//  var b = make([]byte, n)
	//  copy(b, p)
	//  return b, err
	//
	ReadBinary(n int) (p []byte, err error)

	// Release the memory space occupied by all read slices. This method needs to be executed actively to
	// recycle the memory after confirming that the previously read data is no longer in use.
	// After invoking Release, the slices obtained by the method such as Next, Peek, Skip will
	// become an invalid address and cannot be used anymore.
	Release() (err error)

	// Len returns the total length of the readable data in the reader.
	Len() (length int)
}
