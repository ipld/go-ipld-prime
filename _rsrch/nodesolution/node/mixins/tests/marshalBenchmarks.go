package tests

import (
	"bytes"
	"testing"

	refmtjson "github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/tok"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/codec"
	"github.com/ipld/go-ipld-prime/must"
)

// All of the marshalling and unmarshalling benchmark specs use JSON.
// This does mean we're measuring a bunch of stuff that has nothing to do
//  with the core operations of the Node/NodeBuilder interface.
// We do this so that:
// - we get a reasonable picture of how much time is spent in the IPLD Data Model
//    versus how much time is spent in the serialization efforts;
// - we can make direct comparisons to the standard library json marshalling
//    and unmarshalling, thus having a back-of-the-envelope baseline to compare.

func SpecBenchmarkMarshalMapStrInt_3n(b *testing.B, ns ipld.NodeStyle) {
	nb := ns.NewBuilder()
	must.NotError(codec.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(`{"whee":1,"woot":2,"waga":3}`))))
	n := nb.Build()
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err = codec.Marshal(n, refmtjson.NewEncoder(&buf, refmtjson.EncodeOptions{}))
		sink = buf
	}
	if err != nil {
		panic(err)
	}
}

func SpecBenchmarkMarshalToNullMapStrInt_3n(b *testing.B, ns ipld.NodeStyle) {
	nb := ns.NewBuilder()
	must.NotError(codec.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(`{"whee":1,"woot":2,"waga":3}`))))
	n := nb.Build()
	b.ResetTimer()
	var err error
	encoder := &nullTokenSink{}
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err = codec.Marshal(n, encoder)
		sink = buf
	}
	if err != nil {
		panic(err)
	}
}

type nullTokenSink struct{}

func (nullTokenSink) Step(_ *tok.Token) (bool, error) {
	return false, nil
}
