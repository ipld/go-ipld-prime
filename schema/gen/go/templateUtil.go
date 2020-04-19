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
			"MaybeUsesPtr":     adjCfg.MaybeUsesPtr,
			"add":              func(a, b int) int { return a + b },
		}).
		// Seriously consider prepending `{{ $dot := . }}` (or 'top', or something).
		// Or a func into the map that causes `dot` to mean `func() interface{} { return data }`.
		// The number of times that the range feature has a dummy capture line above it is... not reasonable.
		//  (Grep for "/* ranging modifies dot, unhelpfully */" -- empirically, it's over 20 times already.)
		Parse(wish.Dedent(tmplstr)))
	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}

// We really need to do some more composable stuff around here.
// Generators should probably be carrying down their own doTemplate methods that curry customizations.
// E.g., map generators would benefit hugely from being able to make a clause for "entTypeStrung", "mTypeStrung", etc.
//
// Open question: how exactly?  Should some of this stuff should be composed by:
//   - composing template fragments;
//   - amending the funcmap;
//   - computing the whole result and injecting it as a string;
//   - ... combinations of the above?
// Adding to the complexity of the question is that sometimes we want to be
//  doing composition inside the output (e.g. DRY by functions in the result,
//   rather than by DRY'ing the templates).
// Best practice to make this evolve nicely is not at all obvious to this author.
//
