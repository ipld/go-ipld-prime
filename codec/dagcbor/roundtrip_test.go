package dagcbor

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"math"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	cid "github.com/ipfs/go-cid"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
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
	na.AssembleEntry("map").CreateMap(2, func(na fluent.MapAssembler) {
		na.AssembleEntry("one").AssignInt(1)
		na.AssembleEntry("two").AssignInt(2)
	})
	na.AssembleEntry("list").CreateList(2, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignString("three")
		na.AssembleValue().AssignString("four")
	})
	na.AssembleEntry("plain").AssignString("olde string")
	na.AssembleEntry("nested").CreateMap(1, func(na fluent.MapAssembler) {
		na.AssembleEntry("deeper").CreateList(1, func(na fluent.ListAssembler) {
			na.AssembleValue().AssignString("things")
		})
	})
})
var serial = "\xa4cmap\xa2cone\x01ctwo\x02dlist\x82ethreedfoureplainkolde stringfnested\xa1fdeeper\x81fthings"

func TestRoundtrip(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encode(n, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, serial)
	})
	t.Run("length", func(t *testing.T) {
		length, err := EncodedLength(n)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, length, qt.Equals, int64(len(serial)))
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(serial)
		nb := basicnode.Prototype.Map.NewBuilder()
		err := Decode(nb, buf)
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
		err := Encode(simple, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, `japplesauce`)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(`japplesauce`)
		nb := basicnode.Prototype__String{}.NewBuilder()
		err := Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, nb.Build(), nodetests.NodeContentEquals, simple)
	})
}

func TestRoundtripLinksAndBytes(t *testing.T) {
	lnk := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    0x71,
		MhType:   0x13,
		MhLength: 4,
	}}.BuildLink([]byte{1, 2, 3, 4}) // dummy value, content does not matter to this test.

	var linkByteNode = fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
		nva := na.AssembleEntry("Link")
		nva.AssignLink(lnk)
		nva = na.AssembleEntry("Bytes")
		bytes := make([]byte, 100)
		_, _ = rand.Read(bytes)
		nva.AssignBytes(bytes)
	})

	buf := bytes.Buffer{}
	err := Encode(linkByteNode, &buf)
	qt.Assert(t, err, qt.IsNil)
	nb := basicnode.Prototype.Map.NewBuilder()
	err = Decode(nb, &buf)
	qt.Assert(t, err, qt.IsNil)
	reconstructed := nb.Build()
	qt.Check(t, reconstructed, nodetests.NodeContentEquals, linkByteNode)
}

func TestInts(t *testing.T) {
	data := []struct {
		name      string
		hex       string
		value     uint64
		intValue  int64
		intErr    string
		decodeErr string
	}{
		{"max uint64", "1bffffffffffffffff", math.MaxUint64, 0, "unsigned integer out of range of int64 type", ""},
		{"max int64", "1b7fffffffffffffff", math.MaxInt64, math.MaxInt64, "", ""},
		{"1", "01", 1, 1, "", ""},
		{"0", "00", 0, 0, "", ""},
		{"-1", "20", 0, -1, "", ""},
		{"min int64", "3b7fffffffffffffff", 0, math.MinInt64, "", ""},
		{"~min uint64", "3bfffffffffffffffe", 0, 0, "", "cbor: negative integer out of rage of int64 type"},
		// TODO: 3bffffffffffffffff isn't properly handled by refmt, it's coerced to zero
		// MaxUint64 gets overflowed here: https://github.com/polydawn/refmt/blob/30ac6d18308e584ca6a2e74ba81475559db94c5f/cbor/cborDecoderTerminals.go#L75
	}

	for _, td := range data {
		t.Run(td.name, func(t *testing.T) {
			buf, err := hex.DecodeString(td.hex) // max uint64
			qt.Assert(t, err, qt.IsNil)
			nb := basicnode.Prototype.Any.NewBuilder()
			err = Decode(nb, bytes.NewReader(buf))
			if td.decodeErr != "" {
				qt.Assert(t, err, qt.IsNotNil)
				qt.Assert(t, err.Error(), qt.Equals, td.decodeErr)
				return
			}
			qt.Assert(t, err, qt.IsNil)
			n := nb.Build()

			ii, err := n.AsInt()
			if td.intErr != "" {
				qt.Assert(t, err.Error(), qt.Equals, td.intErr)
			} else {
				qt.Assert(t, err, qt.IsNil)
				qt.Assert(t, ii, qt.Equals, int64(td.intValue))
			}

			// if the number is outside of the positive int64 range, we should be able
			// to access it as a UintNode and be able to access the full int64 range
			uin, ok := n.(datamodel.UintNode)
			if td.value <= math.MaxInt64 {
				qt.Assert(t, ok, qt.IsFalse)
			} else {
				qt.Assert(t, ok, qt.IsTrue)
				val, err := uin.AsUint()
				qt.Assert(t, err, qt.IsNil)
				qt.Assert(t, val, qt.Equals, uint64(td.value))
			}

			var byts bytes.Buffer
			err = Encode(n, &byts)
			qt.Assert(t, err, qt.IsNil)
			qt.Assert(t, hex.EncodeToString(byts.Bytes()), qt.Equals, td.hex)
		})
	}
}
