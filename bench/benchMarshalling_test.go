package bench

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt/json"

	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/impl/free"
)

// This is identical to something in the refmt project benchmarks.
// We have a *radically* different approach than refmt's "obj" package does, so this will be... interesting.
// ... update: actually, they're... pretty close.
// Unmarshal is within 20% of the same performance, before we've done any work to optimize go-ipld-prime (but a bit slower).
// Marshal is actually taking only 70% as long as refmt did on a map -- our node accessors are *more* efficient than reflection.
// Okay then!  Not bad.
var fixture_structAlpha_json = []byte(`{"B":{"R":{"M":"quir","R":{"M":"asdf","R":{"M":"","R":null}}}},"C":{"M":13,"N":"n"},"C2":{"M":14,"N":"n2"},"W":"4","X":1,"Y":2,"Z":"3"}`)

var sink interface{}

func Benchmark_MapAlpha_UnmarshalIntoFreenode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nb := ipldfree.NodeBuilder()
		n, err := encoding.Unmarshal(nb, json.NewDecoder(bytes.NewReader(fixture_structAlpha_json)))
		if err != nil {
			panic(err)
		}
		sink = n
	}
}

func Benchmark_MapAlpha_MarshalFromFreenode(b *testing.B) {
	nb := ipldfree.NodeBuilder()
	n, err := encoding.Unmarshal(nb, json.NewDecoder(bytes.NewReader(fixture_structAlpha_json)))
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		encoding.Marshal(n, json.NewEncoder(&buf, json.EncodeOptions{}))
	}
}
