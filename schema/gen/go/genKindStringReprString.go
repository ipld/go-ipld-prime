package gengo

import (
	"io"
	"text/template"

	wish "github.com/warpfork/go-wish"
)

func (gk generateKindString) EmitNodeType(w io.Writer) {
	template.Must(template.New("").Parse("\n"+wish.Dedent(`
		var _ ipld.Node = {{ .Name }}{}

		type {{ .Name }} struct { x string }
	`))).Execute(w, gk)
}

func (gk generateKindString) EmitNodeMethodReprKind(w io.Writer) {
	template.Must(template.New("").Parse("\n"+wish.Dedent(`
		func ({{ .Name }}) ReprKind() ipld.ReprKind {
			return ipld.ReprKind_String
		}
	`))).Execute(w, gk)
}
