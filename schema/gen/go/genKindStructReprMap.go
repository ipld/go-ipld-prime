package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func getStructRepresentationMapNodeGen(t schema.TypeStruct) nodeGenerator {
	return generateStructReprMapNode{
		t,
		generateKindedRejections_Map{
			mungeTypeReprNodeIdent(t),
			string(t.Name()) + ".Representation",
		},
	}
}

type generateStructReprMapNode struct {
	Type schema.TypeStruct
	generateKindedRejections_Map
}

func (gk generateStructReprMapNode) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeReprNodeIdent }}{}

		type {{ .Type | mungeTypeReprNodeIdent }} struct{
			n *{{ .Type | mungeTypeNodeIdent }}
		}

	`, w, gk)
}

func (gk generateStructReprMapNode) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeReprNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_Map
		}
	`, w, gk)
}

func (gk generateStructReprMapNode) EmitNodeMethodLookupString(w io.Writer) {
	// almost idential to the type-level one, just with different strings in the switch.
	// TODO : support for implicits is missing.
	doTemplate(`
		func (rn {{ .Type | mungeTypeReprNodeIdent }}) LookupString(key string) (ipld.Node, error) {
			switch key {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $field := .Type.Fields }}
			case "{{ $field | $type.RepresentationStrategy.GetFieldKey }}":
				{{- if and $field.IsOptional $field.IsNullable }}
				if !rn.n.{{ $field.Name }}__exists {
					return ipld.Undef, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
				}
				if rn.n.{{ $field.Name }} == nil {
					return ipld.Null, nil
				}
				{{- else if $field.IsOptional }}
				if rn.n.{{ $field.Name }} == nil {
					return ipld.Undef, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
				}
				{{- else if $field.IsNullable }}
				if rn.n.{{ $field.Name }} == nil {
					return ipld.Null, nil
				}
				{{- end}}
				{{- if or $field.IsOptional $field.IsNullable }}
				return *rn.n.{{ $field.Name }}, nil
				{{- else}}
				return rn.n.{{ $field.Name }}, nil
				{{- end}}
			{{- end}}
			default:
				return nil, typed.ErrNoSuchField{Type: nil /*TODO*/, FieldName: key}
			}
		}
	`, w, gk)
}

func (gk generateStructReprMapNode) EmitNodeMethodLookup(w io.Writer) {
	doTemplate(`
		func (rn {{ .Type | mungeTypeReprNodeIdent }}) Lookup(key ipld.Node) (ipld.Node, error) {
			ks, err := key.AsString()
			if err != nil {
				return nil, ipld.ErrInvalidKey{"got " + key.ReprKind().String() + ", need string"}
			}
			return rn.LookupString(ks)
		}
	`, w, gk)
}

func (gk generateStructReprMapNode) EmitNodeMethodMapIterator(w io.Writer) {
	// Amusing note, the iterator ends up with a loop in its body, even though
	//  it only yields one entry pair at a time -- this is needed so we can
	//   use 'continue' statements to skip past optionals which are undefined.
	// TODO : support for implicits is missing.
	doTemplate(`
		func (rn {{ .Type | mungeTypeReprNodeIdent }}) MapIterator() ipld.MapIterator {
			return &{{ .Type | mungeTypeReprNodeItrIdent }}{rn.n, 0}
		}

		type {{ .Type | mungeTypeReprNodeItrIdent }} struct {
			node *{{ .Type | mungeTypeNodeIdent }}
			idx  int
		}

		func (itr *{{ .Type | mungeTypeReprNodeItrIdent }}) Next() (k ipld.Node, v ipld.Node, _ error) {
			if itr.idx >= {{ len .Type.Fields }} {
				return nil, nil, ipld.ErrIteratorOverread{}
			}
			for {
				switch itr.idx {
				{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
				{{- range $i, $field := .Type.Fields }}
				case {{ $i }}:
					k = String{"{{ $field | $type.RepresentationStrategy.GetFieldKey }}"}
					{{- if and $field.IsOptional $field.IsNullable }}
					if !itr.node.{{ $field.Name }}__exists {
						itr.idx++
						continue
					}
					if itr.node.{{ $field.Name }} == nil {
						v = ipld.Null
						break
					}
					{{- else if $field.IsOptional }}
					if itr.node.{{ $field.Name }} == nil {
						itr.idx++
						continue
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
			}
			itr.idx++
			return
		}
		func (itr *{{ .Type | mungeTypeReprNodeItrIdent }}) Done() bool {
			return itr.idx >= {{ len .Type.Fields }}
		}

	`, w, gk)

}

func (gk generateStructReprMapNode) EmitNodeMethodLength(w io.Writer) {
	// This is fun: it has to count down for any unset optional fields.
	// TODO : support for implicits is missing.
	doTemplate(`
		func (rn {{ .Type | mungeTypeReprNodeIdent }}) Length() int {
			l := {{ len .Type.Fields }}
			{{- range $field := .Type.Fields }}
			{{- if and $field.IsOptional $field.IsNullable }}
			if !rn.n.{{ $field.Name }}__exists {
				l--
			}
			{{- else if $field.IsOptional }}
			if rn.n.{{ $field.Name }} == nil {
				l--
			}
			{{- end}}
			{{- end}}
			return l
		}
	`, w, gk)
}

