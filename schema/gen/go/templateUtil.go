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
		}).
		Parse(wish.Dedent(tmplstr)))
	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}
