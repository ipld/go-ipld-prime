package amend

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/warpfork/go-testmark"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/traversal/patch"
)

func TestSpecFixtures(t *testing.T) {
	dir := "../../.ipld/specs/patch/fixtures/"
	testOneSpecFixtureFile(t, dir+"fixtures-1.md")
}

func testOneSpecFixtureFile(t *testing.T, filename string) {
	doc, err := testmark.ReadFile(filename)
	if os.IsNotExist(err) {
		t.Skipf("not running spec suite: %s (did you clone the submodule with the data?)", err)
	}
	if err != nil {
		t.Fatalf("spec file parse failed?!: %s", err)
	}

	// Data hunk in this spec file are in "directories" of a test scenario each.
	doc.BuildDirIndex()

	for _, dir := range doc.DirEnt.ChildrenList {
		t.Run(dir.Name, func(t *testing.T) {
			// Grab all the data hunks.
			//  Each "directory" contains three piece of data:
			//   - `initial` -- this is the "block".  It's arbitrary example data.  They're all in json (or dag-json) format, for simplicity.
			//   - `patch` -- this is a list of patch ops.  Again, as json.
			//   - `result` -- this is the expected result object.  Again, as json.
			initialBlob := dir.Children["initial"].Hunk.Body
			patchBlob := dir.Children["patch"].Hunk.Body
			resultBlob := dir.Children["result"].Hunk.Body

			// Parse everything.
			initial, err := ipld.Decode(initialBlob, dagjson.Decode)
			if err != nil {
				t.Fatalf("failed to parse fixture data: %s", err)
			}
			ops, err := patch.ParseBytes(patchBlob, dagjson.Decode)
			if err != nil {
				t.Fatalf("failed to parse fixture patch: %s", err)
			}
			// We don't actually keep the decoded result object.  We're just gonna serialize the result and textually diff that instead.
			_, err = ipld.Decode(resultBlob, dagjson.Decode)
			if err != nil {
				t.Fatalf("failed to parse fixture data: %s", err)
			}

			// Do the thing!
			actualResult, err := Eval(initial, ops)
			if strings.HasSuffix(dir.Name, "-fail") {
				if err == nil {
					t.Fatalf("patch was expected to fail")
				} else {
					return
				}
			} else {
				if err != nil {
					t.Fatalf("patch did not apply: %s", err)
				}
			}

			// Serialize (and pretty print) result, so that we can diff it.
			actualResultBlob, err := ipld.Encode(actualResult, dagjson.EncodeOptions{
				EncodeLinks: true,
				EncodeBytes: true,
				MapSortMode: codec.MapSortMode_None,
			}.Encode)
			if err != nil {
				t.Errorf("failed to reserialize result: %s", err)
			}
			var actualResultBlobPretty bytes.Buffer
			json.Indent(&actualResultBlobPretty, actualResultBlob, "", "\t")

			// Diff!
			qt.Assert(t, actualResultBlobPretty.String()+"\n", qt.Equals, string(resultBlob))
		})
	}
}
