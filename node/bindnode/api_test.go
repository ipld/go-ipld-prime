package bindnode_test

import (
	"encoding/hex"
	"math"
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/node/bindnode"

	qt "github.com/frankban/quicktest"
)

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
