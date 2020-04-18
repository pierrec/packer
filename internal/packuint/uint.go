package packuint

import (
	"encoding/binary"
	"io"
	"math/bits"

	"github.com/pierrec/packer/iobyte"
)

//go:generate go run uintgen.go

// Pack and unpack integers:
//  - first byte contains a bitmap of the non zero bytes found in the integer
//  - following bytes are the integer non zero bytes

// PackUint64 packs x into buf and returns the number of bytes used.
// buf must be at least 9 bytes long.
func PackUint64(buf []byte, x uint64) int {
	_ = buf[8]
	if x == 0 {
		buf[0] = 0
		return 1
	}
	const (
		shift = 8
		max   = 1 << shift
	)
	var (
		bitmap uint8
		i      int
	)

	if x := byte(x); x > 0 {
		bitmap = 1 << 7
		i = 1
		buf[1] = x
	}
	if x < max {
		buf[0] = bitmap
		return 2
	}
	x >>= shift

	if x := byte(x); x > 0 {
		bitmap |= 1 << 6
		i++
		buf[i] = x
	}
	if x < max {
		buf[0] = bitmap
		return i + 1
	}
	x >>= shift

	if x := byte(x); x > 0 {
		bitmap |= 1 << 5
		i++
		buf[i] = x
	}
	if x < max {
		buf[0] = bitmap
		return i + 1
	}
	x >>= shift

	if x := byte(x); x > 0 {
		bitmap |= 1 << 4
		i++
		buf[i] = x
	}
	if x < max {
		buf[0] = bitmap
		return i + 1
	}
	x >>= shift

	if x := byte(x); x > 0 {
		bitmap |= 1 << 3
		i++
		buf[i] = x
	}
	if x < max {
		buf[0] = bitmap
		return i + 1
	}
	x >>= shift

	if x := byte(x); x > 0 {
		bitmap |= 1 << 2
		i++
		buf[i] = x
	}
	if x < max {
		buf[0] = bitmap
		return i + 1
	}
	x >>= shift

	if x := byte(x); x > 0 {
		bitmap |= 1 << 1
		i++
		buf[i] = x
	}
	if x < max {
		buf[0] = bitmap
		return i + 1
	}
	buf[0] = bitmap | 1
	i++
	buf[i] = byte(x >> shift)
	return i + 1
}

func PackUint64To(w io.Writer, buf []byte, x uint64) error {
	n := PackUint64(buf, x)
	_, err := w.Write(buf[:n])
	return err
}

// UnpackUint64 unpacks buf and returns the value.
func UnpackUint64(bitmap byte, buf []byte) uint64 {
	switch bitmap {
	case 0:
		return 0
	case 0xFF:
		return binary.LittleEndian.Uint64(buf)
	}
	entry := unpackTable[bitmap-1]
	a, b, c, d, e, f, g := entry.A()*8, entry.B()*8, entry.C()*8, entry.D()*8, entry.E()*8, entry.F()*8, entry.G()*8
	switch entry.Num() {
	case 1:
		return uint64(buf[0]) << a
	case 2:
		_ = buf[1]
		return uint64(buf[0])<<a | uint64(buf[1])<<b
	case 3:
		_ = buf[2]
		return uint64(buf[0])<<a | uint64(buf[1])<<b | uint64(buf[2])<<c
	case 4:
		_ = buf[3]
		return uint64(buf[0])<<a | uint64(buf[1])<<b | uint64(buf[2])<<c | uint64(buf[3])<<d
	case 5:
		_ = buf[4]
		return uint64(buf[0])<<a | uint64(buf[1])<<b | uint64(buf[2])<<c | uint64(buf[3])<<d |
			uint64(buf[4])<<e
	case 6:
		_ = buf[5]
		return uint64(buf[0])<<a | uint64(buf[1])<<b | uint64(buf[2])<<c | uint64(buf[3])<<d |
			uint64(buf[4])<<e | uint64(buf[5])<<f
	}
	_ = buf[6]
	return uint64(buf[0])<<a | uint64(buf[1])<<b | uint64(buf[2])<<c | uint64(buf[3])<<d |
		uint64(buf[4])<<e | uint64(buf[5])<<f | uint64(buf[6])<<g
}

func UnpackUint64From(r iobyte.ByteReader, buf []byte) (x uint64, err error) {
	bitmap, err := r.ReadByte()
	if err != nil {
		return
	}
	if bitmap == 0 {
		return
	}
	n := bits.OnesCount8(bitmap)
	if n == 1 {
		buf[0], err = r.ReadByte()
	} else {
		_, err = io.ReadFull(r, buf[:n])
	}
	if err != nil {
		return
	}
	return UnpackUint64(bitmap, buf), nil
}

