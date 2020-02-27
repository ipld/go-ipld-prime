package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindBytes(t schema.Type) typedNodeGenerator {
	return generateKindBytes{
		t.(schema.TypeBytes),
		generateKindedRejections_Bytes{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindBytes struct {
	Type schema.TypeBytes
	generateKindedRejections_Bytes
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindBytes) EmitNativeType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }} struct{ x []byte }

	`, w, gk)
}

func (gk generateKindBytes) EmitNativeAccessors(w io.Writer) {
	doTemplate(`
		// TODO generateKindBytes.EmitNativeAccessors
	`, w, gk)
}

func (gk generateKindBytes) EmitNativeBuilder(w io.Writer) {
	doTemplate(`
		// TODO generateKindBytes.EmitNativeBuilder
	`, w, gk)
}

func (gk generateKindBytes) EmitNativeMaybe(w io.Writer) {
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
