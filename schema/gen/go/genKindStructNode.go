package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

// --- type-semantics node interface satisfaction --->

func (gk generateKindStruct) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

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
				{{- if $field.IsOptional }}
				if x.d.{{ $field.Name | titlize }}.Maybe == typed.Maybe_Absent {
					return ipld.Undef, nil
				}
				{{- end}}
				{{- if $field.IsNullable }}
				if x.d.{{ $field.Name | titlize }}.Maybe == typed.Maybe_Null {
					return ipld.Null, nil
				}
				{{- end}}
				{{- if or $field.IsOptional $field.IsNullable }}
				return x.d.{{ $field.Name | titlize}}.Value, nil
				{{- else}}
				return x.d.{{ $field.Name | titlize}}, nil
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
	// Note that the typed iterator will report absent fields.
	//  The representation iterator, if maplike, will not.
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
				{{- if $field.IsOptional }}
				if itr.node.d.{{ $field.Name | titlize  }}.Maybe == typed.Maybe_Absent {
					v = ipld.Undef
					break
				}
				{{- end}}
				{{- if $field.IsNullable }}
				if itr.node.d.{{ $field.Name | titlize }}.Maybe == typed.Maybe_Null {
					v = ipld.Null
					break
				}
				{{- end}}
				{{- if or $field.IsOptional $field.IsNullable }}
				v = itr.node.d.{{ $field.Name | titlize}}.Value
				{{- else}}
				v = itr.node.d.{{ $field.Name | titlize}}
				{{- end}}
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

// --- type-semantics nodebuilder --->

func (gk generateKindStruct) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateKindStruct) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindStruct{
		gk.Type,
		genKindedNbRejections_Map{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
}

type generateNbKindStruct struct {
	Type schema.TypeStruct
	genKindedNbRejections_Map
}

func (gk generateNbKindStruct) EmitNodebuilderType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodebuilderIdent }} struct{}

	`, w, gk)
}

func (gk generateNbKindStruct) EmitNodebuilderConstructor(w io.Writer) {
	doTemplate(`
		func {{ .Type | mungeNodebuilderConstructorIdent }}() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateNbKindStruct) EmitNodebuilderMethodCreateMap(w io.Writer) {
	// Some interesting edge cases to note:
	//  - This builder, being all about semantics and not at all about serialization,
	//      is order-insensitive.
	//  - We don't specially handle being given 'undef' as a value.
	//      It just falls into the "need a typed.Node" error bucket.
	//  - We only accept *codegenerated values* -- a typed.Node created
	//      in the same schema universe *isn't accepted*.
	//        REVIEW: We could try to accept those, but it might have perf/sloc costs,
	//          and it's hard to imagine a user story that gets here.
	//  - The has-been-set-if-required validation is fun; it only requires state
	//     for non-optional fields, and that often gets a little hard to follow
	//       because it gets wedged in with other logic tables around optionality.
	// REVIEW: 'x, ok := v.({{ $field.Type.Name }})' might need some stars in it... sometimes.
	// FIXME : since the typed iterator yields undefineds, this should probably also accept them...?
	// TODO : review the panic of `ErrNoSuchField` in `BuilderForValue` --
	//  see the comments in the NodeBuilder interface for the open questions on this topic.
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) CreateMap() (ipld.MapBuilder, error) {
			mb := &{{ .Type | mungeTypeNodeMapBuilderIdent }}{v:&{{ .Type | mungeTypeNodeIdent }}{}}
			{{- range $field := .Type.Fields }}
			{{- if $field.IsOptional }}
			mb.v.d.{{ $field.Name | titlize }}.Maybe = typed.Maybe_Absent
			{{- end}}
			{{- end}}
			return mb, nil
		}

		type {{ .Type | mungeTypeNodeMapBuilderIdent }} struct{
			v *{{ .Type | mungeTypeNodeIdent }}
			{{- range $field := .Type.Fields }}
			{{- if not $field.IsOptional }}
			{{ $field.Name }}__isset bool
			{{- end}}
			{{- end}}
		}

		func (mb *{{ .Type | mungeTypeNodeMapBuilderIdent }}) Insert(k, v ipld.Node) error {
			ks, err := k.AsString()
			if err != nil {
				return ipld.ErrInvalidKey{"not a string: " + err.Error()}
			}
			switch ks {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $field := .Type.Fields }}
			case "{{ $field.Name }}":
				{{- if $field.IsNullable }}
				if v.IsNull() {
					mb.v.d.{{ $field.Name | titlize}}.Maybe = typed.Maybe_Null
					{{- if not $field.IsOptional }}
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
					panic("field '{{$field.Name}}' in type {{$type.Name}} is type {{$field.Type.Name}}; cannot assign "+tv.Type().Name()) // FIXME need an error type for this
				}

				{{- if or $field.IsOptional $field.IsNullable }}
				mb.v.d.{{ $field.Name | titlize}}.Value = x
				{{- else}}
				mb.v.d.{{ $field.Name | titlize}} = x
				{{- end}}
				{{- if $field.IsOptional }}
				mb.v.d.{{ $field.Name | titlize}}.Maybe = typed.Maybe_Value
				{{- else}}
				mb.{{ $field.Name }}__isset = true
				{{- end}}
			{{- end}}
			default:
				return typed.ErrNoSuchField{Type: nil /*TODO:typelit*/, FieldName: ks}
			}
			return nil
		}
		func (mb *{{ .Type | mungeTypeNodeMapBuilderIdent }}) Delete(k ipld.Node) error {
			panic("TODO later")
		}
		func (mb *{{ .Type | mungeTypeNodeMapBuilderIdent }}) Build() (ipld.Node, error) {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $field := .Type.Fields }}
			{{- if not $field.IsOptional }}
			if !mb.{{ $field.Name }}__isset {
				panic("missing required field '{{$field.Name}}' in building struct {{ $type.Name }}") // FIXME need an error type for this
			}
			{{- end}}
			{{- end}}
			v := mb.v
			mb = nil
			return v, nil
		}
		func (mb *{{ .Type | mungeTypeNodeMapBuilderIdent }}) BuilderForKeys() ipld.NodeBuilder {
			return _String__NodeBuilder{}
		}
		func (mb *{{ .Type | mungeTypeNodeMapBuilderIdent }}) BuilderForValue(ks string) ipld.NodeBuilder {
			switch ks {
			{{- range $field := .Type.Fields }}
			case "{{ $field.Name }}":
				return {{ $field.Type | mungeNodebuilderConstructorIdent }}()
			{{- end}}
			default:
				panic(typed.ErrNoSuchField{Type: nil /*TODO:typelit*/, FieldName: ks})
			}
			return nil
		}

	`, w, gk)
}

func (gk generateNbKindStruct) EmitNodebuilderMethodAmendMap(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) AmendMap() (ipld.MapBuilder, error) {
			panic("TODO later")
		}
	`, w, gk)
}

// --- entrypoints to representation --->

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
