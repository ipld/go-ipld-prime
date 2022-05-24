package bindnode_test

import (
	"encoding/hex"
	"math"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
)

func TestEnumError(t *testing.T) {
	type Action string
	const (
		ActionPresent = Action("p")
		ActionMissing = Action("m")
	)
	type S struct{ Action Action }

	schema := `
		type S struct {
			Action Action
		} representation tuple
		type Action enum {
			| Present             ("p")
			| Missing             ("m")
		} representation string
 	`

	typeSystem, err := ipld.LoadSchemaBytes([]byte(schema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("S")

	node := bindnode.Wrap(&S{Action: ActionPresent}, schemaType).Representation()
	_, err = ipld.Encode(node, dagcbor.Encode)
	qt.Assert(t, err, qt.IsNotNil)
	qt.Assert(t, err.Error(), qt.Equals, `AsString: "p" is not a valid member of enum Action (bindnode works at the type level; did you mean "Present"?)`)
}

func TestSubNodeWalkAndUnwrap(t *testing.T) {
	type F struct {
		F bool
	}
	type B struct {
		B int
	}
	type A struct {
		A string
	}
	type S struct {
		F      F
		B      B
		A      *A
		Any    datamodel.Node
		Bytes  []byte
		String string
	}

	schema := `
		type F struct {
			F Bool
		} representation tuple
		type B struct {
			B Int
		} representation tuple
		type A struct {
			A String
		} representation tuple
		type S struct {
			F F
			B B
			A optional A
			Any Any
			Bytes Bytes
			String String
		} representation tuple
	`

	encodedHex := "8681f581186581624141fb4069466666666666430102036f636f6e7374616e7420737472696e67" // [[true],[101],["AA"],202.2,[]byte{1,2,3},"constant string"]
	byts := []byte{0, 1, 2, 3}
	const constStr = "constant string"

	expected := S{
		F:      F{true},
		B:      B{101},
		A:      &A{"AAA"[1:]},
		Any:    basicnode.NewFloat(202.2),
		Bytes:  byts[1:],
		String: constStr,
	}

	typeSystem, err := ipld.LoadSchemaBytes([]byte(schema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("S")

	verifyMap := func(node datamodel.Node) {
		mi := node.MapIterator()

		key, value, err := mi.Next()
		qt.Assert(t, err, qt.IsNil)

		str, err := key.AsString()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, str, qt.Equals, "F")

		typ := bindnode.Unwrap(value)
		instF, ok := typ.(*F)
		qt.Assert(t, ok, qt.IsTrue)
		qt.Assert(t, *instF, qt.Equals, F{true})

		key, value, err = mi.Next()
		qt.Assert(t, err, qt.IsNil)

		str, err = key.AsString()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, str, qt.Equals, "B")

		typ = bindnode.Unwrap(value)
		instB, ok := typ.(*B)
		qt.Assert(t, ok, qt.IsTrue)
		qt.Assert(t, *instB, qt.Equals, B{101})

		key, value, err = mi.Next()
		qt.Assert(t, err, qt.IsNil)

		str, err = key.AsString()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, str, qt.Equals, "A")

		typ = bindnode.Unwrap(value)
		instA, ok := typ.(*A)
		qt.Assert(t, ok, qt.IsTrue)
		qt.Assert(t, *instA, qt.Equals, A{"AA"})

		key, value, err = mi.Next()
		qt.Assert(t, err, qt.IsNil)

		str, err = key.AsString()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, str, qt.Equals, "Any")

		qt.Assert(t, ipld.DeepEqual(basicnode.NewFloat(202.2), value), qt.IsTrue)

		key, value, err = mi.Next()
		qt.Assert(t, err, qt.IsNil)

		str, err = key.AsString()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, str, qt.Equals, "Bytes")

		typ = bindnode.Unwrap(value)
		instByts, ok := typ.(*[]byte)
		qt.Assert(t, ok, qt.IsTrue)
		qt.Assert(t, *instByts, qt.DeepEquals, []byte{1, 2, 3})

		key, value, err = mi.Next()
		qt.Assert(t, err, qt.IsNil)

		str, err = key.AsString()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, str, qt.Equals, "String")

		typ = bindnode.Unwrap(value)
		instStr, ok := typ.(*string)
		qt.Assert(t, ok, qt.IsTrue)
		qt.Assert(t, *instStr, qt.DeepEquals, "constant string")
	}

	t.Run("decode", func(t *testing.T) {
		encoded, _ := hex.DecodeString(encodedHex)

		proto := bindnode.Prototype(&S{}, schemaType)

		node, err := ipld.DecodeUsingPrototype([]byte(encoded), dagcbor.Decode, proto)
		qt.Assert(t, err, qt.IsNil)

		typ := bindnode.Unwrap(node)
		instS, ok := typ.(*S)
		qt.Assert(t, ok, qt.IsTrue)

		qt.Assert(t, *instS, qt.DeepEquals, expected)

		verifyMap(node)
	})

	t.Run("encode", func(t *testing.T) {
		node := bindnode.Wrap(&expected, schemaType)

		byts, err := ipld.Encode(node, dagcbor.Encode)
		qt.Assert(t, err, qt.IsNil)

		qt.Assert(t, hex.EncodeToString(byts), qt.Equals, encodedHex)

		verifyMap(node)
	})
}

