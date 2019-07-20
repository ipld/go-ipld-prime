package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func NewGeneratorForKindString(t schema.Type) typeGenerator {
	return generateKindString{
		t,
		generateKindedRejections_String{t},
	}
}

type generateKindString struct {
	schema.Type
	generateKindedRejections_String
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}

func (gk generateKindString) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func (x {{ .Name }}) AsString() (string, error) {
			return x.x, nil
		}
	`, w, gk)
}
