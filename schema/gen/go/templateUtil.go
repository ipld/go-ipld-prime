package gengo

import (
	"io"
	"text/template"

	ipld "github.com/ipld/go-ipld-prime"

	wish "github.com/warpfork/go-wish"
)

func doTemplate(tmplstr string, w io.Writer, data interface{}) {
	tmpl := template.Must(template.New("").
		Funcs(template.FuncMap{
			// 'ReprKindConst' returns the source-string for "ipld.ReprKind_{{Kind}}".
			"ReprKindConst": func(k ipld.ReprKind) string {
				return "ipld.ReprKind_" + k.String() // happens to be fairly trivial.
			},

			// 'Add' does what it says on the tin.
			"Add": func(a, b int) int {
				return a + b
			},

			"mungeTypeNodeIdent":            mungeTypeNodeIdent,
			"mungeTypeNodebuilderIdent":     mungeTypeNodebuilderIdent,
			"mungeTypeReprNodeIdent":        mungeTypeReprNodeIdent,
			"mungeTypeReprNodebuilderIdent": mungeTypeReprNodebuilderIdent,
		}).
		Parse(wish.Dedent(tmplstr)))
	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}
