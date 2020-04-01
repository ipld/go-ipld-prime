package gengo

import (
	"io"
	"text/template"

	wish "github.com/warpfork/go-wish"
)

func doTemplate(tmplstr string, w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	tmpl := template.Must(template.New("").
		Funcs(template.FuncMap{
			"TypeSymbol":       adjCfg.TypeSymbol,
			"FieldSymbolLower": adjCfg.FieldSymbolLower,
			"FieldSymbolUpper": adjCfg.FieldSymbolUpper,
			"add":              func(a, b int) int { return a + b },
		}).
		Parse(wish.Dedent(tmplstr)))
	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}
