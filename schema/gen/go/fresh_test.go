package gengo

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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
			cmd := exec.Command("go", "build", "./_test/"+prefix)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				panic(err)
			}

			t.Run("test", func(t *testing.T) {
				cmd := exec.Command("go", "test", "-json", "./_test/"+prefix)
				var p io.Reader
				p, _ = cmd.StdoutPipe()
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); err != nil {
					panic(err)
				}
				//p = io.TeeReader(p, os.Stdout) // uncomment to see it plain
				dec := json.NewDecoder(p)
				recurse(t, dec)
				if err := cmd.Wait(); err != nil {
					panic(err)
				}
			})
		})
	})
}

// This function tries to reconstruct the 't.Run' tree from a child process,
// but it's not at all correct.
// Since the messages come out of order (more than one thing starts before its
// siblings finish, etc), a lot of buffering would be needed to make this right.
func recurse(t *testing.T, dec *json.Decoder) {
	for dec.More() {
		type msg struct {
			Action  string
			Package string
			Test    string
			Output  string
			Elapsed float32
		}
		var st msg
		if err := dec.Decode(&st); err != nil {
			panic(err)
		}
		switch st.Action {
		case "run":
			t.Run(path.Base(st.Test), func(t *testing.T) {
				// munch one message, because it's reliably the thing announcing "=== RUN" for itself on output.
				if err := dec.Decode(&st); err != nil {
					panic(err)
				}
				recurse(t, dec)
			})
		case "pass":
			return
		case "output":
			trimmed := strings.TrimLeft(st.Output, " ")
			// Filter out "--- PASS" output.  i don't see any unambiguous way to do this;
			//  'println' in the tests can easily write the same magic strings.
			//   (Double fun?  It seems 'go test -json' will even itself be
			//    confused by such strings, and turn them into nonsense json anyway.
			//     Really, there's no winning possible anywhere near here.)
			switch {
			case strings.HasPrefix(trimmed, "--- PASS: "):
				continue
			case trimmed == "PASS\n":
				continue
			case strings.HasPrefix(trimmed, "ok  \t"):
				continue
			}
			t.Logf("%s", trimmed)
		default:
			panic(st.Action)
		}
	}
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
