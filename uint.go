package packer

import (
	"io"

	"github.com/pierrec/packer/internal/packuint"
	"github.com/pierrec/packer/iobyte"
)

// PackUint64 packs x into buf and returns the number of bytes used.
// buf must be at least 9 bytes long.
func PackUint64(buf []byte, x uint64) int {
	return packuint.PackUint64(buf, x)
}

// PackUint64To packs x to w.
func PackUint64To(w io.Writer, buf []byte, x uint64) error {
	return packuint.PackUint64To(w, buf, x)
}

// UnpackUint64 unpacks buf and returns the value.
func UnpackUint64(bitmap byte, buf []byte) uint64 {
	return packuint.UnpackUint64(bitmap, buf)
}

// UnpackUint64From unpacks an uint32 from r.
func UnpackUint64From(r io.Reader, buf []byte) (uint64, error) {
	return packuint.UnpackUint64From(iobyte.NewReader(r), buf)
}

// PackUint32 packs x into buf and returns the number of bytes used.
// buf must be at least 5 bytes long.
func PackUint32(buf []byte, x uint32) int {
	return packuint.PackUint32(buf, x)
}

// PackUint32To packs x to w.
func PackUint32To(w io.Writer, buf []byte, x uint32) error {
	n := packuint.PackUint32(buf, x)
	_, err := w.Write(buf[:n])
	return err
}

// UnpackUint32 unpacks buf and returns the value.
func UnpackUint32(bitmap byte, buf []byte) uint32 {
	return packuint.UnpackUint32(bitmap, buf)
}

// UnpackUint32From unpacks an uint32 from r.
func UnpackUint32From(r io.Reader, buf []byte) (uint32, error) {
	return packuint.UnpackUint32From(iobyte.NewReader(r), buf)
}
