package fluent_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

var roundTripTestCases = []struct {
	name  string
	value interface{}
}{
	{name: "Number", value: int64(100)},
	{name: "String", value: "hi"},
	{name: "Bool", value: true},
	{name: "Bytes", value: []byte("hi")},
	{name: "Map", value: map[string]interface{}{"a": "1", "b": int64(2), "c": 3.14, "d": true}},
	{name: "Array", value: []interface{}{"a", "b", "c"}},
	{name: "Nil", value: nil},
}

func TestRoundTrip(t *testing.T) {
	for _, testCase := range roundTripTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := qt.New(t)
			n, err := fluent.Reflect(basicnode.Prototype.Any, testCase.value)
			c.Assert(err, qt.IsNil)
			out, err := fluent.ToInterface(n)
			c.Assert(err, qt.IsNil)
			c.Check(out, qt.DeepEquals, testCase.value)
		})
	}
}

func TestLink(t *testing.T) {
	c := qt.New(t)
	someCid, err := cid.Parse("bafybeihrqe2hmfauph5yfbd6ucv7njqpiy4tvbewlvhzjl4bhnyiu6h7pm")
	c.Assert(err, qt.IsNil)
	link := cidlink.Link{Cid: someCid}
	v, err := fluent.ToInterface(basicnode.NewLink(link))
	c.Assert(err, qt.IsNil)
	c.Assert(v.(cidlink.Link), qt.Equals, link)
}
