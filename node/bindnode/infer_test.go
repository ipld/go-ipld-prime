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
	"runtime/debug"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
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
				uint   Int
				float  Float
				string String
				bytes  Bytes
			}`,
		ptrType: (*struct {
			Bool   bool
			Int    int64
			Uint   uint32
			Float  float64
			String string
			Bytes  []byte
		})(nil),
		prettyDagJSON: `{
			"bool":   true,
			"bytes":  {"/": {"bytes": "34cd"}},
			"float":  12.5,
			"int":    3,
			"string": "foo",
			"uint":   50
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
				uintAsInt            EnumAsInt
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
			UintAsInt            uint16
		})(nil),
		prettyDagJSON: `{
			"intAsInt":             12,
			"stringAsInt":          10,
			"stringAsString":       "Maybe",
			"stringAsStringCustom": "Yes",
			"uintAsInt":            11
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
		{"Any_Node_Int", "Any", (*datamodel.Node)(nil), `23`},
		// TODO: reenable once we don't require pointers for nullable
		// {"Any_Pointer_Int", "{String: nullable Any}",
		// 	(*struct {
		// 		Keys   []string
		// 		Values map[string]datamodel.Node
		// 	})(nil), `{"x":3,"y":"bar","z":[2.3,4.5]}`},
		{"Map_String_Int", "{String:Int}", (*struct {
			Keys   []string
			Values map[string]int64
		})(nil), `{"x":3,"y":4}`},
	}

	// For each IPLD kind, we test a matrix of combinations for IPLD's optional
	// and nullable fields alongside pointer usage on the Go field side.
	modifiers := []struct {
		schemaField string // "", "optional", "nullable", "optional nullable"
		goPointers  int    // 0 (T), 1 (*T), 2 (**T)
	}{
		{"", 0},                  // regular IPLD field with Go's T
		{"", 1},                  // regular IPLD field with Go's *T
		{"optional", 0},          // optional IPLD field with Go's T (skipped unless T is nilable)
		{"optional", 1},          // optional IPLD field with Go's *T
		{"nullable", 0},          // nullable IPLD field with Go's T (skipped unless T is nilable)
		{"nullable", 1},          // nullable IPLD field with Go's *T
		{"optional nullable", 2}, // optional and nullable IPLD field with Go's **T
	}
	for _, kindTest := range kindTests {
		for _, modifier := range modifiers {
			// don't reuse range vars
			kindTest := kindTest
			modifier := modifier
			goFieldType := reflect.TypeOf(kindTest.fieldPtrType)
			switch modifier.goPointers {
			case 0:
				goFieldType = goFieldType.Elem() // dereference fieldPtrType
			case 1:
				// fieldPtrType already uses one pointer
			case 2:
				goFieldType = reflect.PtrTo(goFieldType) // dereference fieldPtrType
			}
			if modifier.schemaField != "" && !nilable(goFieldType.Kind()) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s-%dptr", kindTest.name, modifier.schemaField, modifier.goPointers), func(t *testing.T) {
				t.Parallel()

				var buf bytes.Buffer
				err := template.Must(template.New("").Parse(`
						type Root struct {
							field {{.Modifier}} {{.Type}}
						}`)).Execute(&buf,
					struct {
						Type, Modifier string
					}{kindTest.schemaType, modifier.schemaField})
				qt.Assert(t, err, qt.IsNil)
				schemaSrc := buf.String()
				t.Logf("IPLD schema: %s", schemaSrc)

				// *struct { Field {{.goFieldType}} }
				goType := reflect.Zero(reflect.PtrTo(reflect.StructOf([]reflect.StructField{
					{Name: "Field", Type: goFieldType},
				}))).Interface()
				t.Logf("Go type: %T", goType)

				ts, err := ipld.LoadSchemaBytes([]byte(schemaSrc))
				qt.Assert(t, err, qt.IsNil)
				schemaType := ts.TypeByName("Root")
				qt.Assert(t, schemaType, qt.Not(qt.IsNil))

				proto := bindnode.Prototype(goType, schemaType)
				wantEncodedBytes, err := json.Marshal(map[string]interface{}{"field": json.RawMessage(kindTest.fieldDagJSON)})
				qt.Assert(t, err, qt.IsNil)
				wantEncoded := string(wantEncodedBytes)

				node := dagjsonDecode(t, proto.Representation(), wantEncoded).(schema.TypedNode)

				encoded := dagjsonEncode(t, node.Representation())
				qt.Assert(t, encoded, qt.Equals, wantEncoded)

				// Assigning with the missing field should only work with optional.
				nb := proto.NewBuilder()
				err = dagjson.Decode(nb, strings.NewReader(`{}`))
				switch modifier.schemaField {
				case "optional", "optional nullable":
					qt.Assert(t, err, qt.IsNil)
					node := nb.Build()
					// The resulting node should be non-nil with a nil field.
					nodeVal := reflect.ValueOf(bindnode.Unwrap(node))
					qt.Assert(t, nodeVal.Elem().FieldByName("Field").IsNil(), qt.IsTrue)
				default:
					qt.Assert(t, err, qt.Not(qt.IsNil))
				}

				// Assigning with a null field should only work with nullable.
				nb = proto.NewBuilder()
				err = dagjson.Decode(nb, strings.NewReader(`{"field":null}`))
				switch modifier.schemaField {
				case "nullable", "optional nullable":
					qt.Assert(t, err, qt.IsNil)
					node := nb.Build()
					// The resulting node should be non-nil with a nil field.
					nodeVal := reflect.ValueOf(bindnode.Unwrap(node))
					if modifier.schemaField == "nullable" {
						qt.Assert(t, nodeVal.Elem().FieldByName("Field").IsNil(), qt.IsTrue)
					} else {
						qt.Assert(t, nodeVal.Elem().FieldByName("Field").Elem().IsNil(), qt.IsTrue)
					}
				default:
					qt.Assert(t, err, qt.Not(qt.IsNil))
				}
			})
		}
	}
}

func nilable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	default:
		return false
	}
}

func assembleAsKind(proto datamodel.NodePrototype, schemaType schema.Type, asKind datamodel.Kind) (ipld.Node, error) {
	nb := proto.NewBuilder()
	switch asKind {
	case datamodel.Kind_Bool:
		if err := nb.AssignBool(true); err != nil {
			return nil, err
		}
	case datamodel.Kind_Int:
		if err := nb.AssignInt(123); err != nil {
			return nil, err
		}
	case datamodel.Kind_Float:
		if err := nb.AssignFloat(12.5); err != nil {
			return nil, err
		}
	case datamodel.Kind_String:
		if err := nb.AssignString("foo"); err != nil {
			return nil, err
		}
	case datamodel.Kind_Bytes:
		if err := nb.AssignBytes([]byte("\x00bar")); err != nil {
			return nil, err
		}
	case datamodel.Kind_Link:
		someCid, err := cid.Decode("bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi")
		if err != nil {
			return nil, err
		}
		if err := nb.AssignLink(cidlink.Link{Cid: someCid}); err != nil {
			return nil, err
		}
	case datamodel.Kind_Map:
		asm, err := nb.BeginMap(-1)
		if err != nil {
			return nil, err
		}
		// First via AssembleKey.
		if err := asm.AssembleKey().AssignString("F1"); err != nil {
			return nil, err
		}
		if err := asm.AssembleValue().AssignInt(101); err != nil {
			return nil, err
		}
		// Then via AssembleEntry.
		entryAsm, err := asm.AssembleEntry("F2")
		if err != nil {
			return nil, err
		}
		if err := entryAsm.AssignInt(102); err != nil {
			return nil, err
		}

		// If this is a struct, using a missing field should error.
		if _, ok := schemaType.(*schema.TypeStruct); ok {
			if err := asm.AssembleKey().AssignString("MissingKey"); err != nil {
				return nil, err
			}
			if err := asm.AssembleValue().AssignInt(101); err == nil {
				return nil, fmt.Errorf("expected error on missing struct key")
			}
		}
		if err := asm.Finish(); err != nil {
			return nil, err
		}
	case datamodel.Kind_List:
		asm, err := nb.BeginList(-1)
		if err != nil {
			return nil, err
		}
		// Note that we want the list to have two integer elements,
		// which matches the map entries above,
		// so that the struct with tuple repr just works too.
		if err := asm.AssembleValue().AssignInt(101); err != nil {
			return nil, err
		}
		if err := asm.AssembleValue().AssignInt(102); err != nil {
			return nil, err
		}
		// If this is a struct, assembling one more tuple entry should error.
		if _, ok := schemaType.(*schema.TypeStruct); ok {
			if err := asm.AssembleValue().AssignInt(103); err == nil {
				return nil, fmt.Errorf("expected error on extra tuple entry")
			}
		}
		if err := asm.Finish(); err != nil {
			return nil, err
		}
	}
	node := nb.Build()
	if node == nil {
		// If we succeeded, node must never be nil.
		return nil, fmt.Errorf("built node is nil")
	}
	return node, nil
}

func useNodeAsKind(node datamodel.Node, asKind datamodel.Kind) error {
	if gotKind := node.Kind(); gotKind != asKind {
		// Return a dummy error to signal when the kind doesn't match.
		return datamodel.ErrWrongKind{MethodName: "TestKindMismatches_Dummy_Kind"}
	}
	// Just check that IsAbsent and IsNull don't panic, for now.
	_ = node.IsAbsent()
	_ = node.IsNull()

	proto := node.Prototype()
	if proto == nil {
		return fmt.Errorf("got a null Prototype")
	}

	// TODO: also check LookupByNode, LookupBySegment
	switch asKind {
	case datamodel.Kind_Bool:
		if _, err := node.AsBool(); err != nil {
			return err
		}
	case datamodel.Kind_Int:
		if _, err := node.AsInt(); err != nil {
			return err
		}
	case datamodel.Kind_Float:
		if _, err := node.AsFloat(); err != nil {
			return err
		}
	case datamodel.Kind_String:
		if _, err := node.AsString(); err != nil {
			return err
		}
	case datamodel.Kind_Bytes:
		if _, err := node.AsBytes(); err != nil {
			return err
		}
	case datamodel.Kind_Link:
		if _, err := node.AsLink(); err != nil {
			return err
		}
	case datamodel.Kind_Map:
		iter := node.MapIterator()
		if iter == nil {
			// Return a dummy error to signal whether iter is nil or not.
			return datamodel.ErrWrongKind{MethodName: "TestKindMismatches_Dummy_MapIterator"}
		}
		for !iter.Done() {
			_, _, err := iter.Next()
			if err != nil {
				return err
			}
		}

		// valid element
		entryNode, err := node.LookupByString("F1")
		if err != nil {
			return err
		}
		if err := useNodeAsKind(entryNode, datamodel.Kind_Int); err != nil {
			return err
		}

		// missing element
		_, missingErr := node.LookupByString("MissingKey")
		switch err := missingErr.(type) {
		case nil:
			return fmt.Errorf("lookup of a missing key succeeded")
		case datamodel.ErrNotExists: // expected for maps
		case schema.ErrInvalidKey: // expected for structs
		default:
			return err
		}

		switch l := node.Length(); l {
		case 2:
		case -1:
			// Return a dummy error to signal whether Length failed.
			return datamodel.ErrWrongKind{MethodName: "TestKindMismatches_Dummy_Length"}
		default:
			return fmt.Errorf("unexpected Length: %d", l)
		}
	case datamodel.Kind_List:
		iter := node.ListIterator()
		if iter == nil {
			// Return a dummy error to signal whether iter is nil or not.
			return datamodel.ErrWrongKind{MethodName: "TestKindMismatches_Dummy_ListIterator"}
		}
		for !iter.Done() {
			_, _, err := iter.Next()
			if err != nil {
				return err
			}
		}

		// valid element
		entryNode, err := node.LookupByIndex(1)
		if err != nil {
			return err
		}
		if err := useNodeAsKind(entryNode, datamodel.Kind_Int); err != nil {
			return err
		}

		// missing element
		_, missingErr := node.LookupByIndex(30)
		switch err := missingErr.(type) {
		case nil:
			return fmt.Errorf("lookup of a missing key succeeded")
		case datamodel.ErrNotExists: // expected for maps
		case schema.ErrInvalidKey: // expected for structs
		default:
			return err
		}

		switch l := node.Length(); l {
		case 2:
		case -1:
			// Return a dummy error to signal whether Length failed.
			return datamodel.ErrWrongKind{MethodName: "TestKindMismatches_Dummy_Length"}
		default:
			return fmt.Errorf("unexpected Length: %d", l)
		}
	}
	return nil
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
		{"Enum", `
			type Root enum {
				| Foo ("foo")
				| Bar ("bar")
				| Either
			}
		`},
		{"Union_Kinded_onlyString", `
			type Root union {
				| String string
			} representation kinded
		`},
		{"Union_Kinded_onlyList", `
			type Root union {
				| List list
			} representation kinded
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

	// TODO: also test for non-repr assemblers and nodes

	for _, kindTest := range kindTests {
		// don't reuse range vars
		kindTest := kindTest
		t.Run(kindTest.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					// Note that debug.Stack inside the recover will include the
					// stack trace for the original panic call.
					t.Errorf("caught panic:\n%v\n%s", r, debug.Stack())
				}
			}()

			ts, err := ipld.LoadSchemaBytes([]byte(kindTest.schemaSrc))
			qt.Assert(t, err, qt.IsNil)
			schemaType := ts.TypeByName("Root")
			qt.Assert(t, schemaType, qt.Not(qt.IsNil))

			// Note that the Go type is inferred.
			proto := bindnode.Prototype(nil, schemaType).Representation()

			reprBehaviorKind := schemaType.RepresentationBehavior()
			if reprBehaviorKind == datamodel.Kind_Invalid {
				// For now, this only applies to kinded unions.
				// We'll need to modify this when we test with Any.
				members := schemaType.(*schema.TypeUnion).Members()
				qt.Assert(t, members, qt.HasLen, 1)
				reprBehaviorKind = members[0].RepresentationBehavior()
			}
			qt.Assert(t, reprBehaviorKind, qt.Not(qt.Equals), datamodel.Kind_Invalid)

			for _, kind := range allKinds {
				_, err := assembleAsKind(proto, schemaType, kind)
				comment := qt.Commentf("assigned as %v", kind)
				// Assembling should succed iff we used the right kind.
				if kind == reprBehaviorKind {
					qt.Assert(t, err, qt.IsNil, comment)
				} else {
					qt.Assert(t, err, qt.Not(qt.IsNil), comment)
					qt.Assert(t, err, qt.ErrorAs, new(datamodel.ErrWrongKind), comment)
				}
			}

			node, err := assembleAsKind(proto, schemaType, reprBehaviorKind)
			qt.Assert(t, err, qt.IsNil)
			node = node.(schema.TypedNode).Representation()
			nodeKind := node.Kind()

			for _, kind := range allKinds {
				err := useNodeAsKind(node, kind)
				comment := qt.Commentf("used as %v", kind)
				// Using the node should succed iff we used the right kind.
				if kind == nodeKind {
					qt.Assert(t, err, qt.IsNil, comment)
				} else {
					qt.Assert(t, err, qt.Not(qt.IsNil), comment)
					qt.Assert(t, err, qt.ErrorAs, new(datamodel.ErrWrongKind), comment)
				}
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
		name:      "MapNullableAny",
		schemaSrc: `type Root {String:nullable Any}`,
		goodTypes: []interface{}{
			(*struct {
				Keys   []string
				Values map[string]*datamodel.Node
			})(nil),
			(*struct {
				Keys   []string
				Values map[string]datamodel.Node
			})(nil),
		},
		badTypes: []verifyBadType{
			{(*string)(nil), `.*type Root .* type string: kind mismatch;.*`},
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
				List   []string
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
				List   []string
				String string
			})(nil), `.*type Root .*: union members must be nilable`},
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

