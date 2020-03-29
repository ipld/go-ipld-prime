package gengo

import (
	"io"
	"text/template"

	wish "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/schema"
)

func doTemplate(tmplstr string, w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	tmpl := template.Must(template.New("").
		Funcs(template.FuncMap{
			"TypeSymbol":       adjCfg.TypeSymbol,
			"FieldSymbolLower": adjCfg.FieldSymbolLower,
			"FieldSymbolUpper": adjCfg.FieldSymbolUpper,
			"FieldTypeOrMaybe": func(f schema.StructField) string {
				// Returns the symbol used for embedding the field's type, or, the MaybeT of that type.
				// Shorthand for:
				//  `{{if or $field.IsOptional $field.IsNullable }}Maybe{{else}}_{{end}}{{ $field.Type | TypeSymbol }}`
				// REVIEW: still not sure if this is gonna be worth it.  Thought it would appear in more places; actually, is only one.
				if f.IsOptional() || f.IsNullable() {
					return "Maybe" + adjCfg.TypeSymbol(f.Type())
				}
				return "_" + adjCfg.TypeSymbol(f.Type())
			},
		}).
		Parse(wish.Dedent(tmplstr)))
	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}
