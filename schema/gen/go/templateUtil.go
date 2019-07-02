package gengo

import (
	"io"
	"text/template"

	wish "github.com/warpfork/go-wish"
)

func doTemplate(tmplstr string, w io.Writer, data interface{}) {
	tmpl := template.Must(template.New("").Parse("\n" + wish.Dedent(tmplstr)))
	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}
