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
	// Observe that we get a '*' if a field is *either* nullable *or* optional;
	//  and we get an extra bool for the second cardinality +1'er if both are true.
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }} struct{
			{{- range $field := .Type.Fields }}
			{{ $field.Name }} {{if or $field.IsOptional $field.IsNullable }}*{{end}}{{ $field.Type | mungeTypeNodeIdent }}
			{{- end}}
			{{ range $field := .Type.Fields }}
			{{- if and $field.IsOptional $field.IsNullable }}
			{{ $field.Name }}__exists bool
			{{- end}}
			{{- end}}
		}

	`, w, gk)
}

func (gk generateKindStruct) EmitNativeAccessors(w io.Writer) {
	doTemplate(`
		{{- range $field := .Type.Fields -}}
		func (x {{ .Type | mungeTypeNodeIdent }}) Field{{ $field.Name | titlize }}() {{ $field.Type | mungeTypeNodeIdent }} {
			// TODO going to tear through here with changes to Maybe system in a moment anyway
			return {{ $field.Type | mungeTypeNodeIdent }}{}
		}
		{{end}}

	`, w, gk)
}

func (gk generateKindStruct) EmitNativeBuilder(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }}__Content struct {
			{{- range $field := .Type.Fields }}
			// TODO
			{{- end}}
		}

		func (b {{ .Type | mungeTypeNodeIdent }}__Content) Build() ({{ .Type | mungeTypeNodeIdent }}, error) {
			x := {{ .Type | mungeTypeNodeIdent }}{
				// TODO
			}
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
			Maybe typed.Maybe
			Value {{ .Type | mungeTypeNodeIdent }}
		}

		func (m Maybe{{ .Type | mungeTypeNodeIdent }}) Must() {{ .Type | mungeTypeNodeIdent }} {
			if m.Maybe != typed.Maybe_Value {
				panic("unbox of a maybe rejected")
			}
			return m.Value
		}

	`, w, gk)
}
