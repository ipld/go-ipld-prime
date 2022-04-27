package schemadmt_test

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	ipldjson "github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"

	qt "github.com/frankban/quicktest"
)

func TestRoundtripSchemaSchema(t *testing.T) {
	t.Parallel()

	input := "../../.ipld/specs/schemas/schema-schema.ipldsch.json"

	src, err := ioutil.ReadFile(input)
	qt.Assert(t, err, qt.IsNil)
	testRoundtrip(t, string(src), func(updated string) {
		err := ioutil.WriteFile(input, []byte(updated), 0o777)
		qt.Assert(t, err, qt.IsNil)
	})
}

func testRoundtrip(t *testing.T, want string, updateFn func(string)) {
	t.Helper()

	crre := regexp.MustCompile(`\r?\n`)
	want = crre.ReplaceAllString(want, "\n")
	nb := schemadmt.Type.Schema.Representation().NewBuilder()
	err := ipldjson.Decode(nb, strings.NewReader(want))
	qt.Assert(t, err, qt.IsNil)
	node := nb.Build().(schema.TypedNode)

	// Ensure the decoded schema compiles as expected.
	{
		sch := bindnode.Unwrap(node).(*schemadmt.Schema)

		var ts schema.TypeSystem
		ts.Init()
		err := schemadmt.Compile(&ts, sch)
		qt.Assert(t, err, qt.IsNil)

		typeStruct := ts.TypeByName("TypeDefnStruct")
		if typeStruct == nil {
			t.Fatal("TypeStruct not found")
		}
	}

	// Ensure we can re-encode the schema as dag-json,
	// and that it results in the same bytes as prettified by encoding/json.
	{
		var buf bytes.Buffer
		err := ipldjson.Encode(node.Representation(), &buf)
		qt.Assert(t, err, qt.IsNil)

		got := buf.String()
		qt.Assert(t, got, qt.Equals, want)
	}

	// For the sake of completeness, check that we can encode the non-repr node.
	// This just ensures we don't panic or error.
	{
		var buf bytes.Buffer
		err := ipldjson.Encode(node, &buf)
		qt.Assert(t, err, qt.IsNil)
	}
}
