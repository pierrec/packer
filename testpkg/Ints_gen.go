package testpkg

// Ints is defined as follow:
//   field     bits
//   -----     ----
//   Int8      8
//   Int16     16
//   Int32     32
//   (unused)  8
type Ints uint64

// Getters.
func (x Ints) Int8() int8   { return int8(x & 0xFF) }
func (x Ints) Int16() int16 { return int16(x >> 8 & 0xFFFF) }
func (x Ints) Int32() int32 { return int32(x >> 24 & 0xFFFFFFFF) }

// Setters.
func (x *Ints) Int8Set(v int8) *Ints   { *x = *x&^0xFF | Ints(v)&0xFF; return x }
func (x *Ints) Int16Set(v int16) *Ints { *x = *x&^(0xFFFF<<8) | (Ints(v) & 0xFFFF << 8); return x }
func (x *Ints) Int32Set(v int32) *Ints {
	*x = *x&^(0xFFFFFFFF<<24) | (Ints(v) & 0xFFFFFFFF << 24)
	return x
}
