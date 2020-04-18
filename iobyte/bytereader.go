package iobyte

import (
	"bufio"
	"io"
)

// ByteReader combines several io.Reader related interfaces.
type ByteReader interface {
	io.Reader
	io.ByteReader
}

// NewReader efficiently returns a ByteReader based on the given io.Reader.
func NewReader(r io.Reader) ByteReader {
	if br, ok := r.(ByteReader); ok {
		return br
	}
	return bufio.NewReader(r)
}

// ByteWriter combines several io.Writer related interfaces.
type ByteWriter interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
}

type byteWriter struct {
	*bufio.Writer
}

func noop(*error) {}

// NewWriter efficiently returns a ByteWriter based on the given io.Writer.
// ``done`` is to be called once writing to w is complete.
func NewWriter(w io.Writer) (_ ByteWriter, done func(*error)) {
	if bw, ok := w.(ByteWriter); ok {
		return bw, noop
	}
	b := bufio.NewWriter(w)
	return &byteWriter{b}, func(errp *error) {
		// Flush the bufio.Writer if there is no error already.
		if errp == nil || *errp != nil {
			return
		}
		if err := b.Flush(); err != nil {
			*errp = err
		}
	}
}
