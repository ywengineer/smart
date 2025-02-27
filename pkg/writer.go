package pkg

type Writer interface {
	// WriteBinary is a faster implementation of Malloc when a slice needs to be written.
	// It replaces:
	//
	//  var buf, err = Malloc(len(b))
	//  n = copy(buf, b)
	//  return n, err
	//
	// The argument slice b will be referenced based on the original address and will not be copied,
	// so make sure that the slice b will not be changed.
	WriteBinary(b []byte) (n int, err error)

	// Flush will submit all malloc data and must confirm that the allocated bytes have been correctly assigned.
	// Its behavior is equivalent to the io.Writer hat already has parameters(slice b).
	Flush() (err error)
}
