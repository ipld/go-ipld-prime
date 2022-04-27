package schemadsl_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	ipldjson "github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"
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

			testParse(t, fixt.Schema, fixt.Expected, func(updated string) {
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

	inJSON = strings.Replace(inJSON, "  ", "\t", -1) // fix non-tab indenting
	crre := regexp.MustCompile(`\r?\n`)
	inJSON = crre.ReplaceAllString(inJSON, "\n") // fix windows carriage-return

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
