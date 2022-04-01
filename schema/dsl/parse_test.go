package schemadsl_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	ipldjson "github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"
	"github.com/warpfork/go-testmark"
	"gopkg.in/yaml.v2"

	qt "github.com/frankban/quicktest"
)

var update = flag.Bool("update", false, "update testdata files in-place")

func TestParseSchemaSchema(t *testing.T) {
	t.Parallel()

	inputSchema := "../../.ipld/specs/schemas/schema-schema.ipldsch"
	inputJSON := "../../.ipld/specs/schemas/schema-schema.ipldsch.json"

	src, err := ioutil.ReadFile(inputSchema)
	qt.Assert(t, err, qt.IsNil)

	srcJSON, err := ioutil.ReadFile(inputJSON)
	qt.Assert(t, err, qt.IsNil)

	testParse(t, string(src), string(srcJSON), func(updated string) {
		err := ioutil.WriteFile(inputJSON, []byte(updated), 0o777)
		qt.Assert(t, err, qt.IsNil)
	})
}

type yamlFixture struct {
	Schema          string
	Canonical       string `yaml:",omitempty"`
	Expected        string
	ExpectedParsed  interface{}        `yaml:",omitempty"`
	Blocks          []yamlFixtureBlock `yaml:",omitempty"`
	BadBlocks       []string           `yaml:"badBlocks,omitempty"`
	BadBlocksParsed []interface{}      `yaml:",omitempty"`
}

type yamlFixtureBlock struct {
	Actual         string      `yaml:",omitempty"`
	ActualParsed   interface{} `yaml:",omitempty"`
	Expected       string      `yaml:",omitempty"`
	ExpectedParsed interface{} `yaml:",omitempty"`
}

// Test the DSL parsing, as well as the DMT compile, and its normalization,
// against fixtures that are in testmark data hunks in markdown files in the ipld/ipld repo.
func TestFromTestmark(t *testing.T) {
	// Glob all the markdown files in the fixtures directory.  Most do contain fixtures.
	matches, err := filepath.Glob("../../.ipld/specs/schemas/fixtures/*.md")
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, matches, qt.Not(qt.HasLen), 0)

	for _, pth := range matches {
		pth := pth // do not reuse range var since we're about to make parallel use of that value.

		// Fixture names tend to be unique, but let's bracket them with the filename they came from anyway.
		t.Run(filepath.Base(pth), func(t *testing.T) {
			t.Parallel()
			doc, err := testmark.ReadFile(pth)
			qt.Assert(t, err, qt.IsNil)

			// Data hunks in these spec files are in "directories" with one schema each.
			doc.BuildDirIndex()
			for _, dir := range doc.DirEnt.ChildrenList {
				t.Run(dir.Name, func(t *testing.T) {
					// First, check if this is marked for skipping.
					// TODO eh, regexp?  how should this work?

					// Each "directory" can contain many pieces of data, but only two we care about for this test:
					//  - `schema.ipldsch` -- the DSL form of the schema.
					//  - `schema.dmt.json` -- the logically equivalent Data Model tree of that same schema, in JSON.
					dslBlob := dir.Children["schema.ipldsch"].Hunk.Body
					dmtBlob := dir.Children["schema.dmt.json"].Hunk.Body

					var sch *schemadmt.Schema
					t.Run("parseable", func(t *testing.T) {
						sch, err = schemadsl.ParseBytes(dslBlob)
						qt.Assert(t, err, qt.IsNil)
					})
					if t.Failed() {
						t.FailNow()
					}

					t.Run("compilable", func(t *testing.T) {
						var ts schema.TypeSystem
						ts.Init()
						err := schemadmt.Compile(&ts, sch)
						qt.Assert(t, err, qt.IsNil)
						qt.Assert(t, ts.Names(), qt.Not(qt.HasLen), 0)
					})

					t.Run("matches-fixture-dmt", func(t *testing.T) {
						node := bindnode.Wrap(sch, schemadmt.Type.Schema.Type())

						var buf bytes.Buffer
						err := ipldjson.Encode(node.Representation(), &buf)
						qt.Assert(t, err, qt.IsNil)
						qt.Assert(t, buf.String(), qt.Equals, string(dmtBlob))
						// TODO add update feature
						//  ... the existing one works on the schema-schema too, that's nice, not exactly sure how to port that
					})

					// TODO: ensure that doing a json codec decode results in the same Schema Go
					// value that we got by parsing the DSL.
				})
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	matches, err := filepath.Glob("../../.ipld/specs/schemas/tests/*.yml")
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, matches, qt.Not(qt.HasLen), 0)

	for _, ymlPath := range matches {
		ymlPath := ymlPath // do not reuse range var
		name := filepath.Base(ymlPath)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := ioutil.ReadFile(ymlPath)
			qt.Assert(t, err, qt.IsNil)

			var fixt yamlFixture
			err = yaml.Unmarshal(data, &fixt)
			qt.Assert(t, err, qt.IsNil)

			inJSON := strings.Replace(fixt.Expected, "  ", "\t", -1)

			testParse(t, fixt.Schema, inJSON, func(updated string) {
				updated = strings.Replace(updated, "\t", "  ", -1)
				fixt.Expected = updated

				data, err = yaml.Marshal(&fixt)
				qt.Assert(t, err, qt.IsNil)

				// Note that this will strip comments.
				// Probably don't commit its changes straight away.
				err = ioutil.WriteFile(ymlPath, data, 0777)
				qt.Assert(t, err, qt.IsNil)
			})
		})
	}
}

func testParse(t *testing.T, inSchema, inJSON string, updateFn func(string)) {
	t.Helper()

	sch, err := schemadsl.ParseBytes([]byte(inSchema))
	qt.Assert(t, err, qt.IsNil)

	// Ensure the parsed schema compiles as expected.
	{
		var ts schema.TypeSystem
		ts.Init()
		err := schemadmt.Compile(&ts, sch)
		qt.Assert(t, err, qt.IsNil)

		qt.Assert(t, ts.Names(), qt.Not(qt.HasLen), 0)
	}

	// Ensure we can encode the schema as the json codec,
	// and that it results in the same bytes as the ipldsch.json file.
	{
		node := bindnode.Wrap(sch, schemadmt.Type.Schema.Type())

		var buf bytes.Buffer
		err := ipldjson.Encode(node.Representation(), &buf)
		qt.Assert(t, err, qt.IsNil)

		// If we're updating, write to the file.
		// Otherwise, expect the files to be equal.
		got := buf.String()
		if *update {
			updateFn(got)
			return
		}
		qt.Assert(t, got, qt.Equals, inJSON,
			qt.Commentf("run 'go test -update' to write to the iplsch.json file"))
	}

	// TODO: ensure that doing a json codec decode results in the same Schema Go
	// value that we got by parsing the DSL.
}
