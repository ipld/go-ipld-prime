package gengo

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func withFile(filename string, fn func(io.Writer)) {
	// Rm-rf the whole "./_test" dir at your leasure.
	//  We don't by default because it's nicer to let go's builds of things cache.
	//  If you change the names of types, though, you'll have garbage files leftover,
	//   and that's currently a manual cleanup problem.  Sorry.
	os.Mkdir(filepath.Dir("./_test/"), 0755)
	os.Mkdir(filepath.Dir("./_test/"+filename), 0755)
	f, err := os.OpenFile("./_test/"+filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fn(f)
}

// behavioralTests describes the interface of a function that can be supplied
// in order to run tests on generated code.
//
// The getStyleByName function can get its job done using only interface types
// that we already know from outside any generated code, so you can write tests
// that have no _compile time_ dependency on the generated code.  This makes it
// easier for IDEs and suchlike to help you write and check the test functions.
//
// Ask for styles using the type name alone (no package prefix);
// their representation styles can be obtained by appending ".Repr".
type behavioralTests func(t *testing.T, getStyleByName func(string) ipld.NodeStyle)

func genAndCompileAndTest(
	t *testing.T,
	prefix string,
	pkgName string,
	ts schema.TypeSystem,
	adjCfg *AdjunctCfg,
	testsFn behavioralTests,
) {
	t.Run("generate", func(t *testing.T) {
		// Emit fixed bits.
		withFile(prefix+"/minima.go", func(f io.Writer) {
			EmitInternalEnums(pkgName, f)
		})

		// Emit a file for each type.
		//  This contains a bunch of big switches for type and representation strategy,
		//   which will probably get hoisted out to an exported feature at some point.
		for _, typ := range ts.GetTypes() {
			withFile(prefix+"/t"+typ.Name().String()+".go", func(f io.Writer) {
				EmitFileHeader(pkgName, f)
				switch t2 := typ.(type) {
				case schema.TypeString:
					EmitEntireType(NewStringReprStringGenerator(pkgName, t2, adjCfg), f)
				case schema.TypeStruct:
					switch t2.RepresentationStrategy().(type) {
					case schema.StructRepresentation_Map:
						EmitEntireType(NewStructReprMapGenerator(pkgName, t2, adjCfg), f)
					default:
						panic("unrecognized struct representation strategy")
					}
				default:
					panic("add more type switches here :)")
				}
			})
		}

		// Emit an exported top level function for getting nodestyles.
		//  This part isn't necessary except for a special need we have with this plugin trick;
		//   normally, user code uses the `{pkgname}.Style.{TypeName}` constant (so-to-speak, anyway) to get a hold of nodestyles...
		//   but for plugins, we need a top-level exported symbol to grab ahold of, and we can't easily look through the `Style` value
		//    without an interface... so we generate this function to fit the bill instead.
		withFile(prefix+"/styleGetter.go", func(w io.Writer) {
			doTemplate(`
				package `+pkgName+`

				import "github.com/ipld/go-ipld-prime"

				func GetStyleByName(name string) ipld.NodeStyle {
					switch name {
					{{- range . }}
					case "{{ .Name }}":
						return _{{ . | TypeSymbol }}__Style{}
					case "{{ .Name }}.Repr":
						return _{{ . | TypeSymbol }}__ReprStyle{}
					{{- end}}
					default:
						return nil
					}
				}
			`, w, adjCfg, ts.GetTypes())
		})

		t.Run("compile", func(t *testing.T) {
			// Build the genned code.
			//  This will either make a plugin (which we can run behavioral tests on next!),
			//  or just build it quietly just to see if there are compile-time errors,
			//  depending on your build tags.
			// See 'HACKME_testing.md' for discussion.
			buildGennedCode(t, prefix, pkgName)

			// This will either load the plugin and run behavioral tests,
			//  or emit a dummy t.Run and a skip,
			//  depending on your build tags.
			// See 'HACKME_testing.md' for discussion.
			runBehavioralTests(t, prefix, testsFn)
		})
	})
}
