package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindLink(t schema.Type) typedNodeGenerator {
	return generateKindLink{
		t.(schema.TypeLink),
		generateKindedRejections_Link{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindLink struct {
	Type schema.TypeLink
	generateKindedRejections_Link
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindLink) EmitNativeType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }} struct{ x ipld.Link }

	`, w, gk)
}

func (gk generateKindLink) EmitNativeAccessors(w io.Writer) {
	doTemplate(`
		// TODO generateKindLink.EmitNativeAccessors
	`, w, gk)
}

func (gk generateKindLink) EmitNativeBuilder(w io.Writer) {
	doTemplate(`
		// TODO generateKindLink.EmitNativeBuilder
	`, w, gk)
}

func (gk generateKindLink) EmitNativeMaybe(w io.Writer) {
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
