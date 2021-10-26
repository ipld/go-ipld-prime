package bindnode_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

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

				ts, err := ipld.LoadSchemaBytes([]byte(test.schemaSrc))
				qt.Assert(t, err, qt.IsNil)
				schemaType := ts.TypeByName("Root")
				qt.Assert(t, schemaType, qt.Not(qt.IsNil))

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

type verifyBadType struct {
	ptrType     interface{}
	panicRegexp string
}

type (
	namedBool    bool
	namedInt64   int64
	namedFloat64 float64
	namedString  string
	namedBytes   []byte
)

var verifyTests = []struct {
	name      string
	schemaSrc string

	goodTypes []interface{}
	badTypes  []verifyBadType
}{
	{
		name:      "Bool",
		schemaSrc: `type Root bool`,
		goodTypes: []interface{}{
			(*bool)(nil),
			(*namedBool)(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch`},
		},
	},
	{
		name:      "Int",
		schemaSrc: `type Root int`,
		goodTypes: []interface{}{
			(*int)(nil),
			(*namedInt64)(nil),
			(*int8)(nil),
			(*int16)(nil),
			(*int32)(nil),
			(*int64)(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch`},
		},
	},
	{
		name:      "Float",
		schemaSrc: `type Root float`,
		goodTypes: []interface{}{
			(*float64)(nil),
			(*namedFloat64)(nil),
			(*float32)(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch`},
		},
	},
	{
		name:      "String",
		schemaSrc: `type Root string`,
		goodTypes: []interface{}{
			(*string)(nil),
			(*namedString)(nil),
		},
		badTypes: []verifyBadType{
			{(*int)(nil), `.*type Root .* type int: kind mismatch`},
		},
	},
	{
		name:      "Bytes",
		schemaSrc: `type Root bytes`,
		goodTypes: []interface{}{
			(*[]byte)(nil),
			(*namedBytes)(nil),
			(*[]uint8)(nil), // alias of byte
		},
		badTypes: []verifyBadType{
			{(*int)(nil), `.*type Root .* type int: kind mismatch`},
			{(*[]int)(nil), `.*type Root .* type \[\]int: kind mismatch`},
		},
	},
	{
		name:      "List",
		schemaSrc: `type Root [String]`,
		goodTypes: []interface{}{
			(*[]string)(nil),
			(*[]namedString)(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch`},
			{(*[]int)(nil), `.*type String .* type int: kind mismatch`},
			{(*[3]string)(nil), `.*type Root .* type \[3\]string: kind mismatch`},
		},
	},
	{
		name: "Struct",
		schemaSrc: `type Root struct {
				int Int
			}`,
		goodTypes: []interface{}{
			(*struct{ Int int })(nil),
			(*struct{ Int namedInt64 })(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch`},
			{(*struct{ Int bool })(nil), `.*type Int .* type bool: kind mismatch`},
			{(*struct{ Int1, Int2 int })(nil), `.*type Root .* type struct {.*}: 2 vs 1 fields`},
		},
	},
	{
		name:      "Map",
		schemaSrc: `type Root {String:Int}`,
		goodTypes: []interface{}{
			(*struct {
				Keys   []string
				Values map[string]int
			})(nil),
			(*struct {
				Keys   []namedString
				Values map[namedString]namedInt64
			})(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch`},
			{(*struct{ Keys []string })(nil), `.*type Root .*: 1 vs 2 fields`},
			{(*struct{ Values map[string]int })(nil), `.*type Root .*: 1 vs 2 fields`},
			{(*struct {
				Keys   string
				Values map[string]int
			})(nil), `.*type Root .*: kind mismatch`},
			{(*struct {
				Keys   []string
				Values string
			})(nil), `.*type Root .*: kind mismatch`},
		},
	},
	{
		name: "Union",
		schemaSrc: `type Root union {
				| List_String list
				| String      string
			} representation kinded

			type List_String [String]
			`,
		goodTypes: []interface{}{
			(*struct {
				List   *[]string
				String *string
			})(nil),
			(*struct {
				List   *[]namedString
				String *namedString
			})(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch`},
			{(*struct{ List *[]string })(nil), `.*type Root .*: 1 vs 2 members`},
			{(*struct {
				List   *[]string
				String string
			})(nil), `.*type Root .*: union members must be pointers`},
			{(*struct {
				List   *[]string
				String *int
			})(nil), `.*type String .*: kind mismatch`},
		},
	},
}

func TestSchemaVerify(t *testing.T) {
	t.Parallel()

	for _, test := range verifyTests {
		test := test // don't reuse the range var

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ts, err := ipld.LoadSchemaBytes([]byte(test.schemaSrc))
			qt.Assert(t, err, qt.IsNil)
			schemaType := ts.TypeByName("Root")
			qt.Assert(t, schemaType, qt.Not(qt.IsNil))

			for _, ptrType := range test.goodTypes {
				proto := bindnode.Prototype(ptrType, schemaType)
				qt.Assert(t, proto, qt.Not(qt.IsNil))

				ptrVal := reflect.New(reflect.TypeOf(ptrType).Elem()).Interface()
				node := bindnode.Wrap(ptrVal, schemaType)
				qt.Assert(t, node, qt.Not(qt.IsNil))
			}

			for _, bad := range test.badTypes {
				qt.Check(t, func() { bindnode.Prototype(bad.ptrType, schemaType) },
					qt.PanicMatches, bad.panicRegexp)

				ptrVal := reflect.New(reflect.TypeOf(bad.ptrType).Elem()).Interface()
				qt.Check(t, func() { bindnode.Wrap(ptrVal, schemaType) },
					qt.PanicMatches, bad.panicRegexp)
			}
		})
	}
}
