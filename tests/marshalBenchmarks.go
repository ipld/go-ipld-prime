package tests

import (
	"bytes"
	"fmt"
	"testing"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/tests/corpus"
)

// All of the marshalling and unmarshalling benchmark specs use JSON.
// This does mean we're measuring a bunch of stuff that has nothing to do
//  with the core operations of the Node/NodeBuilder interface.
// We do this so that:
// - we get a reasonable picture of how much time is spent in the IPLD Data Model
//    versus how much time is spent in the serialization efforts;
// - we can make direct comparisons to the standard library json marshalling
//    and unmarshalling, thus having a back-of-the-envelope baseline to compare.

func BenchmarkSpec_Marshal_Map3StrInt(b *testing.B, nb ipld.NodeBuilder) {
	msg := corpus.Map3StrInt()
	node := mustNodeFromJsonString(nb, msg)
	b.ResetTimer()

	var buf bytes.Buffer
	var err error
	for i := 0; i < b.N; i++ {
		buf = bytes.Buffer{}
		err = encoding.Marshal(node, refmtjson.NewEncoder(&buf, refmtjson.EncodeOptions{}))
	}

	b.StopTimer()
	if err != nil {
		b.Fatalf("marshal errored: %s", err)
	}
	if buf.String() != msg {
		b.Fatalf("marshal didn't match corpus")
	}
}

func BenchmarkSpec_Marshal_MapNStrMap3StrInt(b *testing.B, nb ipld.NodeBuilder) {
	for _, n := range []int{0, 1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			msg := corpus.MapNStrMap3StrInt(n)
			node := mustNodeFromJsonString(nb, msg)
			b.ResetTimer()

			var buf bytes.Buffer
			var err error
			for i := 0; i < b.N; i++ {
				buf = bytes.Buffer{}
				err = encoding.Marshal(node, refmtjson.NewEncoder(&buf, refmtjson.EncodeOptions{}))
			}

			b.StopTimer()
			if err != nil {
				b.Fatalf("marshal errored: %s", err)
			}
			if buf.String() != msg {
				b.Fatalf("marshal didn't match corpus")
			}
		})
	}
}
