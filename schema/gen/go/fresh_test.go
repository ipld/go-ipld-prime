package gengo

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func invokeBuildPlugin(prefix string) {
	cmd := exec.Command("go", "build", "-o=./_test/"+prefix+"/obj.so", "-buildmode=plugin", "./_test/"+prefix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
func loadPlugin(prefix string) *plugin.Plugin {
	plg, err := plugin.Open("./_test/" + prefix + "/obj.so")
	if err != nil {
		panic(err)
	}
	return plg
}

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
	tests func(t *testing.T, getStyleByName func(string) ipld.NodeStyle),
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
			invokeBuildPlugin(prefix)
			plg := loadPlugin(prefix)

			sym, err := plg.Lookup("GetStyleByName")
			if err != nil {
				panic(err)
			}
			getStyleByName := sym.(func(string) ipld.NodeStyle)

			t.Run("test", func(t *testing.T) {
				tests(t, getStyleByName)
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
	pkgName := "main" // has to be 'main' for plugins to work.  this stricture makes little sense to me, but i didn't write the rules.
	genAndCompilerAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		ns := getStyleByName("String")
		t.Run("string operations work", func(t *testing.T) {
			nb := ns.NewBuilder()
			nb.AssignString("woiu")
			n := nb.Build()
			t.Logf("%v\n", n)
		})
		t.Run("null is rejected", func(t *testing.T) {
			nb := ns.NewBuilder()
			nb.AssignNull()

		})
	})
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
	pkgName := "main" // has to be 'main' for plugins to work.  this stricture makes little sense to me, but i didn't write the rules.
	genAndCompilerAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		ns := getStyleByName("String")
		t.Run("string operations work", func(t *testing.T) {
			nb := ns.NewBuilder()
			nb.AssignString("woiu")
			n := nb.Build()
			t.Logf("%v\n", n)
		})
		t.Run("null is rejected", func(t *testing.T) {
			nb := ns.NewBuilder()
			nb.AssignNull()

		})
	})
}
