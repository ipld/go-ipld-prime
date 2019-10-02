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

func (gk generateKindStruct) EmitNodeType(w io.Writer) {
	// Observe that we get a '*' if a field is *either* nullable *or* optional;
	//  and we get an extra bool for the second cardinality +1'er if both are true.
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

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

func (gk generateKindStruct) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Map
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodLookupString(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) LookupString(key string) (ipld.Node, error) {
			switch key {
			{{- range $field := .Type.Fields }}
			case "{{ $field.Name }}":
				{{- if and $field.IsOptional $field.IsNullable }}
				if !x.{{ $field.Name }}__exists {
					return ipld.Undef, nil
				}
				if x.{{ $field.Name }} == nil {
					return ipld.Null, nil
				}
				{{- else if $field.IsOptional }}
				if x.{{ $field.Name }} == nil {
					return ipld.Undef, nil
				}
				{{- else if $field.IsNullable }}
				if x.{{ $field.Name }} == nil {
					return ipld.Null, nil
				}
				{{- end}}
				{{- if or $field.IsOptional $field.IsNullable }}
				return *x.{{ $field.Name }}, nil
				{{- else}}
				return x.{{ $field.Name }}, nil
				{{- end}}
			{{- end}}
			default:
				return nil, typed.ErrNoSuchField{Type: nil /*TODO*/, FieldName: key}
			}
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodLookup(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) Lookup(key ipld.Node) (ipld.Node, error) {
			ks, err := key.AsString()
			if err != nil {
				return nil, ipld.ErrInvalidKey{"got " + key.ReprKind().String() + ", need string"}
			}
			return x.LookupString(ks)
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodMapIterator(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) MapIterator() ipld.MapIterator {
			return &{{ .Type | mungeTypeNodeItrIdent }}{&x, 0}
		}

		type {{ .Type | mungeTypeNodeItrIdent }} struct {
			node *{{ .Type | mungeTypeNodeIdent }}
			idx  int
		}

		func (itr *{{ .Type | mungeTypeNodeItrIdent }}) Next() (k ipld.Node, v ipld.Node, _ error) {
			if itr.idx >= {{ len .Type.Fields }} {
				return nil, nil, ipld.ErrIteratorOverread{}
			}
			switch itr.idx {
			{{- range $i, $field := .Type.Fields }}
			case {{ $i }}:
				k = String{"{{ $field.Name }}"}
				{{- if and $field.IsOptional $field.IsNullable }}
				if !itr.node.{{ $field.Name }}__exists {
					v = ipld.Undef
					break
				}
				if itr.node.{{ $field.Name }} == nil {
					v = ipld.Null
					break
				}
				{{- else if $field.IsOptional }}
				if itr.node.{{ $field.Name }} == nil {
					v = ipld.Undef
					break
				}
				{{- else if $field.IsNullable }}
				if itr.node.{{ $field.Name }} == nil {
					v = ipld.Null
					break
				}
				{{- end}}
				v = itr.node.{{ $field.Name }}
			{{- end}}
			default:
				panic("unreachable")
			}
			itr.idx++
			return
		}
		func (itr *{{ .Type | mungeTypeNodeItrIdent }}) Done() bool {
			return itr.idx >= {{ len .Type.Fields }}
		}

	`, w, gk)
}

func (gk generateKindStruct) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Length() int {
			return {{ len .Type.Fields }}
		}
	`, w, gk)
}

func (gk generateKindStruct) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | mungeTypeNodeIdent }}) Representation() ipld.Node {
			return {{ .Type | mungeTypeReprNodeIdent }}{&n}
		}
	`, w, gk)
}

func (gk generateKindStruct) GetRepresentationNodeGen() nodeGenerator {
	switch gk.Type.RepresentationStrategy().(type) {
	case schema.StructRepresentation_Map:
		return getStructRepresentationMapNodeGen(gk.Type)
	default:
		panic("missing case in switch for repr strategy for structs")
	}
}
