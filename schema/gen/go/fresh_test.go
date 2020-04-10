package gengo

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

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

func genAndCompilerAndTest(
	t *testing.T,
	prefix string,
	pkgName string,
	ts schema.TypeSystem,
	adjCfg *AdjunctCfg,
	testFnName string, // disturbing, i know
) {
	t.Run("generate", func(t *testing.T) {
		// Emit fixed bits.
		withFile(prefix+"/minima.go", func(f io.Writer) {
			EmitInternalEnums(pkgName, f)
		})

		// Emit a file for each type.
		for _, typ := range ts.GetTypes() {
			withFile(prefix+"/t"+typ.Name().String()+".go", func(f io.Writer) {
				EmitFileHeader(pkgName, f)
				switch t2 := typ.(type) {
				case schema.TypeString:
					EmitEntireType(NewStringReprStringGenerator(pkgName, t2, adjCfg), f)
				}
			})
		}

		// Emit an exported top level function for getting nodestyles.
		//  (This part isn't necessary except for a special need we have with this plugin trick;
		//   normally, user code uses the `{pkgname}.Style.{TypeName}` constant access.)
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
			withFile(prefix+"/test_test.go", func(w io.Writer) {
				doTemplate(`
					package `+pkgName+`

					import "testing"
					import "github.com/ipld/go-ipld-prime/schema/gen/go/tests"

					func TestAll(t *testing.T)  {
						tests.{{ . }}(t, GetStyleByName)
					}
				`, w, adjCfg, testFnName)
			})
			cmd := exec.Command("go", "test", "-run=^$", "./_test/"+prefix)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				panic(err)
			}

			t.Run("test", func(t *testing.T) {
				cmd := exec.Command("go", "test", "-v", "./_test/"+prefix)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					panic(err)
				}
			})
		})
	})
}

func TestFancier(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}

	ts.Accumulate(schema.SpawnString("String"))
	adjCfg.maybeUsesPtr["String"] = false

	prefix := "foo"
	pkgName := prefix
	genAndCompilerAndTest(t, prefix, pkgName, ts, adjCfg, "ExerciseString")
}

func TestFanciest(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}

	ts.Accumulate(schema.SpawnString("String"))
	adjCfg.maybeUsesPtr["String"] = true

	prefix := "bar"
	pkgName := prefix
	genAndCompilerAndTest(t, prefix, pkgName, ts, adjCfg, "ExerciseString")
}
