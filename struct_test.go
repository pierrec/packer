package packer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

//go:generate go run gen.go

func TestStruct(t *testing.T) {
	type tcase struct {
		in  interface{}
		err error
	}

	type (
		Ints struct {
			Int8  int8
			Int16 int16
			Int32 int32
		}
		Uints struct {
			Uint8  uint8
			Uint16 uint16
			Uint32 uint32
		}
		Version1 struct {
			version [4]uint
			flag    [1]bool
		}
		Version2 struct {
			version [4]uint
			flag    bool
			Len     [16]int
		}
		Version3 struct {
			version  [4]uint
			flag     bool
			_        [7]int // reserved
			Len      [16]int
			_        [4]int // reserved
			Checksum [32]uint32
		}
		Broken1 struct {
			X [64]int64
			Y [64]int64
		}
		Broken2 struct {
			Data [32]int8
		}
		Broken3 struct {
			Data [2]float32
		}
		Broken4 struct {
			Data float32
		}
		Broken5 struct {
			Version1
		}
		Broken6 struct{}
	)

	for _, tc := range []tcase{
		{Ints{}, nil},
		{Uints{}, nil},
		{Version1{}, nil},
		{Version2{}, nil},
		{Version3{}, nil},
		{0, ErrNotAStruct},
		{Broken1{}, ErrStructOverflow},
		{Broken2{}, ErrFieldOverflow},
		{Broken3{}, ErrFieldType},
		{Broken4{}, ErrFieldBadType},
		{Broken5{}, ErrEmbeddedField},
		{Broken6{}, ErrEmptyStruct},
	} {
		label := fmt.Sprintf("testpkg/%s_gen.go", reflect.TypeOf(tc.in).Name())
		t.Run(label, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := GenPackedStruct(buf, &Config{PkgName: "testpkg"}, tc.in)
			switch {
			case tc.err == nil && err != nil:
				t.Fatal(err)
			case tc.err != nil && err == nil:
				t.Fatal("expected error not found")
			case tc.err != nil && err != nil:
				var serr Error
				switch {
				case !errors.As(err, &serr):
					t.Fatalf("got %T; want %T", err, serr)
				case !errors.Is(err, tc.err):
					t.Fatalf("got %q; want %q", serr, tc.err.Error())
				}
			}

			if tc.err != nil || err != nil {
				// Do not create invalid files.
				return
			}

			out, err := os.Create(label)
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()
			_, err = io.Copy(out, buf)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func ExampleStruct() {
	// type Header struct {
	//   version [4]uint
	//   Flag    bool
	//   Len     [8]int
	// }

	var h Header

	h.versionSet(15).FlagSet(true).LenSet(1000)

	fmt.Printf("version = %d\n", h.version())
	fmt.Printf("flag = %v\n", h.Flag())
	fmt.Printf("len = %d\n", h.Len())
	// Output:
	// version = 15
	// flag = true
	// len = 1000
}
