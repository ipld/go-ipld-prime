package dagjson_test

import (
	"bytes"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	nodetests "github.com/ipld/go-ipld-prime/node/tests"
)

var byteNode = fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("plain").AssignString("olde string")
	na.AssembleEntry("bytes").AssignBytes([]byte("deadbeef"))
})
var byteNodeSorted = fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("bytes").AssignBytes([]byte("deadbeef"))
	na.AssembleEntry("plain").AssignString("olde string")
})
var byteSerial = `{"bytes":{"/":{"bytes":"ZGVhZGJlZWY"}},"plain":"olde string"}`

func TestRoundtripBytes(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagjson.Encode(byteNode, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, byteSerial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(byteSerial)
		nb := basicnode.Prototype.Map.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, nb.Build(), nodetests.NodeContentEquals, byteNodeSorted)
	})
}

var encapsulatedNode = fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
	na.AssembleEntry("/").CreateMap(1, func(sa fluent.MapAssembler) {
		sa.AssembleEntry("bytes").AssignBytes([]byte("deadbeef"))
	})
})
var encapsulatedSerial = `{"/":{"bytes":{"/":{"bytes":"ZGVhZGJlZWY"}}}}`

func TestEncapsulatedBytes(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagjson.Encode(encapsulatedNode, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, encapsulatedSerial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(encapsulatedSerial)
		nb := basicnode.Prototype.Map.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, nb.Build(), nodetests.NodeContentEquals, encapsulatedNode)
	})
}

var withPadding = `{"/": {"bytes": "Bxrk96XO8cwr3hrcL4VeWtVdYudzHv47BbBl7CesWvmjRrRPOLZp9Ukg6sivn5Nqg4V5X2w43mk4Ppuzr+M+DA=="}}`

func TestPaddedBytes(t *testing.T) {
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(withPadding)
		nb := basicnode.Prototype.Bytes.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
	})
}
