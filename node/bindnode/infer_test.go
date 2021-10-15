package bindnode_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"
)

// TODO: tests where an IPLD schema and Go type do not match

var prototypeTests = []struct {
	name      string
	schemaSrc string
	ptrType   interface{}

	// Prettified for maintainability, valid DAG-JSON when compacted.
	prettyDagJSON string
}{
	{
		name: "Scalars",
		schemaSrc: `type Root struct {
			bool   Bool
			int    Int
			float  Float
			string String
			bytes  Bytes
		}`,
		ptrType: (*struct {
			Bool   bool
			Int    int64
			Float  float64
			String string
			Bytes  []byte
		})(nil),
		prettyDagJSON: `{
			"bool":   true,
			"bytes":  {"/": {"bytes": "34cd"}},
			"float":  12.5,
			"int":    3,
			"string": "foo"
		}`,
	},
	{
		name: "Links",
		schemaSrc: `type Root struct {
			linkGeneric Link
			linkCID     Link
		}`,
		ptrType: (*struct {
			LinkGeneric ipld.Link
			LinkCID     cid.Cid
		})(nil),
		prettyDagJSON: `{
			"linkCID":     {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"},
			"linkGeneric": {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}
		}`,
	},
}

func loadSchema(t *testing.T, src string) schema.Type {
	t.Helper()

	dmt, err := schemadsl.Parse("", strings.NewReader(src))
	qt.Assert(t, err, qt.IsNil)

	var ts schema.TypeSystem
	ts.Init()
	err = schemadmt.Compile(&ts, dmt)
	qt.Assert(t, err, qt.IsNil)

	typ := ts.TypeByName("Root")
	qt.Assert(t, typ, qt.Not(qt.IsNil))
	return typ
}

func compactJSON(t *testing.T, pretty string) string {
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(pretty))
	qt.Assert(t, err, qt.IsNil)
	return buf.String()
}

func dagjsonEncode(t *testing.T, node ipld.Node) string {
	var sb strings.Builder
	err := dagjson.Encode(node, &sb)
	qt.Assert(t, err, qt.IsNil)
	return sb.String()
}

func dagjsonDecode(t *testing.T, proto ipld.NodePrototype, src string) ipld.Node {
	nb := proto.NewBuilder()
	err := dagjson.Decode(nb, strings.NewReader(src))
	qt.Assert(t, err, qt.IsNil)
	return nb.Build()
}

func TestPrototype(t *testing.T) {
	t.Parallel()

	for _, test := range prototypeTests {
		test := test // don't reuse the range var

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			schemaType := loadSchema(t, test.schemaSrc)
			proto := bindnode.Prototype(test.ptrType, schemaType)

			wantEncoded := compactJSON(t, test.prettyDagJSON)
			node := dagjsonDecode(t, proto, wantEncoded)
			// TODO: assert node type matches ptrType

			encoded := dagjsonEncode(t, node)
			qt.Assert(t, encoded, qt.Equals, wantEncoded)

			// TODO: also check that just using the schema works?
		})
	}
}