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

func (gk generateKindList) EmitNativeType(w io.Writer) {
	// Observe that we get a '*' if the values are nullable.
	//  FUTURE: worth reviewing if this could or should use 'maybe' structs instead of pointers
	//   (which would effectively trade alloc count vs size for very different performance characteristics).
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }} struct{
			x []{{if .Type.ValueIsNullable}}*{{end}}{{.Type.ValueType | mungeTypeNodeIdent}}
		}
	`, w, gk)
}

func (gk generateKindList) EmitNativeAccessors(w io.Writer) {
	doTemplate(`
		// TODO generateKindList.EmitNativeAccessors
	`, w, gk)
}

func (gk generateKindList) EmitNativeBuilder(w io.Writer) {
	doTemplate(`
		// TODO generateKindList.EmitNativeBuilder
	`, w, gk)
}

func (gk generateKindList) EmitNativeMaybe(w io.Writer) {
	// TODO this can most likely be extracted and DRY'd, just not 100% sure yet
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
