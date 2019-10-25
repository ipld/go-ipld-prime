package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindInt(t schema.Type) typedNodeGenerator {
	return generateKindInt{
		t.(schema.TypeInt),
		generateKindedRejections_Int{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindInt struct {
	Type schema.TypeInt
	generateKindedRejections_Int
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindInt) EmitNativeType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }} struct{ x int }

	`, w, gk)
}

func (gk generateKindInt) EmitNativeAccessors(w io.Writer) {
	// The node interface's `AsInt` method is almost sufficient... but
	//  this method unboxes without needing to return an error that's statically impossible,
	//   which makes it easier to use in chaining.
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) Int() int {
			return x.x
		}
	`, w, gk)
}

func (gk generateKindInt) EmitNativeBuilder(w io.Writer) {
	doTemplate(`
		// TODO generateKindInt.EmitNativeBuilder
	`, w, gk)
}

func (gk generateKindInt) EmitNativeMaybe(w io.Writer) {
	// TODO this can most likely be extracted and DRY'd, just not 100% sure yet
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
