package packer

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
	"text/template"
)

type (
	// Error defines the interface satisfied by all errors generated by this package.
	Error interface {
		error
		// Make sure errors cannot be created outside this package.
		private()
	}
	structError string
)

func (e structError) private()      {}
func (e structError) Error() string { return string(e) }

// The package functions wrap one of the following errors.
const (
	ErrNotAStruct     structError = "not a struct"
	ErrEmptyStruct    structError = "empty struct"
	ErrEmbeddedField  structError = "embedded field not supported"
	ErrFieldBadType   structError = "field must be one of array, bool, uint{8,16,32}"
	ErrFieldType      structError = "unsupported field type"
	ErrFieldOverflow  structError = "too many bits for field type"
	ErrStructOverflow structError = "struct overflows uint64"
)

// Struct packs a struct into an uint{8, 16, 32, 64} and generates the code to access its members.
// The struct must be defined as follow:
//  - field name is used as the method name to access its value
//  - fields named _ do not produce any method
//  - field type must be either:
//     - bool
//     - uint{8, 16, 32}
//     - [n]T where:
//       - T is the type returned by the field method
//       - T is one of bool, {u}int or {u}int{8, 16, 32, 64}
//       - n defines the number of bits used by the value
//
// It returns an error if the struct overflows uint64.
//
// ``pkg`` defines the package name used for the generated code. If empty, the package clause is not generated.
//
// Example:
//  type Header struct{
//    version [4]uint
//    Flag    bool
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
func Struct(w io.Writer, config *Config, s interface{}) error {
	werr := func(err error) error { return fmt.Errorf("packer: type %T: %w", s, err) }
	werrf := func(f string, err error) error { return fmt.Errorf("packer: type %T.%s: %w", s, f, err) }

	config.Init()

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
	var descr []struct {
		string
		int
	}

	var size int
	for i, n := 0, typ.NumField(); i < n; i++ {
		field := typ.Field(i)
		if field.Anonymous {
			return werrf(field.Name, ErrEmbeddedField)
		}

		out := field.Type
		var outBits int
		switch field.Type.Kind() {
		case reflect.Bool:
			outBits = 1
		case reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Uint8, reflect.Uint16, reflect.Uint32:
			outBits = out.Bits()
		case reflect.Array:
			out = field.Type.Elem()
			outBits = field.Type.Len()
			var n int
			switch out.Kind() {
			case reflect.Bool:
				n = 1
			case reflect.Int, reflect.Uint:
				// Code generated on 64bits platforms must work on 32bits ones.
				n = 32
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				n = out.Bits()
			default:
				return werrf(field.Name, ErrFieldType)
			}
			// Make sure that the extracted bits fit into the returned type.
			if outBits > n {
				return werrf(field.Name, ErrFieldOverflow)
			}
		default:
			return werrf(field.Name, ErrFieldBadType)
		}

		fields = append(fields, _Field{
			Name:  field.Name,
			Out:   out.String(),
			Shift: size,
			Mask:  fmt.Sprintf("0x%X", 1<<outBits-1),
		})
		size += outBits
		descr = append(descr, struct {
			string
			int
		}{field.Name, outBits})
	}

	var unused int
	switch {
	case size <= 0:
		return werr(ErrEmptyStruct)
	case size <= 8:
		unused += 8 - size
		size = 8
	case size <= 16:
		unused += 16 - size
		size = 16
	case size <= 32:
		unused += 32 - size
		size = 32
	case size <= 64:
		unused += 64 - size
		size = 64
	default:
		return werr(ErrStructOverflow)
	}
	typname := fmt.Sprintf("uint%d", size)
	for i := range fields {
		fields[i].TypeName = typ.Name()
		fields[i].Type = typname
	}

	// Package header.
	header := []string{config.TopComments}
	if config.PkgName != "" {
		line := fmt.Sprintf("package %s\n", config.PkgName)
		header = append(header, line)
	}
	if _, err := fmt.Fprintf(w, strings.Join(header, "\n")); err != nil {
		return werr(err)
	}

	// Type comments.
	buf := new(strings.Builder)
	tw := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	_, _ = fmt.Fprintf(tw, "//   field\t\tbits\n")
	_, _ = fmt.Fprintf(tw, "//   -----\t\t----\n")
	for _, c := range descr {
		_, _ = fmt.Fprintf(tw, "//   %s\t\t%d\n", c.string, c.int)
	}
	if unused > 0 {
		_, _ = fmt.Fprintf(tw, "//   (unused)\t\t%d\n", unused)
	}
	_ = tw.Flush()
	comments := fmt.Sprintf("// %s is defined as follow:\n%s", typ.Name(), buf.String())

	err := structTemplate.Execute(w, struct {
		Comments string
		TypeName string
		Type     string
		Fields   []_Field
	}{
		comments,
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
{{- if and (eq .Out "bool") (eq .Shift 0) -}} return x&1 != 0
{{- else if eq .Out "bool" -}} return x>>{{.Shift}}&1 != 0
{{- else if eq .Shift 0 -}} return {{.Out}}(x&{{.Mask}})
{{- else -}} return {{.Out}}(x>>{{.Shift}}&{{.Mask}}) {{- end}}
{{- end}}
{{- define "body_set"}}
{{- if and (eq .Out "bool") (eq .Shift 0) -}} if v { *x |= 1 } else { *x &^= 1 }; return x
{{- else if eq .Out "bool" -}} const b = 1<<{{.Shift}}; if v { *x = *x&^b | b } else { *x &^= b }; return x
{{- else if eq .Shift 0 -}} *x = *x&^{{.Mask}} | {{.TypeName}}(v)&{{.Mask}}; return x
{{- else -}} *x = *x&^({{.Mask}}<<{{.Shift}}) | ({{.TypeName}}(v)&{{.Mask}}<<{{.Shift}}); return x {{- end}}
{{- end}}
{{.Comments -}}
type {{.TypeName}} {{.Type}}

// Getters.
{{range .Fields}}
{{- if not (eq .Name "_") -}}
func (x {{.TypeName}}) {{.Name}}() {{.Out}} { {{template "body_get" .}} }
{{ end -}}
{{end}}
// Setters.
{{range .Fields}}
{{- if not (eq .Name "_") -}}
func (x *{{.TypeName}}) {{.Name}}Set(v {{.Out}}) *{{.TypeName}} { {{template "body_set" .}} }
{{ end -}}
{{end}}`
