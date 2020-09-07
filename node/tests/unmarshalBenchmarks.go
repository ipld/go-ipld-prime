package tests

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
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

func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B, np ipld.NodePrototype) {
	var err error
	for i := 0; i < b.N; i++ {
		nb := np.NewBuilder()
		err = codec.Unmarshal(nb, refmtjson.NewDecoder(strings.NewReader(`{"whee":1,"woot":2,"waga":3}`)))
		sink = nb.Build()
	}
	if err != nil {
		panic(err)
	}
}

func BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b *testing.B, np ipld.NodePrototype) {
	for _, n := range []int{0, 1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			msg := corpus.MapNStrMap3StrInt(n)
			b.ResetTimer()

			var node ipld.Node
			var err error
			nb := np.NewBuilder()
			for i := 0; i < b.N; i++ {
				err = codec.Unmarshal(nb, refmtjson.NewDecoder(strings.NewReader(msg)))
				node = nb.Build()
				nb.Reset()
			}

			b.StopTimer()
			if err != nil {
				b.Fatalf("unmarshal errored: %s", err)
			}
			var buf bytes.Buffer
			codec.Marshal(node, refmtjson.NewEncoder(&buf, refmtjson.EncodeOptions{}))
			if buf.String() != msg {
				b.Fatalf("remarshal didn't match corpus")
			}
		})
	}
}
