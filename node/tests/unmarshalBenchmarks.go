package tests

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/tests/corpus"
)

// All of the marshalling and unmarshalling benchmark specs use JSON.
// This does mean we're measuring a bunch of stuff that has nothing to do
//  with the core operations of the Node/NodeBuilder interface.
// We do this so that:
// - we get a reasonable picture of how much time is spent in the IPLD Data Model
//    versus how much time is spent in the serialization efforts;
// - we can make direct comparisons to the standard library json marshalling
//    and unmarshalling, thus having a back-of-the-envelope baseline to compare.

func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B, np datamodel.NodePrototype) {
	var err error
	for i := 0; i < b.N; i++ {
		nb := np.NewBuilder()
		err = json.Decode(nb, strings.NewReader(`{"whee":1,"woot":2,"waga":3}`))
		sink = nb.Build()
	}
	if err != nil {
		panic(err)
	}
}

func BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b *testing.B, np datamodel.NodePrototype) {
	for _, n := range []int{0, 1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			msg := corpus.MapNStrMap3StrInt(n)
			b.ResetTimer()

			var node datamodel.Node
			var err error
			nb := np.NewBuilder()
			for i := 0; i < b.N; i++ {
				err = json.Decode(nb, strings.NewReader(msg))
				node = nb.Build()
				nb.Reset()
			}

			b.StopTimer()
			if err != nil {
				b.Fatalf("decode errored: %s", err)
			}
			var buf bytes.Buffer
			json.Encode(node, &buf)
			if buf.String() != msg {
				b.Fatalf("re-encode result didn't match corpus")
			}
		})
	}
}
