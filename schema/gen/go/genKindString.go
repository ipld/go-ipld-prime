package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindString(t schema.Type) typedNodeGenerator {
	return generateKindString{
		t.(schema.TypeString),
		generateKindedRejections_String{
			mungeTypeNodeIdent(t),
			string(t.Name()),
		},
	}
}

type generateKindString struct {
	Type schema.TypeString
	generateKindedRejections_String
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindString) EmitNativeType(w io.Writer) {
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }} struct{ x string }

	`, w, gk)
}

func (gk generateKindString) EmitNativeAccessors(w io.Writer) {
	// The node interface's `AsString` method is almost sufficient... but
	//  this method unboxes without needing to return an error that's statically impossible,
	//   which makes it easier to use in chaining.
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) String() string {
			return x.x
		}
	`, w, gk)
}

func (gk generateKindString) EmitNativeBuilder(w io.Writer) {
	// Having a builder for scalar kinds seems overkill, but it gives us a place to do validations,
	//  it keeps things consistent, and it lets us do 'Build' and 'MustBuild' without two top-level symbols.
	//   The ergonomics of this will be worth reviewing once we get a more holistic overview of the finished system, though.
	doTemplate(`
		type {{ .Type | mungeTypeNodeIdent }}__Content struct {
			Value string
		}

		func (b {{ .Type | mungeTypeNodeIdent }}__Content) Build() ({{ .Type | mungeTypeNodeIdent }}, error) {
			x := {{ .Type | mungeTypeNodeIdent }}{
				b.Value,
			}
			// FUTURE : want to support customizable validation.
			//   but 'if v, ok := x.(schema.Validatable); ok {' doesn't fly: need a way to work on concrete types.
			return x, nil
		}
		func (b {{ .Type | mungeTypeNodeIdent }}__Content) MustBuild() {{ .Type | mungeTypeNodeIdent }} {
			if x, err := b.Build(); err != nil {
				panic(err)
			} else {
				return x
			}
		}

	`, w, gk)
}

func (gk generateKindString) EmitNativeMaybe(w io.Writer) {
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

func (gk generateKindString) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Type | mungeTypeNodeIdent }}{}
		var _ typed.Node = {{ .Type | mungeTypeNodeIdent }}{}

	`, w, gk)
}

func (gk generateKindString) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func (x {{ .Type | mungeTypeNodeIdent }}) AsString() (string, error) {
			return x.x, nil
		}
	`, w, gk)
}

func (gk generateKindString) EmitTypedNodeMethodRepresentation(w io.Writer) {
	doTemplate(`
		func ({{ .Type | mungeTypeNodeIdent }}) Representation() ipld.Node {
			panic("TODO representation")
		}
	`, w, gk)
}

func (gk generateKindString) GetRepresentationNodeGen() nodeGenerator {
	return nil // TODO of course
}
