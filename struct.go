package packer

import (
	"fmt"
	"io"
	"math/bits"
	"reflect"
	"text/template"
)

type (
	StructError interface {
		error
		private()
	}
	structError string
)

func (e structError) private()      {}
func (e structError) Error() string { return string(e) }

// The following errors satisfy the StructError interface.
const (
	ErrNotAStruct     structError = "not a struct"
	ErrEmptyStruct    structError = "empty struct"
	ErrEmbeddedField  structError = "embedded field not supported"
	ErrFieldNotArray  structError = "field must be an array"
	ErrFieldType      structError = "unsupported field type"
	ErrFieldOverflow  structError = "too many bits for field type"
	ErrStructOverflow structError = "struct overflows uint64"
)

// Struct packs a struct into an uint{8, 16, 32, 64} and generates the code to access its members.
// The struct must be defined as follow:
//  - field name is used as the method name to access its value
//  - field type must be an array of T, where:
//     - T is the type returned by the field method
//     - T is one of bool, {u}int or {u}int{8, 16, 32, 64}
//     - the size of the array defines the number of bits used by the value
//
// It returns an error if the struct overflows uint64.
//
// ``pkg`` defines the package name used for the generated code.
//
// Example:
//  type Header struct{
//    version [4]uint
//    Flag    [1]bool
//    Len     [16]int
//  }
// results in the following type:
//    type Header uint32
// with methods:
//       Header.version() uint
//       Header.Flag() bool
//       Header.Len() int
//    (*Header).versionSet(uint)
//    (*Header).Flag(bool)
//    (*Header).LenSet(int)
func Struct(w io.Writer, pkg string, s interface{}) error {
	werr := func(err error) error { return fmt.Errorf("packer: type %T: %w", s, err) }
	werrf := func(f string, err error) error { return fmt.Errorf("packer: type %T.%s: %w", s, f, err) }

	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return werr(ErrNotAStruct)
	}

	type _Field struct {
		TypeName string // overall type name
		Type     string // underlying type name
		Name     string // method name
		Out      string // returned type name
		Shift    int
		Mask     string
	}
	var fields []_Field
	var size int
	for i, n := 0, typ.NumField(); i < n; i++ {
		field := typ.Field(i)
		if field.Anonymous {
			return werrf(field.Name, ErrEmbeddedField)
		}
		if field.Type.Kind() != reflect.Array {
			return werrf(field.Name, ErrFieldNotArray)
		}
		out := field.Type.Elem()
		on := field.Type.Len()

		var n int
		switch kind := out.Kind(); kind {
		case reflect.Bool:
			n = 1
		case reflect.Int, reflect.Uint:
			n = bits.OnesCount(^uint(0))
		case reflect.Int8, reflect.Uint8:
			n = 8
		case reflect.Int16, reflect.Uint16:
			n = 16
		case reflect.Int32, reflect.Uint32:
			n = 32
		case reflect.Int64, reflect.Uint64:
			n = 64
		default:
			return werrf(field.Name, ErrFieldType)
		}
		// Make sure that the extracted bits fit into the returned type.
		if on > n {
			return werrf(field.Name, ErrFieldOverflow)
		}
		fields = append(fields, _Field{
			Name:  field.Name,
			Out:   out.String(),
			Shift: size,
			Mask:  fmt.Sprintf("0x%X", 1<<on-1),
		})
		size += on
	}

	switch {
	case size <= 0:
		return werr(ErrEmptyStruct)
	case size <= 8:
		size = 8
	case size <= 16:
		size = 16
	case size <= 32:
		size = 32
	case size <= 64:
		size = 64
	default:
		return werr(ErrStructOverflow)
	}
	typname := fmt.Sprintf("uint%d", size)
	for i := range fields {
		fields[i].TypeName = typ.Name()
		fields[i].Type = typname
	}

	const header = `package %s
`
	if _, err := fmt.Fprintf(w, header, pkg); err != nil {
		return werr(err)
	}

	err := structTemplate.Execute(w, struct {
		TypeName string
		Type     string
		Fields   []_Field
	}{
		typ.Name(),
		typname,
		fields,
	})
	if err != nil {
		return werr(err)
	}
	return nil
}

var structTemplate = template.Must(template.New("struct code gen").Parse(structSource))

const structSource = `
{{- define "body_get"}}
{{- if eq .Out "bool" -}} return x>>{{.Shift}}&1 != 0
{{- else if eq .Shift 0 -}} return {{.Out}}(x&{{.Mask}})
{{- else -}} return {{.Out}}(x>>{{.Shift}}&{{.Mask}}) {{- end}}
{{- end}}
{{- define "body_set"}}
{{- if eq .Out "bool" -}} const b = 1<<{{.Shift}}; if v { *x = *x&^b | b } else { *x &^= b }; return x
{{- else if eq .Shift 0 -}} *x = *x&^{{.Mask}} | {{.TypeName}}(v)&{{.Mask}}; return x
{{- else -}} *x = *x&^({{.Mask}}<<{{.Shift}}) | ({{.TypeName}}(v)&{{.Mask}}<<{{.Shift}}); return x {{- end}}
{{- end}}
type {{.TypeName}} {{.Type}}

// Getters.
{{range .Fields}}func (x {{.TypeName}}) {{.Name}}() {{.Out}} { {{template "body_get" .}} }
{{end}}
// Setters.
{{range .Fields}}func (x *{{.TypeName}}) {{.Name}}Set(v {{.Out}}) *{{.TypeName}} { {{template "body_set" .}} }
{{end}}`
