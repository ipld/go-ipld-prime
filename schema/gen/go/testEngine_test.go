package gengo

import (
	"io"
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

// behavioralTests describes the interface of a function that can be supplied
// in order to run tests on generated code.
//
// The getPrototypeByName function can get its job done using only interface types
// that we already know from outside any generated code, so you can write tests
// that have no _compile time_ dependency on the generated code.  This makes it
// easier for IDEs and suchlike to help you write and check the test functions.
//
// Ask for prototypes using the type name alone (no package prefix);
// their representation prototypes can be obtained by appending ".Repr".
type behavioralTests func(t *testing.T, getPrototypeByName func(string) ipld.NodePrototype)

func genAndCompileAndTest(
	t *testing.T,
	prefix string,
	pkgName string,
	ts schema.TypeSystem,
	adjCfg *AdjunctCfg,
	testsFn behavioralTests,
) {
	t.Run("generate", func(t *testing.T) {
		// Make directories for the package we're about to generate.
		//  Everything will be prefixed with "./_test".
		// You can rm-rf the whole "./_test" dir at your leisure.
		//  We don't by default because it's nicer to let go's builds of things cache.
		//  If you change the names of types, though, you'll have garbage files leftover,
		//   and that's currently a manual cleanup problem.  Sorry.
		os.Mkdir("./_test/", 0755)
		os.Mkdir("./_test/"+prefix, 0755)

		// Generate... everything, really.
		Generate("./_test/"+prefix, pkgName, ts, adjCfg)

		// Emit an exported top level function for getting NodePrototype.
		//  This part isn't necessary except for a special need we have with this plugin trick;
		//   normally, user code uses the `{pkgname}.Prototype.{TypeName}` constant (so-to-speak, anyway) to get a hold of NodePrototypes...
		//   but for plugins, we need a top-level exported symbol to grab ahold of, and we can't easily look through the `Prototype` value
		//    without an interface... so we generate this function to fit the bill instead.
		withFile("./_test/"+prefix+"/prototypeGetter.go", func(w io.Writer) {
			doTemplate(`
				package `+pkgName+`

				import "github.com/ipld/go-ipld-prime"

				func GetPrototypeByName(name string) ipld.NodePrototype {
					switch name {
					{{- range . }}
					case "{{ .Name }}":
						return _{{ . | TypeSymbol }}__Prototype{}
					case "{{ .Name }}.Repr":
						return _{{ . | TypeSymbol }}__ReprPrototype{}
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
