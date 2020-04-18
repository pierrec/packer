package packuint

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestPackUint64(t *testing.T) {
	for _, tc := range []uint64{
		0,
		1,
		127,
		128,
		256,
		1 << 10,
		1 << 20,
		1<<20 | 1<<10,
		1 << 30,
		1 << 63,
		0xF0F0F0F0,
		0xF0F0F0F0F0F0F0F0,
		^uint64(0),
		math.Float64bits(math.Pi),
		math.Float64bits(math.Phi),
		math.Float64bits(math.E),
	} {
		label := fmt.Sprintf("%d", tc)
		t.Run(label, func(t *testing.T) {
			buf := make([]byte, 16)
			n := PackUint64(buf, tc)
			t.Logf("packed size for %d = %d", tc, n)
			x := UnpackUint64(buf[0], buf[1:])
			if got, want := x, tc; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
		t.Run("io "+label, func(t *testing.T) {
			buf := make([]byte, 16)
			rw := new(bytes.Buffer)
			if err := PackUint64To(rw, buf, tc); err != nil {
				t.Fatal(err)
			}

			x, err := UnpackUint64From(rw, buf)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := x, tc; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}

func TestPackUint32(t *testing.T) {
	for _, tc := range []uint32{
		0,
		1,
		127,
		128,
		256,
		1 << 10,
		1 << 20,
		1<<20 | 1<<10,
		1 << 31,
		0xF0F0F0F0,
		^uint32(0),
	} {
		label := fmt.Sprintf("%d", tc)
		t.Run(label, func(t *testing.T) {
			buf := make([]byte, 16)
			n := PackUint32(buf, tc)
			t.Logf("packed size for %d = %d", tc, n)
			x := UnpackUint32(buf[0], buf[1:])
			if got, want := x, tc; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
		t.Run("io "+label, func(t *testing.T) {
			buf := make([]byte, 16)
			rw := new(bytes.Buffer)
			if err := PackUint32To(rw, buf, tc); err != nil {
				t.Fatal(err)
			}

			x, err := UnpackUint32From(rw, buf)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := x, tc; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}

var benchInt int
var benchX64 uint64
var benchX32 uint32

func benchmarkPackUint64(b *testing.B, x uint64) {
	buf := make([]byte, 16)
	for i := 0; i < b.N; i++ {
		benchInt = PackUint64(buf, x)
	}
}

func BenchmarkPackUint64_E(b *testing.B) {
	benchmarkPackUint64(b, math.Float64bits(math.E))
}

func BenchmarkPackUint64_1024(b *testing.B) {
	benchmarkPackUint64(b, 1024)
}

func BenchmarkPackUint64_zebra(b *testing.B) {
	benchmarkPackUint64(b, 0xF0F0F0F0)
}

func BenchmarkPackUint64_zebra2(b *testing.B) {
	benchmarkPackUint64(b, 0xF0F0F0F0F0F0F0F0)
}

func benchmarkUnpackUint64(b *testing.B, x uint64) {
	buf := make([]byte, 16)
	PackUint64(buf, x)
	bitmap := buf[0]
	buf = buf[1:]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchX64 = UnpackUint64(bitmap, buf)
	}
}

func BenchmarkUnpackUint64_E(b *testing.B) {
	benchmarkUnpackUint64(b, math.Float64bits(math.E))
}

func BenchmarkUnpackUint64_1024(b *testing.B) {
	benchmarkUnpackUint64(b, 1024)
}

func BenchmarkUnpackUint64_zebra2(b *testing.B) {
	benchmarkUnpackUint64(b, 0xF0F0)
}

func BenchmarkUnpackUint64_zebra4(b *testing.B) {
	benchmarkUnpackUint64(b, 0xF0F0F0F0)
}

func BenchmarkUnpackUint64_zebra6(b *testing.B) {
	benchmarkUnpackUint64(b, 0xF0F0F0F0F0F0F0)
}

func BenchmarkUnpackUint64_zebra8(b *testing.B) {
	benchmarkUnpackUint64(b, 0xF0F0F0F0F0F0F0F0)
}

func benchmarkPackUint32(b *testing.B, x uint32) {
	buf := make([]byte, 16)
	for i := 0; i < b.N; i++ {
		benchInt = PackUint32(buf, x)
	}
}

func BenchmarkPackUint32_E(b *testing.B) {
	benchmarkPackUint32(b, uint32(math.Float64bits(math.E)))
}

func BenchmarkPackUint32_1024(b *testing.B) {
	benchmarkPackUint32(b, 1024)
}

func BenchmarkPackUint32_zebra(b *testing.B) {
	benchmarkPackUint32(b, 0xF0F0)
}

func BenchmarkPackUint32_zebra2(b *testing.B) {
	benchmarkPackUint32(b, 0xF0F0F0F0)
}

func benchmarkUnpackUint32(b *testing.B, x uint32) {
	buf := make([]byte, 16)
	PackUint32(buf, x)
	bitmap := buf[0]
	buf = buf[1:]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchX32 = UnpackUint32(bitmap, buf)
	}
}

func BenchmarkUnpackUint32_E(b *testing.B) {
	benchmarkPackUint32(b, uint32(math.Float64bits(math.E)))
}

func BenchmarkUnpackUint32_1024(b *testing.B) {
	benchmarkUnpackUint32(b, 1024)
}

func BenchmarkUnpackUint32_zebra2(b *testing.B) {
	benchmarkUnpackUint32(b, 0xF0F0)
}

func BenchmarkUnpackUint32_zebra4(b *testing.B) {
	benchmarkUnpackUint32(b, 0xF0F0F0F0)
}
