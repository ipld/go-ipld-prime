// +build cgo,!skipgenbehavtests

package gengo

import (
	"os"
	"os/exec"
	"plugin"
	"testing"

	"github.com/ipld/go-ipld-prime"
)

func buildGennedCode(t *testing.T, prefix string, _ string) {
	// Invoke `go build` with flags to create a plugin -- we'll be able to
	//  load into this plugin into this selfsame process momentarily.
	cmd := exec.Command("go", "build", "-o=./_test/"+prefix+"/obj.so", "-buildmode=plugin", "./_test/"+prefix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("genned code failed to compile: %s", err)
	}
}

func runBehavioralTests(t *testing.T, prefix string, testsFn behavioralTests) {
	plg, err := plugin.Open("./_test/" + prefix + "/obj.so")
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
