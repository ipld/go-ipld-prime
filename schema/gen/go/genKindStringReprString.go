package gengo

import (
	"io"
)

func (gk generateKindString) EmitNodeType(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = {{ .Name }}{}

		type {{ .Name }} struct { x string }
	`, w, gk)
}

func (gk generateKindString) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func ({{ .Name }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`, w, gk)
}
