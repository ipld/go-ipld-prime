package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindStruct(t schema.Type) typeGenerator {
	return generateKindStruct{
		t.(schema.TypeStruct),
		generateKindedRejections_Map{t},
	}
}

type generateKindStruct struct {
	Type schema.TypeStruct
	generateKindedRejections_Map
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindStruct) EmitNodeType(w io.Writer) {
	// Observe that we get a '*' if a field is *either* nullable *or* optional;
	//  and we get an extra bool for the second cardinality +1'er if both are true.
	doTemplate(`
		var _ ipld.Node = {{ .Type.Name }}{}
		var _ typed.Node = typed.Node(nil) // TODO

		type {{ .Type.Name }} struct{
			{{- range $field := .Type.Fields }}
			{{ $field.Name }} {{if or $field.IsOptional $field.IsNullable }}*{{end}}{{ $field.Type.Name }}
			{{- end}}
			{{ range $field := .Type.Fields }}
			{{- if and $field.IsOptional $field.IsNullable }}
			{{ $field.Name }}__exists bool
			{{- end}}
			{{- end}}
		}

	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Map
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodTraverseField(w io.Writer) {
	doTemplate(`
		func (x {{ .Type.Name }}) TraverseField(key string) (ipld.Node, error) {
			switch key {
			{{- range $field := .Type.Fields }}
			case "{{ $field.Name }}":
				{{- if and $field.IsOptional $field.IsNullable }}
				if !x.{{ $field.Name }}__exists {
					return ipld.Undefined, nil
				}
				if x.{{ $field.Name }} == nil {
					return ipld.Nil, nil
				}
				{{- else if $field.IsOptional }}
				if x.{{ $field.Name }} == nil {
					return ipld.Undefined, nil
				}
				{{- else if $field.IsNullable }}
				if x.{{ $field.Name }} == nil {
					return ipld.Nil, nil
				}
				{{- end}}
				return x.{{ $field.Name }}, nil
			{{- end}}
			default:
				return nil, typed.ErrNoSuchField{Type: nil /*TODO*/, FieldName: key}
			}
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodMapIterator(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) MapIterator() ipld.MapIterator {
			return nil // TODO EmitNodeMethodMapIterator
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) Length() int {
			return {{ len .Type.Fields }}
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type.Name }}) NodeBuilder() ipld.NodeBuilder {
			return nil // TODO EmitNodeMethodNodeBuilder
		}
	`, w, gk)
}