// PackUint32 packs x into buf and returns the number of bytes used.
// buf must be at least 5 bytes long.
func PackUint32(buf []byte, x uint32) int {
	_ = buf[4]
	if x == 0 {
		buf[0] = 0
		return 1
	}
	const (
		shift = 4
		max   = 1 << shift
		mask  = max - 1
	)
	var (
		acc    uint32
		i      int
		bitmap uint8
	)

	if x := x & mask; x > 0 {
		bitmap = 1 << 7
		acc = x
		i = 1
	}
	if x < max {
		buf[0] = bitmap
		buf[1] = byte(acc)
		return 2
	}
	x >>= shift

	if x := x & mask; x > 0 {
		bitmap |= 1 << 6
		acc |= x << (i * 4)
		i++
	}
	if x < max {
		buf[0] = bitmap
		buf[1] = byte(acc)
		return 2
	}
	x >>= shift

	if x := x & mask; x > 0 {
		bitmap |= 1 << 5
		acc |= x << (i * 4)
		i++
	}
	if x < max {
		buf[0] = bitmap
		buf[1] = byte(acc)
		buf[2] = byte(acc >> 8)
		return 1 + (i+1)/2
	}
	x >>= shift

	if x := x & mask; x > 0 {
		bitmap |= 1 << 4
		acc |= x << (i * 4)
		i++
	}
	if x < max {
		buf[0] = bitmap
		buf[1] = byte(acc)
		buf[2] = byte(acc >> 8)
		return 1 + (i+1)/2
	}
	x >>= shift

	if x := x & mask; x > 0 {
		bitmap |= 1 << 3
		acc |= x << (i * 4)
		i++
	}
	if x < max {
		buf[0] = bitmap
		buf[1] = byte(acc)
		buf[2] = byte(acc >> 8)
		buf[3] = byte(acc >> 16)
		return 1 + (i+1)/2
	}
	x >>= shift

	if x := x & mask; x > 0 {
		bitmap |= 1 << 2
		acc |= x << (i * 4)
		i++
	}
	if x < max {
		buf[0] = bitmap
		buf[1] = byte(acc)
		buf[2] = byte(acc >> 8)
		buf[3] = byte(acc >> 16)
		return 1 + (i+1)/2
	}
	x >>= shift

	if x := x & mask; x > 0 {
		bitmap |= 1 << 1
		acc |= x << (i * 4)
		i++
	}
	x >>= shift
	if x := x & mask; x > 0 {
		bitmap |= 1
		acc |= x << (i * 4)
		i++
	}

	buf[0] = bitmap
	buf[1] = byte(acc)
	buf[2] = byte(acc >> 8)
	buf[3] = byte(acc >> 16)
	buf[4] = byte(acc >> 24)
	return 1 + (i+1)/2
}

func PackUint32To(w io.Writer, buf []byte, x uint32) error {
	n := PackUint32(buf, x)
	_, err := w.Write(buf[:n])
	return err
}

// UnpackUint32 unpacks buf and returns the value.
func UnpackUint32(bitmap byte, buf []byte) uint32 {
	switch bitmap {
	case 0:
		return 0
	case 255:
		return binary.LittleEndian.Uint32(buf)
	}
	entry := unpackTable[bitmap-1]
	a, b, c, d, e, f, g := entry.A()*4, entry.B()*4, entry.C()*4, entry.D()*4, entry.E()*4, entry.F()*4, entry.G()*4
	switch entry.Num() {
	case 1:
		return uint32(buf[0]&0xF) << a
	case 2:
		return uint32(buf[0]&0xF)<<a | uint32(buf[0]>>4)<<b
	case 3:
		_ = buf[1]
		return uint32(buf[0]&0xF)<<a | uint32(buf[0]>>4)<<b | uint32(buf[1]&0xF)<<c
	case 4:
		_ = buf[1]
		return uint32(buf[0]&0xF)<<a | uint32(buf[0]>>4)<<b | uint32(buf[1]&0xF)<<c | uint32(buf[1]>>4)<<d
	case 5:
		_ = buf[2]
		return uint32(buf[0]&0xF)<<a | uint32(buf[0]>>4)<<b | uint32(buf[1]&0xF)<<c | uint32(buf[1]>>4)<<d |
			uint32(buf[2]&0xF)<<e
	case 6:
		_ = buf[2]
		return uint32(buf[0]&0xF)<<a | uint32(buf[0]>>4)<<b | uint32(buf[1]&0xF)<<c | uint32(buf[1]>>4)<<d |
			uint32(buf[2]&0xF)<<e | uint32(buf[2]>>4)<<f
	}
	_ = buf[3]
	return uint32(buf[0]&0xF)<<a | uint32(buf[0]>>4)<<b | uint32(buf[1]&0xF)<<c | uint32(buf[1]>>4)<<d |
		uint32(buf[2]&0xF)<<e | uint32(buf[2]>>4)<<f | uint32(buf[3]&0xF)<<g
}

func UnpackUint32From(r iobyte.ByteReader, buf []byte) (x uint32, err error) {
	bitmap, err := r.ReadByte()
	if err != nil {
		return
	}
	if bitmap == 0 {
		return
	}
	// 1 or 2 nibbles per byte.
	n := (bits.OnesCount8(bitmap) + 1) / 2
	if n == 1 {
		buf[0], err = r.ReadByte()
	} else {
		_, err = io.ReadFull(r, buf[:n])
	}
	if err != nil {
		return
	}
	return UnpackUint32(bitmap, buf), nil
}
