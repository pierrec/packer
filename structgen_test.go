package packer

type Header uint32

// Getters.
func (x Header) version() uint { return uint(x&15) }
func (x Header) Flag() bool { return x>>4&1 != 0 }
func (x Header) Len() int { return int(x>>5&65535) }

// Setters.
func (x *Header) versionSet(v uint) *Header { *x = *x&^15 | Header(v)&15; return x }
func (x *Header) FlagSet(v bool) *Header { const b = 1<<4; if v { *x = *x&^b | b } else { *x &^= b }; return x }
func (x *Header) LenSet(v int) *Header { *x = *x&^(65535<<5) | (Header(v)&65535<<5); return x }
