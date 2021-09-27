package basic_test

import (
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime/app"

	"github.com/warpfork/go-testmark"
	"github.com/warpfork/go-testmark/testexec"
)

func TestRead(t *testing.T) {
	// Just real quick compile the whole app.  We need this so we can test it in scripts.
	os.MkdirAll("/tmp/ipld-test/bin/", 0755)
	exec.Command("go", "build", "-o", "/tmp/ipld-test/bin/ipld", "../cmd/ipld/ipld.go").Run()

	filename := "../docs/read.md"
	doc, err := testmark.ReadFile(filename)
	if err != nil {
		t.Fatalf("spec file parse failed?!: %s", err)
	}
	pa := &testmark.PatchAccumulator{}

	// Data hunk in this spec file are in "directories" of a test scenario each.
	doc.BuildDirIndex()
	for _, dir := range doc.DirEnt.ChildrenList {
		t.Run(dir.Name, func(t *testing.T) {
			testexec.Tester{
				ExecFn: app.Main,
				ScriptFn: func(script string, stdin io.Reader, stdout, stderr io.Writer) (exitcode int, oshit error) {
					return testexec.ScriptFn_ExecBash("export PATH=$PATH:/tmp/ipld-test/bin/;\n"+script, stdin, stdout, stderr)
				},
				AssertFn: func(t *testing.T, actual, expect string) {
					quicktest.Assert(t, actual, quicktest.CmpEquals(), expect)
				},
				Patches: pa,
			}.Test(t, dir)
		})
	}

	pa.WriteFileWithPatches(doc, filename)
}
