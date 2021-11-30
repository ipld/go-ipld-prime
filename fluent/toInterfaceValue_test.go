package fluent_test

import (
	"encoding/json"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

var roundTripTestCases = []struct {
	desc  string
	value interface{}
}{
	{desc: "Number", value: 100},
	{desc: "String", value: "hi"},
	{desc: "Bool", value: "hi"},
	{desc: "Bytes", value: []byte("hi")},
	{desc: "Map", value: map[string]interface{}{"a": "1", "b": int64(2), "c": 3.14, "d": true}},
	{desc: "Array", value: []string{"a", "b", "c"}},
}

func TestRoundTrip(t *testing.T) {
	for _, testCase := range roundTripTestCases {
		t.Run(testCase.desc, func(t *testing.T) {
			c := qt.New(t)
			n, err := fluent.Reflect(basicnode.Prototype.Any, testCase.value)
			c.Assert(err, qt.IsNil)
			out, err := fluent.ToInterface(n)
			c.Assert(err, qt.IsNil)
			outJson, err := json.Marshal(out)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, outJson, qt.JSONEquals, testCase.value)
		})
	}
}

func TestLink(t *testing.T) {
	c := qt.New(t)
	var someCid, err = cid.Parse("bafybeihrqe2hmfauph5yfbd6ucv7njqpiy4tvbewlvhzjl4bhnyiu6h7pm")
	c.Assert(err, qt.IsNil)
	var link = cidlink.Link{Cid: someCid}
	var node = basicnode.NewLink(link)
	v, err := fluent.ToInterface(node)
	c.Assert(err, qt.IsNil)

	c.Assert(v.(cidlink.Link), qt.Equals, link)
}
