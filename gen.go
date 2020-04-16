//+build ignore

package main

import (
	"log"
	"os"

	"github.com/pierrec/packer"
)

type Header struct {
	version [4]uint
	Flag    bool
	Len     [16]int
}

func main() {
	out, err := os.Create("structgen_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	err = packer.Struct(out, "packer", Header{})
	if err != nil {
		log.Fatal(err)
	}
}