func (gk generateStructReprMapNode) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeReprNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeReprNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateStructReprMapNode) GetNodeBuilderGen() nodebuilderGenerator {
	return generateStructReprMapNb{
		gk.Type,
		genKindedNbRejections_Map{
			mungeTypeReprNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Representation.Builder",
		},
	}
}

type generateStructReprMapNb struct {
	Type schema.TypeStruct
	genKindedNbRejections_Map
}

func (gk generateStructReprMapNb) EmitNodebuilderType(w io.Writer) {
	// Note there's no need to put the reprKind in the name of the type
	//  we generate here: there's only one representation per type.
	//   (We *could* munge the reprkind in for debug symbol reading,
	//    but at present it hasn't seemed warranted.)
	doTemplate(`
		type {{ .Type | mungeTypeReprNodebuilderIdent }} struct{}

	`, w, gk)
}

func (gk generateStructReprMapNb) EmitNodebuilderConstructor(w io.Writer) {
	doTemplate(`
		func {{ .Type | mungeReprNodebuilderConstructorIdent }}() ipld.NodeBuilder {
			return {{ .Type | mungeTypeReprNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateStructReprMapNb) EmitNodebuilderMethodCreateMap(w io.Writer) {
	// Much of these looks the same as the type-level builders.  Key differences:
	//  - We interact with the rename directives here.
	//  - The "__isset" bools are generated for *all* fields -- we need these
	//     to check if a key is repeated, so we can reject that.
	//      Worth mentioning: we could also choose *not* to check this, instead
	//       insisting it's a codec layer concern.  This needs revisiting;
	//        at present I'm choosing "defense in depth", because trying to
	//         reason out the perf and usability implications in advance has
	//          yielded a huge matrix of concerns and no single clear gradient.
	// TODO : support for implicits is missing.
	doTemplate(`
		func (nb {{ .Type | mungeTypeReprNodebuilderIdent }}) CreateMap() (ipld.MapBuilder, error) {
			return &{{ .Type | mungeTypeReprNodeMapBuilderIdent }}{v:&{{ .Type | mungeTypeNodeIdent }}{}}, nil
		}

		type {{ .Type | mungeTypeReprNodeMapBuilderIdent }} struct{
			v *{{ .Type | mungeTypeNodeIdent }}
			{{- range $field := .Type.Fields }}
			{{ $field.Name }}__isset bool
			{{- end}}
		}

		func (mb *{{ .Type | mungeTypeReprNodeMapBuilderIdent }}) Insert(k, v ipld.Node) error {
			ks, err := k.AsString()
			if err != nil {
				return ipld.ErrInvalidKey{"not a string: " + err.Error()}
			}
			switch ks {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $field := .Type.Fields }}
			case "{{ $field | $type.RepresentationStrategy.GetFieldKey }}":
				if mb.{{ $field.Name }}__isset {
					panic("repeated assignment to field") // FIXME need an error type for this
				}
				{{- if $field.IsNullable }}
				if v.IsNull() {
					mb.v.{{ $field.Name }} = nil
					{{- if $field.IsOptional }}
					mb.v.{{ $field.Name }}__exists = true
					{{- else}}
					mb.{{ $field.Name }}__isset = true
					{{- end}}
					return nil
				}
				{{- else}}
				if v.IsNull() {
					panic("type mismatch on struct field assignment: cannot assign null to non-nullable field") // FIXME need an error type for this
				}
				{{- end}}
				tv, ok := v.(typed.Node)
				if !ok {
					panic("need typed.Node for insertion into struct") // FIXME need an error type for this
				}
				x, ok := v.({{ $field.Type | mungeTypeNodeIdent }})
				if !ok {
					panic("field '{{$field.Name}}' (key: '{{ $field | $type.RepresentationStrategy.GetFieldKey }}') in type {{$type.Name}} is type {{$field.Type.Name}}; cannot assign "+tv.Type().Name()) // FIXME need an error type for this
				}

				{{- if or $field.IsOptional $field.IsNullable }}
				mb.v.{{ $field.Name }} = &x
				{{- else}}
				mb.v.{{ $field.Name }} = x
				{{- end}}
				{{- if and $field.IsOptional $field.IsNullable }}
				mb.v.{{ $field.Name }}__exists = true
				{{- end}}
				mb.{{ $field.Name }}__isset = true
			{{- end}}
			default:
				return typed.ErrNoSuchField{Type: nil /*TODO:typelit*/, FieldName: ks}
			}
			return nil
		}
		func (mb *{{ .Type | mungeTypeReprNodeMapBuilderIdent }}) Delete(k ipld.Node) error {
			panic("TODO later")
		}
		func (mb *{{ .Type | mungeTypeReprNodeMapBuilderIdent }}) Build() (ipld.Node, error) {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $field := .Type.Fields }}
			{{- if not $field.IsOptional }}
			if !mb.{{ $field.Name }}__isset {
				panic("missing required field '{{$field.Name}}' (key: '{{ $field | $type.RepresentationStrategy.GetFieldKey }}') in building struct {{ $type.Name }}") // FIXME need an error type for this
			}
			{{- end}}
			{{- end}}
			v := mb.v
			mb = nil
			return v, nil
		}

	`, w, gk)
}

func (gk generateStructReprMapNb) EmitNodebuilderMethodAmendMap(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type | mungeTypeReprNodebuilderIdent }}) AmendMap() (ipld.MapBuilder, error) {
			panic("TODO later")
		}
	`, w, gk)
}
