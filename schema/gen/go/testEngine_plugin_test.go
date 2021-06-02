//go:build cgo && !skipgenbehavtests
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

func objPath(dirName string) string {
	return filepath.Join(tmpGenBuildDir, dirName, "obj.so")
}

func buildGennedCode(t *testing.T, dirName string, _ string) {
	// Invoke `go build` with flags to create a plugin -- we'll be able to
	//  load into this plugin into this selfsame process momentarily.
	// Use globbing, because these are files outside our module.
	files, err := filepath.Glob(filepath.Join(tmpGenBuildDir, dirName, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"build", "-o=" + objPath(dirName), "-buildmode=plugin"}
	args = append(args, files...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("genned code failed to compile: %s", err)
	}
}

func fnPrototypeByName(dirName string) func(string) ipld.NodePrototype {
	plg, err := plugin.Open(objPath(dirName))
	if err != nil {
		panic(err) // Panic because if this was going to flunk, we expected it to flunk earlier when we ran 'go build'.
	}
	sym, err := plg.Lookup("GetPrototypeByName")
	if err != nil {
		panic(err)
	}
	return sym.(func(string) ipld.NodePrototype)
}
