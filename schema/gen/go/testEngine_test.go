package gengo

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
)

var _ tests.Engine = (*genAndCompileEngine)(nil)

type genAndCompileEngine struct {
	subtestName string
	prefix      string

	adjCfg AdjunctCfg

	prototypeByName func(string) datamodel.NodePrototype
}

var tmpGenBuildDir = filepath.Join(os.TempDir(), "test-go-ipld-prime-gengo")

func (e *genAndCompileEngine) Init(t *testing.T, ts schema.TypeSystem) {
	// Make directories for the package we're about to generate.
	// They will live in a temporary directory, usually
	// /tmp/test-go-ipld-prime-gengo on Linux. It can be removed at any time.
	// We don't by default because it's nicer to let go's builds of things cache.
	// If you change the names of types, though, you'll have garbage files leftover,
	// and that's currently a manual cleanup problem.  Sorry.
	dir := filepath.Join(tmpGenBuildDir, e.prefix)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	pkgName := "main"

	// Generate... everything, really.
	Generate(dir, pkgName, ts, &e.adjCfg)

	// Emit an exported top level function for getting NodePrototype.
	//  This part isn't necessary except for a special need we have with this plugin trick;
	//   normally, user code uses the `{pkgname}.Prototype.{TypeName}` constant (so-to-speak, anyway) to get a hold of NodePrototypes...
	//   but for plugins, we need a top-level exported symbol to grab ahold of, and we can't easily look through the `Prototype` value
	//    without an interface... so we generate this function to fit the bill instead.
	withFile(filepath.Join(dir, "prototypeGetter.go"), func(w io.Writer) {
		doTemplate(`
			package `+pkgName+`

			import "github.com/ipld/go-ipld-prime/datamodel"

			func GetPrototypeByName(name string) datamodel.NodePrototype {
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
		`, w, &e.adjCfg, ts.GetTypes())
	})

	// Build the genned code.
	//  This will either make a plugin (which we can run behavioral tests on next!),
	//  or just build it quietly just to see if there are compile-time errors,
	//  depending on your build tags.
	// See 'HACKME_testing.md' for discussion.
	buildGennedCode(t, e.prefix, pkgName)

	e.prototypeByName = fnPrototypeByName(e.prefix)
}

func (e *genAndCompileEngine) PrototypeByName(name string) datamodel.NodePrototype {
	return e.prototypeByName(name)
}
