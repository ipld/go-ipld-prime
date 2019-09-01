package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func (gk generateKindStruct) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindStruct{
		gk.Type,
		genKindedNbRejections_Map{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
}

func (gk generateKindStruct) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
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
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) CreateMap() (ipld.MapBuilder, error) {
			return &{{ .Type | mungeTypeNodeMapBuilderIdent }}{v:&{{ .Type | mungeTypeNodeIdent }}{}}, nil
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
					panic("field '{{$field.Name}}' in type {{$type.Name}} is type {{$field.Type.Name}}; cannot assign "+tv.Type().Name()) // FIXME need an error type for this
				}

				{{- if or $field.IsOptional $field.IsNullable }}
				mb.v.{{ $field.Name }} = &x
				{{- else}}
				mb.v.{{ $field.Name }} = x
				{{- end}}
				{{- if and $field.IsOptional $field.IsNullable }}
				mb.v.{{ $field.Name }}__exists = true
				{{- else if not $field.IsOptional }}
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

	`, w, gk)
}

func (gk generateNbKindStruct) EmitNodebuilderMethodAmendMap(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) AmendMap() (ipld.MapBuilder, error) {
			panic("TODO later")
		}
	`, w, gk)
}
