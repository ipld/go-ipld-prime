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
			"Comments":         adjCfg.Comments,

			// The whole AdjunctConfig can be accessed.
			//  Access methods like UnionMemlayout through this, as e.g. `.AdjCfg.UnionMemlayout`.
			"AdjCfg": func() *AdjunctCfg { return adjCfg },

			// "dot" is a dummy value that's equal to the original `.` expression, but stays there.
			//  Use this if you're inside a range or other feature that shifted the dot and you want the original.
			//  (This may seem silly, but empirically, I found myself writing a dummy line to store the value of dot before endering a range clause >20 times; that's plenty.)
			"dot": func() interface{} { return data },

			"KindPrim": func(k ipld.Kind) string {
				switch k {
				case ipld.Kind_Map:
					panic("this isn't useful for non-scalars")
				case ipld.Kind_List:
					panic("this isn't useful for non-scalars")
				case ipld.Kind_Null:
					panic("this isn't useful for null")
				case ipld.Kind_Bool:
					return "bool"
				case ipld.Kind_Int:
					return "int64"
				case ipld.Kind_Float:
					return "float64"
				case ipld.Kind_String:
					return "string"
				case ipld.Kind_Bytes:
					return "[]byte"
				case ipld.Kind_Link:
					return "ipld.Link"
				default:
					panic("invalid enumeration value!")
				}
			},
			"Kind": func(s string) ipld.Kind {
				switch s {
				case "map":
					return ipld.Kind_Map
				case "list":
					return ipld.Kind_List
				case "null":
					return ipld.Kind_Null
				case "bool":
					return ipld.Kind_Bool
				case "int":
					return ipld.Kind_Int
				case "float":
					return ipld.Kind_Float
				case "string":
					return ipld.Kind_String
				case "bytes":
					return ipld.Kind_Bytes
				case "link":
					return ipld.Kind_Link
				default:
					panic("invalid enumeration value!")
				}
			},
			"KindSymbol": func(k ipld.Kind) string {
				switch k {
				case ipld.Kind_Map:
					return "ipld.Kind_Map"
				case ipld.Kind_List:
					return "ipld.Kind_List"
				case ipld.Kind_Null:
					return "ipld.Kind_Null"
				case ipld.Kind_Bool:
					return "ipld.Kind_Bool"
				case ipld.Kind_Int:
					return "ipld.Kind_Int"
				case ipld.Kind_Float:
					return "ipld.Kind_Float"
				case ipld.Kind_String:
					return "ipld.Kind_String"
				case ipld.Kind_Bytes:
					return "ipld.Kind_Bytes"
				case ipld.Kind_Link:
					return "ipld.Kind_Link"
				default:
					panic("invalid enumeration value!")
				}
			},
			"add":   func(a, b int) int { return a + b },
			"title": func(s string) string { return strings.Title(s) },
		}).
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