func TestUint64Struct(t *testing.T) {
	t.Run("in struct", func(t *testing.T) {
		type IntHolder struct {
			Int32  int32
			Int64  int64
			Uint64 uint64
		}
		schema := `
			type IntHolder struct {
				Int32 Int
				Int64 Int
				Uint64 Int
			}
		`

		maxExpectedHex := "a365496e7433321a7fffffff65496e7436341b7fffffffffffffff6655696e7436341bffffffffffffffff"
		maxExpected, err := hex.DecodeString(maxExpectedHex)
		qt.Assert(t, err, qt.IsNil)

		typeSystem, err := ipld.LoadSchemaBytes([]byte(schema))
		qt.Assert(t, err, qt.IsNil)
		schemaType := typeSystem.TypeByName("IntHolder")
		proto := bindnode.Prototype(&IntHolder{}, schemaType)

		node, err := ipld.DecodeUsingPrototype([]byte(maxExpected), dagcbor.Decode, proto)
		qt.Assert(t, err, qt.IsNil)

		typ := bindnode.Unwrap(node)
		inst, ok := typ.(*IntHolder)
		qt.Assert(t, ok, qt.IsTrue)

		qt.Assert(t, *inst, qt.DeepEquals, IntHolder{
			Int32:  math.MaxInt32,
			Int64:  math.MaxInt64,
			Uint64: math.MaxUint64,
		})

		node = bindnode.Wrap(inst, schemaType).Representation()
		byt, err := ipld.Encode(node, dagcbor.Encode)
		qt.Assert(t, err, qt.IsNil)

		qt.Assert(t, hex.EncodeToString(byt), qt.Equals, maxExpectedHex)
	})

	t.Run("plain", func(t *testing.T) {
		type IntHolder uint64
		schema := `type IntHolder int`

		maxExpectedHex := "1bffffffffffffffff"
		maxExpected, err := hex.DecodeString(maxExpectedHex)
		qt.Assert(t, err, qt.IsNil)

		typeSystem, err := ipld.LoadSchemaBytes([]byte(schema))
		qt.Assert(t, err, qt.IsNil)
		schemaType := typeSystem.TypeByName("IntHolder")
		proto := bindnode.Prototype((*IntHolder)(nil), schemaType)

		node, err := ipld.DecodeUsingPrototype([]byte(maxExpected), dagcbor.Decode, proto)
		qt.Assert(t, err, qt.IsNil)

		typ := bindnode.Unwrap(node)
		inst, ok := typ.(*IntHolder)
		qt.Assert(t, ok, qt.IsTrue)

		qt.Assert(t, *inst, qt.Equals, IntHolder(math.MaxUint64))

		node = bindnode.Wrap(inst, schemaType).Representation()
		byt, err := ipld.Encode(node, dagcbor.Encode)
		qt.Assert(t, err, qt.IsNil)

		qt.Assert(t, hex.EncodeToString(byt), qt.Equals, maxExpectedHex)
	})
}
