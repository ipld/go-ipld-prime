package basic_test

import (
	"testing"

	"github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime/app"

	"github.com/warpfork/go-testmark"
	"github.com/warpfork/go-testmark/testexec"
)

func TestRead(t *testing.T) {
	filename := "./read.md"
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
				AssertFn: func(t *testing.T, actual, expect string) {
					quicktest.Assert(t, actual, quicktest.CmpEquals(), expect)
				},
				Patches: pa,
			}.TestSequence(t, dir)
		})
	}

	pa.WriteFileWithPatches(doc, filename)
}
