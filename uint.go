package packer

import (
	"io"

	"github.com/pierrec/packer/internal/packuint"
	"github.com/pierrec/packer/iobyte"
)

// PackUint64 packs x into buf and returns the number of bytes used.
//
// buf is used as scratch space if it has at least 9 bytes in capacity.
func PackUint64(buf []byte, x uint64) int {
	if cap(buf) < 9 {
		buf = make([]byte, 9)
	}
	return packuint.PackUint64(buf, x)
}

// PackUint64To packs x to w.
//
// buf is used as scratch space if it has at least 9 bytes in capacity.
func PackUint64To(w io.Writer, buf []byte, x uint64) error {
	n := PackUint64(buf, x)
	_, err := w.Write(buf[:n])
	return err
}

// UnpackUint64 unpacks buf and returns the value.
func UnpackUint64(buf []byte) uint64 {
	return packuint.UnpackUint64(buf[0], buf[1:])
}

// UnpackUint64From unpacks an uint32 from r.
func UnpackUint64From(r io.Reader, buf []byte) (uint64, error) {
	return packuint.UnpackUint64From(iobyte.NewReader(r), buf)
}

// PackUint32 packs x into buf and returns the number of bytes used.
//
// buf is used as scratch space if it has at least 5 bytes in capacity.
func PackUint32(buf []byte, x uint32) int {
	if cap(buf) < 5 {
		buf = make([]byte, 5)
	}
	return packuint.PackUint32(buf, x)
}

// PackUint32To packs x to w.
//
// buf is used as scratch space if it has at least 5 bytes in capacity.
func PackUint32To(w io.Writer, buf []byte, x uint32) error {
	n := PackUint32(buf, x)
	_, err := w.Write(buf[:n])
	return err
}

// UnpackUint32 unpacks buf and returns the value.
func UnpackUint32(buf []byte) uint32 {
	return packuint.UnpackUint32(buf[0], buf[1:])
}

// UnpackUint32From unpacks an uint32 from r.
func UnpackUint32From(r io.Reader, buf []byte) (uint32, error) {
	return packuint.UnpackUint32From(iobyte.NewReader(r), buf)
}
