package gengo

import (
	"io"
	"strings"
	"text/template"

	wish "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
)

func doTemplate(tmplstr string, w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	tmpl := template.Must(template.New("").
		Funcs(template.FuncMap{

			// These methods are used for symbol munging and appear constantly, so they need to be short.
			//  (You could also get at them through `.AdjCfg`, but going direct saves some screen real estate.)
			"TypeSymbol":       adjCfg.TypeSymbol,
			"FieldSymbolLower": adjCfg.FieldSymbolLower,
			"FieldSymbolUpper": adjCfg.FieldSymbolUpper,
			"MaybeUsesPtr":     adjCfg.MaybeUsesPtr,

			// The whole AdjunctConfig can be accessed.
			//  Access methods like UnionMemlayout through this, as e.g. `.AdjCfg.UnionMemlayout`.
			"AdjCfg": func() *AdjunctCfg { return adjCfg },

			"KindPrim": func(k ipld.ReprKind) string {
				switch k {
				case ipld.ReprKind_Map:
					panic("this isn't useful for non-scalars")
				case ipld.ReprKind_List:
					panic("this isn't useful for non-scalars")
				case ipld.ReprKind_Null:
					panic("this isn't useful for null")
				case ipld.ReprKind_Bool:
					return "bool"
				case ipld.ReprKind_Int:
					return "int"
				case ipld.ReprKind_Float:
					return "float64"
				case ipld.ReprKind_String:
					return "string"
				case ipld.ReprKind_Bytes:
					return "[]byte"
				case ipld.ReprKind_Link:
					return "ipld.Link"
				default:
					panic("invalid enumeration value!")
				}
			},
			"add":   func(a, b int) int { return a + b },
			"title": func(s string) string { return strings.Title(s) },
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
