package bindnode_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"
)

// TODO: tests where an IPLD schema and Go type do not match

type anyScalar struct {
	Bool   *bool
	Int    *int64
	Float  *float64
	String *string
	Bytes  *[]byte
	Link   *datamodel.Link
}

type anyRecursive struct {
	List *[]string
	Map  *map[string]string
}

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
			LinkGeneric datamodel.Link
			LinkCID     cid.Cid
		})(nil),
		prettyDagJSON: `{
			"linkCID":     {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"},
			"linkGeneric": {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}
		}`,
	},
	{
		name: "ScalarKindedUnions",
		// TODO: should we use an "Any" type from the prelude?
		schemaSrc: `type Root struct {
				boolAny   AnyScalar
				intAny    AnyScalar
				floatAny  AnyScalar
				stringAny AnyScalar
				bytesAny  AnyScalar
				linkAny   AnyScalar
			}

			type AnyScalar union {
				| Bool   bool
				| Int    int
				| Float  float
				| String string
				| Bytes  bytes
				| Link   link
			} representation kinded`,
		ptrType: (*struct {
			BoolAny   anyScalar
			IntAny    anyScalar
			FloatAny  anyScalar
			StringAny anyScalar
			BytesAny  anyScalar
			LinkAny   anyScalar
		})(nil),
		prettyDagJSON: `{
			"boolAny":   true,
			"bytesAny":  {"/": {"bytes": "34cd"}},
			"floatAny":  12.5,
			"intAny":    3,
			"linkAny":   {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"},
			"stringAny": "foo"
		}`,
	},
	{
		name: "RecursiveKindedUnions",
		// TODO: should we use an "Any" type from the prelude?
		// Especially since we use String map/list element types.
		// TODO: use inline map/list defs once schema and dsl+dmt support it.
		schemaSrc: `type Root struct {
				listAny AnyRecursive
				mapAny  AnyRecursive
			}

			type List_String [String]
			type Map_String {String:String}

			type AnyRecursive union {
				| List_String list
				| Map_String  map
			} representation kinded`,
		ptrType: (*struct {
			ListAny anyRecursive
			MapAny  anyRecursive
		})(nil),
		prettyDagJSON: `{
			"listAny": ["foo", "bar"],
			"mapAny":  {"a": "x", "b": "y"}
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

func dagjsonEncode(t *testing.T, node datamodel.Node) string {
	var sb strings.Builder
	err := dagjson.Encode(node, &sb)
	qt.Assert(t, err, qt.IsNil)
	return sb.String()
}

func dagjsonDecode(t *testing.T, proto datamodel.NodePrototype, src string) datamodel.Node {
	nb := proto.NewBuilder()
	err := dagjson.Decode(nb, strings.NewReader(src))
	qt.Assert(t, err, qt.IsNil)
	return nb.Build()
}

func TestPrototype(t *testing.T) {
	t.Parallel()

	for _, test := range prototypeTests {
		test := test // don't reuse the range var

		for _, onlySchema := range []bool{false, true} {
			suffix := ""
			if onlySchema {
				suffix = "_onlySchema"
			}
			t.Run(test.name+suffix, func(t *testing.T) {
				t.Parallel()

				schemaType := loadSchema(t, test.schemaSrc)

				if onlySchema {
					test.ptrType = nil
				}
				proto := bindnode.Prototype(test.ptrType, schemaType)

				wantEncoded := compactJSON(t, test.prettyDagJSON)
				node := dagjsonDecode(t, proto.Representation(), wantEncoded).(schema.TypedNode)
				// TODO: assert node type matches ptrType

				encoded := dagjsonEncode(t, node.Representation())
				qt.Assert(t, encoded, qt.Equals, wantEncoded)

				// Verify that doing a dag-json encode of the non-repr node works.
				_ = dagjsonEncode(t, node)
			})
		}
	}
}
