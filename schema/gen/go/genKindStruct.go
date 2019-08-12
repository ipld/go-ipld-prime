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

func (gk generateKindStruct) EmitNodeMethodLookupString(w io.Writer) {
	doTemplate(`
		func (x {{ .Type.Name }}) LookupString(key string) (ipld.Node, error) {
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
		func (x {{ .Type.Name }}) MapIterator() ipld.MapIterator {
			return &_{{ .Type.Name }}__itr{&x, 0}
		}

		type _{{ .Type.Name }}__itr struct {
			node *{{ .Type.Name }}
			idx  int
		}

		func (itr *_{{ .Type.Name }}__itr) Next() (k ipld.Node, v ipld.Node, _ error) {
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
		func (itr *_{{ .Type.Name }}__itr) Done() bool {
			return itr.idx >= {{ len .Type.Fields }}
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

// FIXME We actually need several NodeBuilders.
// There should be one for the type-level understanding of the structure;
// another for the parsing of serial data (this will varies by reprStrat!);
// and possibly a third, if the reprStrat has strict and loose variants.
//
// The body of this method started angling for the parser one, in strict mode...
// so this actually belongs elsewhere.
//
// Parking my cursor for a moment here; we need a better set of abstractions
// around generating nodebuilders before proceding with detangling this.
func (gk generateKindStruct) EmitNodeMethodNodeBuilder(w io.Writer) {
	// FIXME dummy value.  Below code doesn't compile yet.
	doTemplate(`
		func ({{ .Type.Name }}) NodeBuilder() ipld.NodeBuilder {
			return nil // TODO EmitNodeMethodNodeBuilder
		}
	`, w, gk)
	return

	// The 'idx' value is incremented to point at the next field we expect to receive.
	// At some indexes, we might accept multiple values (optional fields cause this);
	// in such cases, we can jump forward (but never back).
	doTemplate(`
		func ({{ .Type.Name }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type.Name }}__NodeBuilder{}
		}

		type {{ .Type.Name }}__NodeBuilder struct{}

		func (nb {{ .Type.Name }}__NodeBuilder) CreateMap() (ipld.MapBuilder, error) {
			return &{{ .Type.Name }}__MapBuilder{}, nil
		}

		type {{ .Type.Name }}__MapBuilder struct{
			node *{{ .Type.Name }}
			idx  int
		}

		func (mb *{{ .Type.Name }}__MapBuilder) Insert(k, v ipld.Node) error {
			if mb.idx >= {{ len .Type.Fields }} {
				return fmt.Errorf("no more fields expected")
			}
			mb.idx++
			switch mb.idx {
			{{- range $i, $field := .Type.Fields }}
			case {{ Add $i 1 }}:
				// TODO lookahead to see how many consecutive option fields there are.
				// TODO make a jumptable for those (fallthrough doesn't seem to help).
				// TODO anything you jump over: make sure it's marked undefined (in that field's idiom).
				// whee.
			{{- end}}
			default:
				panic("unreachable")
			}
		}
		func (mb *{{ .Type.Name }}__MapBuilder) Build() (ipld.Node, error) {
			panic("todo")
		}
	`, w, gk)
}
