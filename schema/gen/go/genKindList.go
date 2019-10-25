package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindList(t schema.Type) typedNodeGenerator {
	return generateKindList{
		t.(schema.TypeList),
		generateKindedRejections_List{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindList struct {
	Type schema.TypeList
	generateKindedRejections_List
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindList) EmitNodeType(w io.Writer) {
	// Observe that we get a '*' if a field is *either* nullable *or* optional;
	//  and we get an extra bool for the second cardinality +1'er if both are true.
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

		type {{ .Type | mungeTypeNodeIdent }} struct{
			x []{{if .Type.ValueIsNullable}}*{{end}}{{.Type.ValueType | mungeTypeNodeIdent}}
		}
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
