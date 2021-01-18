// +build skipgenbehavtests

package gengo

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

func buildGennedCode(t *testing.T, prefix string, pkgName string) {
	// Emit a small file with a 'main' method.
	//  'go build' doesn't like it we're in a package called "main" and there isn't one
	//   (and at the same time, plugins demand that they be in a package called 'main',
	//    so 'pkgName' in practice is almost always "main").
	//  I dunno, friend.  I didn't write the rules.
	if pkgName == "main" {
		withFile(filepath.Join(tmpGenBuildDir, prefix, "main.go"), func(w io.Writer) {
			fmt.Fprintf(w, "package %s\n\n", pkgName)
			fmt.Fprintf(w, "func main() {}\n")
		})
	}

	// Invoke 'go build' -- nothing fancy.
	files, err := filepath.Glob(filepath.Join(tmpGenBuildDir, prefix, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"build"}
	args = append(args, files...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("genned code failed to compile: %s", err)
	}
}

func runBehavioralTests(t *testing.T, _ string, _ behavioralTests) {
	t.Run("bhvtest=skip", func(t *testing.T) {
		t.Skip("behavioral tests for generated code skipped: you used the 'skipgenbehavtests' build tag.")
	})
}
