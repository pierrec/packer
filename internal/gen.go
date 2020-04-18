//+build ignore

package main

import (
	"log"
	"os"

	"github.com/pierrec/packer"
)

type UintEntry struct {
	Num                 [3]uint8
	A, B, C, D, E, F, G [3]uint8
}

func main() {
	out, err := os.Create("types.go")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := packer.GenPackedStruct(out, &packer.Config{PkgName: "internal"}, UintEntry{}); err != nil {
		log.Fatal(err)
	}
}
