package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

// --- type-semantics node interface satisfaction --->

func (gk generateKindList) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ schema.TypedNode = {{ .Type | mungeTypeNodeIdent }}{}

	`, w, gk)
}

func (gk generateKindList) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindList) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_List
		}
	`, w, gk)
}

func (gk generateKindList) EmitNodeMethodLookupIndex(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) LookupIndex(index int) (ipld.Node, error) {
			if index >= len(x.x) {
				return nil, ipld.ErrNotExists{ipld.PathSegmentOfInt(index)}
			}
			{{- if .Type.ValueIsNullable }}
			if x.x[index] == nil {
				return ipld.Null, nil
			}
			return *x.x[index], nil
			{{- else }}
			return x.x[index], nil
			{{- end }}
		}
	`, w, gk)
}

func (gk generateKindList) EmitNodeMethodLookup(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) Lookup(key ipld.Node) (ipld.Node, error) {
			ki, err := key.AsInt()
			if err != nil {
				return nil, ipld.ErrInvalidKey{"got " + key.ReprKind().String() + ", need Int"}
			}
			return x.LookupIndex(ki)
		}
	`, w, gk)
}

func (gk generateKindList) EmitNodeMethodListIterator(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) ListIterator() ipld.ListIterator {
			return &{{ .Type | mungeTypeNodeItrIdent }}{&x, 0}
		}

		type {{ .Type | mungeTypeNodeItrIdent }} struct {
			node *{{ .Type | mungeTypeNodeIdent }}
			idx  int
		}

		func (itr *{{ .Type | mungeTypeNodeItrIdent }}) Next() (idx int, value ipld.Node, _ error)	{
			if itr.idx >= len(itr.node.x) {
				return 0, nil, ipld.ErrIteratorOverread{}
			}
			idx = itr.idx
			{{- if .Type.ValueIsNullable }}
			if itr.node.x[idx] == nil {
				value = ipld.Null
			} else {
				value = *itr.node.x[idx]
			}
			{{- else }}
			value = itr.node.x[idx]
			{{- end }}
			itr.idx++
			return
		}

		func (itr *{{ .Type | mungeTypeNodeItrIdent }}) Done() bool {
			return itr.idx >= len(itr.node.x)
		}

	`, w, gk)
}

func (gk generateKindList) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) Length() int {
			return len(x.x)
		}
	`, w, gk)
}

// --- type-semantics nodebuilder --->

func (gk generateKindList) EmitNodeMethodNodeBuilder(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) NodeBuilder() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateKindList) GetNodeBuilderGen() nodebuilderGenerator {
	return generateNbKindList{
		gk.Type,
		genKindedNbRejections_List{
			mungeTypeNodebuilderIdent(gk.Type),
			string(gk.Type.Name()) + ".Builder",
		},
	}
}

type generateNbKindList struct {
	Type schema.TypeList
	genKindedNbRejections_List
}

