package bindnode_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
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
	Map  *struct {
		Keys   []string
		Values map[string]string
	}
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
				linkCID     Link
				linkGeneric Link
				linkImpl    Link
			}`,
		ptrType: (*struct {
			LinkCID     cid.Cid
			LinkGeneric datamodel.Link
			LinkImpl    cidlink.Link
		})(nil),
		prettyDagJSON: `{
			"linkCID":     {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"},
			"linkGeneric": {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"},
			"linkImpl":    {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}
		}`,
	},
	{
		name: "Any",
		// TODO: also Null
		schemaSrc: `type Root struct {
				anyNodeWithBool   Any
				anyNodeWithInt    Any
				anyNodeWithFloat  Any
				anyNodeWithString Any
				anyNodeWithBytes  Any
				anyNodeWithList   Any
				anyNodeWithMap    Any
				anyNodeWithLink   Any

				anyNodeBehindList [Any]
				anyNodeBehindMap  {String:Any}
			}`,
		ptrType: (*struct {
			AnyNodeWithBool   datamodel.Node
			AnyNodeWithInt    datamodel.Node
			AnyNodeWithFloat  datamodel.Node
			AnyNodeWithString datamodel.Node
			AnyNodeWithBytes  datamodel.Node
			AnyNodeWithList   datamodel.Node
			AnyNodeWithMap    datamodel.Node
			AnyNodeWithLink   datamodel.Node

			AnyNodeBehindList []datamodel.Node
			AnyNodeBehindMap  struct {
				Keys   []string
				Values map[string]datamodel.Node
			}
		})(nil),
		prettyDagJSON: `{
			"anyNodeBehindList": [12.5, {"x": false}],
			"anyNodeBehindMap":  {"x": 123, "y": [true, false]},
			"anyNodeWithBool":   true,
			"anyNodeWithBytes":  {"/": {"bytes": "34cd"}},
			"anyNodeWithFloat":  12.5,
			"anyNodeWithInt":    3,
			"anyNodeWithLink":   {"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"},
			"anyNodeWithList":   [3, 2, 1],
			"anyNodeWithMap":    {"a": "x", "b": "y"},
			"anyNodeWithString": "foo"
		}`,
	},
	{
		name: "Enums",
		schemaSrc: `type Root struct {
				stringAsString       EnumAsString
				stringAsStringCustom EnumAsString
				stringAsInt          EnumAsInt
				intAsInt             EnumAsInt
			}
			type EnumAsString enum {
				| Nope ("No")
				| Yep  ("Yes")
				| Maybe
			}
			type EnumAsInt enum {
				| Nope  ("10")
				| Yep   ("11")
				| Maybe ("12")
			} representation int`,
		ptrType: (*struct {
			StringAsString       string
			StringAsStringCustom string
			StringAsInt          string
			IntAsInt             int32
		})(nil),
		prettyDagJSON: `{
			"intAsInt":             12,
			"stringAsInt":          10,
			"stringAsString":       "Maybe",
			"stringAsStringCustom": "Yes"
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
			onlySchema := onlySchema // don't reuse the range var
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

				ptrType := test.ptrType // don't write to the shared test value
				if onlySchema {
					ptrType = nil
				}
				proto := bindnode.Prototype(ptrType, schemaType)

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

func TestPrototypePointerCombinations(t *testing.T) {
	t.Parallel()

	// TODO: Null
	// TODO: cover more schema types and repr strategies.
	// Some of them are still using w.val directly without "nonPtr" calls.
	kindTests := []struct {
		name         string
		schemaType   string
		fieldPtrType interface{}
		fieldDagJSON string
	}{
		{"Bool", "Bool", (*bool)(nil), `true`},
		{"Int", "Int", (*int64)(nil), `23`},
		{"Float", "Float", (*float64)(nil), `34.5`},
		{"String", "String", (*string)(nil), `"foo"`},
		{"Bytes", "Bytes", (*[]byte)(nil), `{"/": {"bytes": "34cd"}}`},
		{"Link_CID", "Link", (*cid.Cid)(nil), `{"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}`},
		{"Link_Impl", "Link", (*cidlink.Link)(nil), `{"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}`},
		{"Link_Generic", "Link", (*datamodel.Link)(nil), `{"/": "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}`},
		{"List_String", "[String]", (*[]string)(nil), `["foo", "bar"]`},
		{"Map_String_Int", "{String:Int}", (*struct {
			Keys   []string
			Values map[string]int64
		})(nil), `{"x":3,"y":4}`},
	}

	for _, kindTest := range kindTests {
		for _, modifier := range []string{"", "optional", "nullable"} {
			// don't reuse range vars
			kindTest := kindTest
			modifier := modifier
			t.Run(fmt.Sprintf("%s/%s", kindTest.name, modifier), func(t *testing.T) {
				t.Parallel()

				var buf bytes.Buffer
				err := template.Must(template.New("").Parse(`
				type Root struct {
					field {{.Modifier}} {{.Type}}
				}`)).Execute(&buf, struct {
					Type, Modifier string
				}{kindTest.schemaType, modifier})
				qt.Assert(t, err, qt.IsNil)
				schemaSrc := buf.String()

				// *struct { Field {{.fieldPtrType}} }
				ptrType := reflect.Zero(reflect.PtrTo(reflect.StructOf([]reflect.StructField{
					{Name: "Field", Type: reflect.TypeOf(kindTest.fieldPtrType)},
				}))).Interface()

				ts, err := ipld.LoadSchemaBytes([]byte(schemaSrc))
				qt.Assert(t, err, qt.IsNil)
				schemaType := ts.TypeByName("Root")
				qt.Assert(t, schemaType, qt.Not(qt.IsNil))

				proto := bindnode.Prototype(ptrType, schemaType)
				wantEncodedBytes, err := json.Marshal(map[string]interface{}{"field": json.RawMessage(kindTest.fieldDagJSON)})
				qt.Assert(t, err, qt.IsNil)
				wantEncoded := string(wantEncodedBytes)

				node := dagjsonDecode(t, proto.Representation(), wantEncoded).(schema.TypedNode)

				encoded := dagjsonEncode(t, node.Representation())
				qt.Assert(t, encoded, qt.Equals, wantEncoded)

				// Assigning with the missing field should only work with optional.
				nb := proto.NewBuilder()
				err = dagjson.Decode(nb, strings.NewReader(`{}`))
				if modifier == "optional" {
					qt.Assert(t, err, qt.IsNil)
					node := nb.Build()
					// The resulting node should be non-nil with a nil field.
					nodeVal := reflect.ValueOf(bindnode.Unwrap(node))
					qt.Assert(t, nodeVal.Elem().FieldByName("Field").IsNil(), qt.IsTrue)
				} else {
					qt.Assert(t, err, qt.Not(qt.IsNil))
				}

				// Assigning with a null field should only work with nullable.
				nb = proto.NewBuilder()
				err = dagjson.Decode(nb, strings.NewReader(`{"field":null}`))
				if modifier == "nullable" {
					qt.Assert(t, err, qt.IsNil)
					node := nb.Build()
					// The resulting node should be non-nil with a nil field.
					nodeVal := reflect.ValueOf(bindnode.Unwrap(node))
					qt.Assert(t, nodeVal.Elem().FieldByName("Field").IsNil(), qt.IsTrue)
				} else {
					qt.Assert(t, err, qt.Not(qt.IsNil))
				}
			})
		}
	}
}

func TestKindMismatches(t *testing.T) {
	t.Parallel()

	kindTests := []struct {
		name      string
		schemaSrc string
	}{
		{"Bool", "type Root bool"},
		{"Int", "type Root int"},
		{"Float", "type Root float"},
		{"String", "type Root string"},
		{"Bytes", "type Root bytes"},
		{"Map", `
			type Root {String:Int}
		`},
		{"Struct", `
			type Root struct {
				F1 Int
				F2 Int
			}
		`},
		{"Struct_Tuple", `
			type Root struct {
				F1 Int
				F2 Int
			} representation tuple
		`},
		// TODO: more schema types and repr strategies
	}

	allKinds := []datamodel.Kind{
		// datamodel.Kind_Null, TODO
		datamodel.Kind_Bool,
		datamodel.Kind_Int,
		datamodel.Kind_Float,
		datamodel.Kind_String,
		datamodel.Kind_Bytes,
		datamodel.Kind_Link,
		datamodel.Kind_Map,
		datamodel.Kind_List,
	}

	someCid, err := cid.Decode("bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi")
	qt.Assert(t, err, qt.IsNil)
	assembleKind := func(proto datamodel.NodePrototype, kind datamodel.Kind) error {
		nb := proto.NewBuilder()
		switch kind {
		case datamodel.Kind_Bool:
			if err := nb.AssignBool(true); err != nil {
				return err
			}
		case datamodel.Kind_Int:
			if err := nb.AssignInt(123); err != nil {
				return err
			}
		case datamodel.Kind_Float:
			if err := nb.AssignFloat(12.5); err != nil {
				return err
			}
		case datamodel.Kind_String:
			if err := nb.AssignString("foo"); err != nil {
				return err
			}
		case datamodel.Kind_Bytes:
			if err := nb.AssignBytes([]byte("\x00bar")); err != nil {
				return err
			}
		case datamodel.Kind_Link:
			if err := nb.AssignLink(cidlink.Link{Cid: someCid}); err != nil {
				return err
			}
		case datamodel.Kind_Map:
			asm, err := nb.BeginMap(-1)
			if err != nil {
				return err
			}
			// First via AssembleKey.
			if err := asm.AssembleKey().AssignString("F1"); err != nil {
				return err
			}
			if err := asm.AssembleValue().AssignInt(101); err != nil {
				return err
			}
			// Then via AssembleEntry.
			entryAsm, err := asm.AssembleEntry("F2")
			if err != nil {
				return err
			}
			if err := entryAsm.AssignInt(102); err != nil {
				return err
			}
			if err := asm.Finish(); err != nil {
				return err
			}
		case datamodel.Kind_List:
			asm, err := nb.BeginList(-1)
			if err != nil {
				return err
			}
			if err := asm.AssembleValue().AssignInt(101); err != nil {
				return err
			}
			if err := asm.AssembleValue().AssignInt(102); err != nil {
				return err
			}
			if err := asm.Finish(); err != nil {
				return err
			}
		}
		node := nb.Build()
		// If we succeeded, node must never be nil.
		qt.Assert(t, node, qt.Not(qt.IsNil))
		return nil
	}

	// TODO: also test for non-repr assemblers and nodes

	for _, kindTest := range kindTests {
		// don't reuse range vars
		kindTest := kindTest
		t.Run(kindTest.name, func(t *testing.T) {
			t.Parallel()

			ts, err := ipld.LoadSchemaBytes([]byte(kindTest.schemaSrc))
			qt.Assert(t, err, qt.IsNil)
			schemaType := ts.TypeByName("Root")
			qt.Assert(t, schemaType, qt.Not(qt.IsNil))

			// Note that the Go type is inferred.
			proto := bindnode.Prototype(nil, schemaType).Representation()

			actualKind := schemaType.RepresentationBehavior()

			for _, kind := range allKinds {
				err := assembleKind(proto, kind)
				comment := qt.Commentf("Assign of %v", kind)
				// Assembling should succed iff we used the right kind.
				if kind == actualKind {
					qt.Assert(t, err, qt.IsNil, comment)
				} else {
					qt.Assert(t, err, qt.Not(qt.IsNil), comment)
					qt.Assert(t, err, qt.ErrorAs, new(datamodel.ErrWrongKind), comment)
				}

				// TODO: check AsT methods just like AssignT
				// TODO: also check valid methods per kind, like Length
			}
		})
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
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
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
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
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
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
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
			{(*int)(nil), `.*type Root .* type int: kind mismatch;.*`},
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
			{(*int)(nil), `.*type Root .* type int: kind mismatch;.*`},
			{(*[]int)(nil), `.*type Root .* type \[\]int: kind mismatch;.*`},
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
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
			{(*[]int)(nil), `.*type String .* type int: kind mismatch;.*`},
			{(*[3]string)(nil), `.*type Root .* type \[3\]string: kind mismatch;.*`},
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
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
			{(*struct{ Int bool })(nil), `.*type Int .* type bool: kind mismatch;.*`},
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
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
			{(*struct{ Keys []string })(nil), `.*type Root .*: 1 vs 2 fields`},
			{(*struct{ Values map[string]int })(nil), `.*type Root .*: 1 vs 2 fields`},
			{(*struct {
				Keys   string
				Values map[string]int
			})(nil), `.*type Root .*: kind mismatch;.*`},
			{(*struct {
				Keys   []string
				Values string
			})(nil), `.*type Root .*: kind mismatch;.*`},
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
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
			{(*struct{ List *[]string })(nil), `.*type Root .*: 1 vs 2 members`},
			{(*struct {
				List   *[]string
				String string
			})(nil), `.*type Root .*: union members must be pointers`},
			{(*struct {
				List   *[]string
				String *int
			})(nil), `.*type String .*: kind mismatch;.*`},
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

func TestProduceGoTypes(t *testing.T) {
	t.Parallel()

	for _, test := range prototypeTests {
		test := test // don't reuse the range var

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ts, err := ipld.LoadSchemaBytes([]byte(test.schemaSrc))
			qt.Assert(t, err, qt.IsNil)

			// Include a package line and the datamodel import.
			buf := new(bytes.Buffer)
			fmt.Fprintln(buf, `package p`)
			fmt.Fprintln(buf, `import "github.com/ipld/go-ipld-prime/datamodel"`)
			fmt.Fprintln(buf, `var _ datamodel.Link // always used`)
			err = bindnode.ProduceGoTypes(buf, ts)
			qt.Assert(t, err, qt.IsNil)

			// Ensure that the output builds, i.e. typechecks.
			genPath := filepath.Join(t.TempDir(), "gen.go")
			err = ioutil.WriteFile(genPath, buf.Bytes(), 0o666)
			qt.Assert(t, err, qt.IsNil)

			out, err := exec.Command("go", "build", genPath).CombinedOutput()
			qt.Assert(t, err, qt.IsNil, qt.Commentf("output: %s", out))

			// TODO: check that the generated types are compatible with the schema.
		})
	}
}
