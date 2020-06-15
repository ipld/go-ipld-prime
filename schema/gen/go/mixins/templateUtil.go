package mixins

import (
	"io"
	"strings"
	"text/template"

	wish "github.com/warpfork/go-wish"
)

func doTemplate(tmplstr string, w io.Writer, data interface{}) {
	tmpl := template.Must(template.New("").
		Funcs(template.FuncMap{
			"title": func(s string) string { return strings.Title(s) },
		}).
		Parse(wish.Dedent(tmplstr)))
	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}
