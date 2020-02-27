package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindStruct(t schema.Type) typedNodeGenerator {
	return generateKindStruct{
		t.(schema.TypeStruct),
		generateKindedRejections_Map{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindStruct struct {
	Type schema.TypeStruct
	generateKindedRejections_Map
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindStruct) EmitNativeType(w io.Writer) {
	// The data is actually the content type, just embedded in an unexported field,
	//  which means we get immutability, plus initializing the object is essentially a memmove.
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }} struct{
			d {{ .Type | mungeTypeNodeIdent }}__Content
		}

	`, w, gk)
}

func (gk generateKindStruct) EmitNativeAccessors(w io.Writer) {
	doTemplate(`
		{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
		{{- range $field := .Type.Fields -}}
		func (x {{ $type | mungeTypeNodeIdent }}) Field{{ $field.Name | titlize }}()
			{{- if or $field.IsOptional $field.IsNullable }}Maybe{{end}}{{ $field.Type | mungeTypeNodeIdent }} {
			return x.d.{{ $field.Name | titlize }}
		}
		{{end}}

	`, w, gk)
}

func (gk generateKindStruct) EmitNativeBuilder(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }}__Content struct {
			{{- range $field := .Type.Fields}}
			{{ $field.Name | titlize }} {{if or $field.IsOptional $field.IsNullable }}Maybe{{end}}{{ $field.Type | mungeTypeNodeIdent }}
			{{- end}}
		}

		func (b {{ .Type | mungeTypeNodeIdent }}__Content) Build() ({{ .Type | mungeTypeNodeIdent }}, error) {
			{{- range $field := .Type.Fields -}}
			{{- if or $field.IsOptional $field.IsNullable }}
			{{- /* if both modifiers present, anything goes */ -}}
			{{- else if $field.IsOptional }}
			if b.{{ $field.Name | titlize }}.Maybe == schema.Maybe_Null {
				return {{ $field.Type | mungeTypeNodeIdent }}{}, fmt.Errorf("cannot be absent")
			}
			{{- else if $field.IsNullable }}
			if b.{{ $field.Name | titlize }}.Maybe == schema.Maybe_Absent {
				return {{ $field.Type | mungeTypeNodeIdent }}{}, fmt.Errorf("cannot be null")
			}
			{{- end}}
			{{- end}}
			x := {{ .Type | mungeTypeNodeIdent }}{b}
			// FUTURE : want to support customizable validation.
			//   but 'if v, ok := x.(schema.Validatable); ok {' doesn't fly: need a way to work on concrete types.
			return x, nil
		}
		func (b {{ .Type | mungeTypeNodeIdent }}__Content) MustBuild() {{ .Type | mungeTypeNodeIdent }} {
			if x, err := b.Build(); err != nil {
				panic(err)
			} else {
				return x
			}
		}

	`, w, gk)
}

func (gk generateKindStruct) EmitNativeMaybe(w io.Writer) {
	doTemplate(`
		type Maybe{{ .Type | mungeTypeNodeIdent }} struct {
			Maybe schema.Maybe
			Value {{ .Type | mungeTypeNodeIdent }}
		}

		func (m Maybe{{ .Type | mungeTypeNodeIdent }}) Must() {{ .Type | mungeTypeNodeIdent }} {
			if m.Maybe != schema.Maybe_Value {
				panic("unbox of a maybe rejected")
			}
			return m.Value
		}

	`, w, gk)
}
