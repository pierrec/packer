package iobyte

import "io"

// CountWriter records the number of bytes written to its source.
type CountWriter struct {
	N int64 // N contains the number of bytes written to its writer
	W io.Writer
}

func (c *CountWriter) Write(p []byte) (int, error) {
	c.N += int64(len(p))
	return c.W.Write(p)
}

// CountReader records the number of bytes read from its source.
type CountReader struct {
	N int64 // N contains the number of bytes read from its reader
	R io.Reader
}

func (c *CountReader) Read(p []byte) (int, error) {
	n, err := c.R.Read(p)
	c.N += int64(n)
	return n, err
}
