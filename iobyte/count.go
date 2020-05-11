package iobyte

import "io"

// CountWriter records the number of bytes written to its source.
type CountWriter struct {
	N int64 // N contains the number of bytes written to its writer
	io.Writer
}

func (c *CountWriter) Write(p []byte) (int, error) {
	c.N += int64(len(p))
	return c.Writer.Write(p)
}

// CountReader records the number of bytes read from its Reader source.
type CountReader struct {
	N int64 // N contains the number of bytes read from its reader
	io.Reader
}

func (c *CountReader) Read(p []byte) (int, error) {
	n, err := c.Reader.Read(p)
	c.N += int64(n)
	return n, err
}

// CountByteReader records the number of bytes read from its ByteReader source.
type CountByteReader struct {
	N int64
	ByteReader
}

func (c *CountByteReader) Read(p []byte) (int, error) {
	n, err := c.ByteReader.Read(p)
	c.N += int64(n)
	return n, err
}

func (c *CountByteReader) ReadByte() (byte, error) {
	b, err := c.ByteReader.ReadByte()
	c.N++
	return b, err
}
