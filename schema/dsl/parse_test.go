package schemadsl_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"

	ipldjson "github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"

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

		typeStruct := ts.TypeByName("TypeDefnStruct")
		if typeStruct == nil {
			t.Fatal("TypeStruct not found")
		}
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
