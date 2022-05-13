package dagjson_test

import (
	"bytes"
	"encoding/hex"
	"math"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	nodetests "github.com/ipld/go-ipld-prime/node/tests"
)

var n = fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("plain").AssignString("olde string")
	na.AssembleEntry("map").CreateMap(2, func(na fluent.MapAssembler) {
		na.AssembleEntry("one").AssignInt(1)
		na.AssembleEntry("two").AssignInt(2)
	})
	na.AssembleEntry("list").CreateList(2, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignString("three")
		na.AssembleValue().AssignString("four")
	})
	na.AssembleEntry("nested").CreateMap(1, func(na fluent.MapAssembler) {
		na.AssembleEntry("deeper").CreateList(1, func(na fluent.ListAssembler) {
			na.AssembleValue().AssignString("things")
		})
	})
})
var nSorted = fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("list").CreateList(2, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignString("three")
		na.AssembleValue().AssignString("four")
	})
	na.AssembleEntry("map").CreateMap(2, func(na fluent.MapAssembler) {
		na.AssembleEntry("one").AssignInt(1)
		na.AssembleEntry("two").AssignInt(2)
	})
	na.AssembleEntry("nested").CreateMap(1, func(na fluent.MapAssembler) {
		na.AssembleEntry("deeper").CreateList(1, func(na fluent.ListAssembler) {
			na.AssembleValue().AssignString("things")
		})
	})
	na.AssembleEntry("plain").AssignString("olde string")
})
var serial = `{"list":["three","four"],"map":{"one":1,"two":2},"nested":{"deeper":["things"]},"plain":"olde string"}`

func TestRoundtrip(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagjson.Encode(n, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, serial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(serial)
		nb := basicnode.Prototype.Map.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, nb.Build(), nodetests.NodeContentEquals, nSorted)
	})
}

func TestRoundtripScalar(t *testing.T) {
	nb := basicnode.Prototype__String{}.NewBuilder()
	nb.AssignString("applesauce")
	simple := nb.Build()
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagjson.Encode(simple, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, `"applesauce"`)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(`"applesauce"`)
		nb := basicnode.Prototype__String{}.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, nb.Build(), nodetests.NodeContentEquals, simple)
	})
}

func TestInts(t *testing.T) {
	t.Run("max uint64", func(t *testing.T) {
		buf := strings.NewReader("18446744073709551615") // max uint64
		nb := basicnode.Prototype.Any.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		n := nb.Build()

		// the overflowed AsInt() int64 cast
		_, err = n.AsInt()
		qt.Assert(t, err.Error(), qt.Equals, "unsigned integer out of rage of int64 type")

		// get real, underlying value
		uin, ok := n.(datamodel.UintNode)
		qt.Assert(t, ok, qt.IsTrue)
		val, positive, err := uin.AsUint()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, val, qt.Equals, uint64(math.MaxUint64))
		qt.Assert(t, positive, qt.IsTrue)

		var byts bytes.Buffer
		err = dagjson.Encode(n, &byts)
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, hex.EncodeToString(byts.Bytes()), qt.Equals, "18446744073709551615")
	})

	t.Run("max int64", func(t *testing.T) {
		buf := strings.NewReader("9223372036854775807") // max int64
		nb := basicnode.Prototype.Any.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		n := nb.Build()

		ii, err := n.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, ii, qt.Equals, int64(math.MaxInt64))

		// get uint form
		uin, ok := n.(datamodel.UintNode)
		qt.Assert(t, ok, qt.IsTrue)
		val, positive, err := uin.AsUint()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, val, qt.Equals, uint64(math.MaxInt64))
		qt.Assert(t, positive, qt.IsTrue)
	})

	t.Run("min int64", func(t *testing.T) {
		buf := strings.NewReader("-9223372036854775808") // min int64
		nb := basicnode.Prototype.Any.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		n := nb.Build()

		ii, err := n.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, ii, qt.Equals, int64(math.MinInt64))

		// get uint form
		uin, ok := n.(datamodel.UintNode)
		qt.Assert(t, ok, qt.IsTrue)
		val, positive, err := uin.AsUint()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, val, qt.Equals, uint64(math.MaxInt64)+1)
		qt.Assert(t, positive, qt.IsFalse)
	})
}
