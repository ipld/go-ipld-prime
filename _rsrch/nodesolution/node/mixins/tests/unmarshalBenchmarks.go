package tests

import (
	"bytes"
	"testing"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/codec"
)

// All of the marshalling and unmarshalling benchmark specs use JSON.
// This does mean we're measuring a bunch of stuff that has nothing to do
//  with the core operations of the Node/NodeBuilder interface.
// We do this so that:
// - we get a reasonable picture of how much time is spent in the IPLD Data Model
//    versus how much time is spent in the serialization efforts;
// - we can make direct comparisons to the standard library json marshalling
//    and unmarshalling, thus having a back-of-the-envelope baseline to compare.

func SpecBenchmarkUnmarshalMapStrInt_3n(b *testing.B, ns ipld.NodeStyle) {
	var err error
	for i := 0; i < b.N; i++ {
		nb := ns.NewBuilder()
		err = codec.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(`{"whee":1,"woot":2,"waga":3}`)))
		sink = nb.Build()
	}
	if err != nil {
		panic(err)
	}
}