func (gk generateNbKindList) EmitNodebuilderType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodebuilderIdent }} struct{}

	`, w, gk)
}

func (gk generateNbKindList) EmitNodebuilderConstructor(w io.Writer) {
	doTemplate(`
		func {{ .Type | mungeNodebuilderConstructorIdent }}() ipld.NodeBuilder {
			return {{ .Type | mungeTypeNodebuilderIdent }}{}
		}
	`, w, gk)
}

func (gk generateNbKindList) EmitNodebuilderMethodCreateList(w io.Writer) {
	// Some interesting edge cases to note:
	//  - This builder, being all about semantics and not at all about serialization,
	//      is order-insensitive.
	//  - We don't specially handle being given 'undef' as a value.
	//      It just falls into the "need a schema.TypedNode" error bucket.
	//  - We only accept *codegenerated values* -- a schema.TypedNode created
	//      in the same schema universe *isn't accepted*.
	//        REVIEW: We could try to accept those, but it might have perf/sloc costs,
	//          and it's hard to imagine a user story that gets here.
	//  - The has-been-set-if-required validation is fun; it only requires state
	//     for non-optional fields, and that often gets a little hard to follow
	//       because it gets wedged in with other logic tables around optionality.
	// REVIEW: 'x, ok := v.({{ $field.Type.Name }})' might need some stars in it... sometimes.
	// TODO : review the panic of `ErrNoSuchField` in `BuilderForValue` --
	//  see the comments in the NodeBuilder interface for the open questions on this topic.
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) CreateList() (ipld.ListBuilder, error) {
			return &{{ .Type | mungeTypeNodeListBuilderIdent }}{v:&{{ .Type | mungeTypeNodeIdent }}{}}, nil
		}

		type {{ .Type | mungeTypeNodeListBuilderIdent }} struct{
			v *{{ .Type | mungeTypeNodeIdent }}
		}

		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) growList(k int) {
			oldLen := len(lb.v.x)
			minLen := k + 1
			if minLen > oldLen {
				// Grow.
				oldCap := cap(lb.v.x)
				if minLen > oldCap {
					// Out of cap; do whole new backing array allocation.
					//  Growth maths are per stdlib's reflect.grow.
					// First figure out how much growth to do.
					newCap := oldCap
					if newCap == 0 {
						newCap = minLen
					} else {
						for minLen > newCap {
							if minLen < 1024 {
								newCap += newCap
							} else {
								newCap += newCap / 4
							}
						}
					}
					// Now alloc and copy over old.
					newArr := make([]{{if .Type.ValueIsNullable}}*{{end}}{{.Type.ValueType | mungeTypeNodeIdent}}, minLen, newCap)
					copy(newArr, lb.v.x)
					lb.v.x = newArr
				} else {
					// Still have cap, just extend the slice.
					lb.v.x = lb.v.x[0:minLen]
				}
			}
		}

		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) validate(v ipld.Node) error {
			{{- if .Type.ValueIsNullable }}
			if v.IsNull() {
				return nil
			}
			{{- else}}
			if v.IsNull() {
				panic("type mismatch on struct field assignment: cannot assign null to non-nullable field") // FIXME need an error type for this
			}
			{{- end}}
			tv, ok := v.(schema.TypedNode)
			if !ok {
				panic("need schema.TypedNode for insertion into struct") // FIXME need an error type for this
			}
			_, ok = v.({{ .Type.ValueType | mungeTypeNodeIdent }})
			if !ok {
				panic("value for type {{.Type.Name}} is type {{.Type.ValueType.Name}}; cannot assign "+tv.Type().Name()) // FIXME need an error type for this
			}
			return nil
		}

		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) unsafeSet(idx int, v ipld.Node) {
			{{- if .Type.ValueIsNullable }}
			if v.IsNull() {
				lb.v.x[idx] = nil
				return
			}
			{{- end}}
			x := v.({{ .Type.ValueType | mungeTypeNodeIdent }})
			{{- if .Type.ValueIsNullable }}
			lb.v.x[idx] = &x
			{{- else}}
			lb.v.x[idx] = x
			{{- end}}
		}

		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) AppendAll(vs []ipld.Node) error {
			for _, v := range vs {
				err := lb.validate(v)
				if err != nil {
					return err
				}
			}
			off := len(lb.v.x)
			new := off + len(vs)
			lb.growList(new-1)
			for _, v := range vs {
				lb.unsafeSet(off, v)
				off++
			}
			return nil
		}

		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) Append(v ipld.Node) error {
			err := lb.validate(v)
			if err != nil {
				return err
			}
			off := len(lb.v.x)
			lb.growList(off)
			lb.unsafeSet(off, v)
			return nil
		}
		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) Set(idx int, v ipld.Node) error {
			err := lb.validate(v)
			if err != nil {
				return err
			}
			lb.growList(idx)
			lb.unsafeSet(idx, v)
			return nil
		}

		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) Build() (ipld.Node, error) {
			v := *lb.v
			lb = nil
			return v, nil
		}

		func (lb *{{ .Type | mungeTypeNodeListBuilderIdent }}) BuilderForValue(_ int) ipld.NodeBuilder {
			return {{ .Type.ValueType | mungeNodebuilderConstructorIdent }}()
		}

	`, w, gk)
}

func (gk generateNbKindList) EmitNodebuilderMethodAmendList(w io.Writer) {
	doTemplate(`
		func (nb {{ .Type | mungeTypeNodebuilderIdent }}) AmendList() (ipld.ListBuilder, error) {
			panic("TODO later")
		}
	`, w, gk)
}

// --- entrypoints to representation --->

func (gk generateKindList) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | mungeTypeNodeIdent }}) Representation() ipld.Node {
			panic("TODO representation")
		}
	`, w, gk)
}

func (gk generateKindList) GetRepresentationNodeGen() nodeGenerator {
	return nil // TODO of course
}
