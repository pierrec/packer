package packer

// Header is defined as follow:
//   field     bits
//   -----     ----
//   version   4
//   Flag      1
//   Len       16
//   (unused)  11
type Header uint32

// Getters.
func (x Header) version() uint { return uint(x&0xF) }
func (x Header) Flag() bool { return x>>4&1 != 0 }
func (x Header) Len() int { return int(x>>5&0xFFFF) }

// Setters.
func (x *Header) versionSet(v uint) *Header { *x = *x&^0xF | Header(v)&0xF; return x }
func (x *Header) FlagSet(v bool) *Header { const b = 1<<4; if v { *x = *x&^b | b } else { *x &^= b }; return x }
func (x *Header) LenSet(v int) *Header { *x = *x&^(0xFFFF<<5) | (Header(v)&0xFFFF<<5); return x }
