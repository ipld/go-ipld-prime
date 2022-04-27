package selector_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/warpfork/go-testmark"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

func TestSpecFixtures(t *testing.T) {
	dir := "../../.ipld/specs/selectors/fixtures/"
	testOneSpecFixtureFile(t, dir+"selector-fixtures-1.md")
	testOneSpecFixtureFile(t, dir+"selector-fixtures-recursion.md")
}

func testOneSpecFixtureFile(t *testing.T, filename string) {
	data, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		t.Skipf("not running spec suite: %s (did you clone the submodule with the data?)", err)
	}
	crre := regexp.MustCompile(`\r?\n`)
	data = []byte(crre.ReplaceAllString(string(data), "\n")) // fix windows carriage-return

	doc, err := testmark.Parse(data)
	if err != nil {
		t.Fatalf("spec file parse failed?!: %s", err)
	}

	// Data hunk in this spec file are in "directories" of a test scenario each.
	doc.BuildDirIndex()
	for _, dir := range doc.DirEnt.ChildrenList {
		t.Run(dir.Name, func(t *testing.T) {
			// Each "directory" contains three piece of data:
			//  - `data` -- this is the "block". It's arbitrary example data. They're all in json (or dag-json) format, for simplicity.
			//  - `selector` -- this is the selector. Again, as json.
			//  - `expect-visit` -- these are json lines (one json object on each line) containing description of each node that should be visited, in order.
			fixtureData := dir.Children["data"].Hunk.Body
			fixtureSelector := dir.Children["selector"].Hunk.Body
			fixtureExpect := dir.Children["expect-visit"].Hunk.Body

			// Parse data into DMT form.
			nb := basicnode.Prototype.Any.NewBuilder()
			if err := dagjson.Decode(nb, bytes.NewReader(fixtureData)); err != nil {
				t.Errorf("failed to parse fixture data: %s", err)
			}
			dataDmt := nb.Build()

			// Parse and compile Selector.
			// (This is already arguably a test event on its own.
			selector, err := selectorparse.ParseAndCompileJSONSelector(string(fixtureSelector))
			if err != nil {
				t.Errorf("failed to parse+compile selector: %s", err)
			}

			// Go!
			//  We'll store the logs of our visit events as... ipld Nodes, actually.
			//  This will make them easy to serialize, which is good for two reasons:
			//   at the end, we're actually going to... do that, and use string diffs for the final assertion
			//    (because string diffing is actually really nice for aggregate feedback in a system like this);
			//   and also that means we're ready to save updated serial data into the fixture files, if we did want to patch them.
			var visitLogs []datamodel.Node
			traversal.WalkAdv(dataDmt, selector, func(prog traversal.Progress, n datamodel.Node, reason traversal.VisitReason) error {
				// Munge info about where we are into DMT shaped like the expectation records in the fixture.
				visitEventDescr, err := qp.BuildMap(basicnode.Prototype.Any, 3, func(ma datamodel.MapAssembler) {
					qp.MapEntry(ma, "path", qp.String(prog.Path.String()))
					qp.MapEntry(ma, "node", qp.Map(1, func(ma datamodel.MapAssembler) {
						qp.MapEntry(ma, n.Kind().String(), func(na datamodel.NodeAssembler) {
							switch n.Kind() {
							case datamodel.Kind_Map, datamodel.Kind_List:
								na.AssignNull()
							default:
								na.AssignNode(n)
							}
						})
					}))
					qp.MapEntry(ma, "matched", qp.Bool(reason == traversal.VisitReason_SelectionMatch))
				})
				if reason == traversal.VisitReason_SelectionMatch && n.Kind() == datamodel.Kind_Bytes {
					if lbn, ok := n.(datamodel.LargeBytesNode); ok {
						rdr, err := lbn.AsLargeBytes()
						if err == nil {
							io.Copy(io.Discard, rdr)
						}
					}
					_, err := n.AsBytes()
					if err != nil {
						panic("insanity at a deeper level than this test's target")
					}
				}
				if err != nil {
					panic("insanity at a deeper level than this test's target")
				}
				visitLogs = append(visitLogs, visitEventDescr)
				return nil
			})

			// Brief detour -- we're going to bounce the fixture data through our own deserialize and serialize.
			//  Just to normalize the heck out of it.  I'm not really interested in if the fixture files have non-normative whitespace in them.
			var fixtureExpectNormBuf bytes.Buffer
			for _, line := range bytes.Split(fixtureExpect, []byte{'\n'}) {
				if len(line) == 0 {
					continue
				}
				nb := basicnode.Prototype.Any.NewBuilder()
				if err := json.Decode(nb, bytes.NewReader(line)); err != nil {
					t.Errorf("failed to parse fixture visit descriptions: %s", err)
				}
				json.Encode(nb.Build(), &fixtureExpectNormBuf)
				fixtureExpectNormBuf.WriteByte('\n')
			}

			// Serialize our own visit logs now too.
			var visitLogString bytes.Buffer
			for _, logEnt := range visitLogs {
				json.Encode(logEnt, &visitLogString)
				visitLogString.WriteByte('\n')
			}

			// DIFF TIME.
			qt.Assert(t, visitLogString.String(), qt.CmpEquals(), fixtureExpectNormBuf.String())
		})
	}
}