func TestRenameAssignNode(t *testing.T) {
	type Foo struct{ I int }

	ts, _ := ipld.LoadSchemaBytes([]byte(`
type Foo struct {
	I Int (rename "J")
}
`))
	FooProto := bindnode.Prototype((*Foo)(nil), ts.TypeByName("Foo"))

	// Decode straight into bindnode typed builder
	nb := FooProto.Representation().NewBuilder()
	err := dagjson.Decode(nb, bytes.NewReader([]byte(`{"J":100}`)))
	qt.Assert(t, err, qt.IsNil)
	nb.Build()

	// decode into basicnode builder
	nb = basicnode.Prototype.Any.NewBuilder()
	err = dagjson.Decode(nb, bytes.NewReader([]byte(`{"J":100}`)))
	qt.Assert(t, err, qt.IsNil)
	node := nb.Build()

	// AssignNode from the basicnode form
	nb = FooProto.Representation().NewBuilder()
	err = nb.AssignNode(node)
	qt.Assert(t, err, qt.IsNil)
	nb.Build()
}

func TestEmptyTypedAssignNode(t *testing.T) {
	type Foo struct {
		I string
		J string
		K int
	}
	type Foo1Optional struct {
		I string
		J string
		K *int
	}
	type Foo2Optional struct {
		I string
		J *string
		K *int
	}
	tupleSchema := `type Foo struct {
		I String
		J String
		K Int
	} representation tuple`
	tuple1OptionalSchema := `type Foo struct {
		I String
		J String
		K optional Int
	} representation tuple`
	tuple2OptionalSchema := `type Foo struct {
		I String
		J optional String
		K optional Int
	} representation tuple`

	testCases := map[string]struct {
		schema  string
		typ     interface{}
		dagJson string
		err     string
	}{
		"tuple": {
			schema:  tupleSchema,
			typ:     (*Foo)(nil),
			dagJson: `["","",0]`,
		},
		"tuple with 2 absents": {
			schema:  tupleSchema,
			typ:     (*Foo)(nil),
			dagJson: `[""]`,
			err:     "missing required fields: J,K",
		},
		"tuple with 1 optional": {
			schema:  tuple1OptionalSchema,
			typ:     (*Foo1Optional)(nil),
			dagJson: `["","",0]`,
		},
		"tuple with 1 optional and absent": {
			schema:  tuple1OptionalSchema,
			typ:     (*Foo1Optional)(nil),
			dagJson: `["",""]`,
		},
		"tuple with 1 optional and 2 absents": {
			schema:  tuple1OptionalSchema,
			typ:     (*Foo1Optional)(nil),
			dagJson: `[""]`,
			err:     "missing required fields: J",
		},
		"tuple with 2 optional": {
			schema:  tuple2OptionalSchema,
			typ:     (*Foo2Optional)(nil),
			dagJson: `["","",0]`,
		},
		"tuple with 2 optional and 1 absent": {
			schema:  tuple2OptionalSchema,
			typ:     (*Foo2Optional)(nil),
			dagJson: `["",""]`,
		},
		"tuple with 2 optional and 2 absent": {
			schema:  tuple2OptionalSchema,
			typ:     (*Foo2Optional)(nil),
			dagJson: `[""]`,
		},
		"tuple with 2 optional and 3 absent": {
			schema:  tuple2OptionalSchema,
			typ:     (*Foo2Optional)(nil),
			dagJson: `[]`,
			err:     "missing required fields: I",
		},
		"map": {
			schema: `type Foo struct {
				I String
				J String
				K Int
			} representation map
			`,
			typ:     (*Foo)(nil),
			dagJson: `{"I":"","J":"","K":0}`,
		},
	}

	for testCase, data := range testCases {
		t.Run(testCase, func(t *testing.T) {
			ts, _ := ipld.LoadSchemaBytes([]byte(data.schema))
			FooProto := bindnode.Prototype(data.typ, ts.TypeByName("Foo"))

			// decode an "empty" object into Foo, these are all default values
			nb := basicnode.Prototype.Any.NewBuilder()
			err := dagjson.Decode(nb, bytes.NewReader([]byte(data.dagJson)))
			qt.Assert(t, err, qt.IsNil)
			node := nb.Build()

			// AssignNode from the basicnode form
			nb = FooProto.Representation().NewBuilder()
			err = nb.AssignNode(node)
			if data.err == "" {
				qt.Assert(t, err, qt.IsNil)
			} else {
				qt.Assert(t, err, qt.ErrorMatches, data.err)
			}
			nb.Build()

			// make an "empty" form, although none of the fields are optional so we should end up with defaults
			nb = FooProto.Representation().NewBuilder()
			empty := nb.Build()

			// AssignNode from the representation of the "empty" form, which should pass through default values
			nb = FooProto.Representation().NewBuilder()
			err = nb.AssignNode(empty.(schema.TypedNode).Representation())
			qt.Assert(t, err, qt.IsNil)
		})
	}
}
