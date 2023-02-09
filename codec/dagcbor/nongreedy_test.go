package dagcbor

import (
	"bytes"
	"encoding/hex"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestNonGreedy(t *testing.T) {
	// same as JSON version of this test: {"a": 1}{"b": 2}
	buf, err := hex.DecodeString("a1616101a1616202")
	qt.Assert(t, err, qt.IsNil)
	r := bytes.NewReader(buf)
	opts := DecodeOptions{
		DontParseBeyondEnd: true,
	}

	// first object
	nb1 := basicnode.Prototype.Map.NewBuilder()
	err = opts.Decode(nb1, r)
	qt.Assert(t, err, qt.IsNil)
	expected, err := qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "a", qp.Int(1))
	})
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, ipld.DeepEqual(nb1.Build(), expected), qt.IsTrue)

	// second object
	nb2 := basicnode.Prototype.Map.NewBuilder()
	err = opts.Decode(nb2, r)
	qt.Assert(t, err, qt.IsNil)
	expected, err = qp.BuildMap(basicnode.Prototype.Any, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "b", qp.Int(2))
	})
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, ipld.DeepEqual(nb2.Build(), expected), qt.IsTrue)
}
