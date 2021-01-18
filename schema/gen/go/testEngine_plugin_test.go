// +build cgo,!skipgenbehavtests

package gengo

import (
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"testing"

	"github.com/ipld/go-ipld-prime"
)

func objPath(prefix string) string {
	return filepath.Join(tmpGenBuildDir, prefix, "obj.so")
}

func buildGennedCode(t *testing.T, prefix string, _ string) {
	// Invoke `go build` with flags to create a plugin -- we'll be able to
	//  load into this plugin into this selfsame process momentarily.
	// Use globbing, because these are files outside our module.
	files, err := filepath.Glob(filepath.Join(tmpGenBuildDir, prefix, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"build", "-o=" + objPath(prefix), "-buildmode=plugin"}
	args = append(args, files...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("genned code failed to compile: %s", err)
	}
}

func runBehavioralTests(t *testing.T, prefix string, testsFn behavioralTests) {
	plg, err := plugin.Open(objPath(prefix))
	if err != nil {
		panic(err) // Panic because if this was going to flunk, we expected it to flunk earlier when we ran 'go build'.
	}
	sym, err := plg.Lookup("GetPrototypeByName")
	if err != nil {
		panic(err)
	}
	getPrototypeByName := sym.(func(string) ipld.NodePrototype)

	t.Run("bhvtest", func(t *testing.T) {
		testsFn(t, getPrototypeByName)
	})
}
