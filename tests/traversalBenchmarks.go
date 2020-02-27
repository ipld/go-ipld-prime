package tests

import (
	"bytes"
	"testing"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/tests/corpus"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func SpecBenchmarkWalkMapStrInt_3n(b *testing.B, nb ipld.NodeBuilder) {
	n := must.Node(encoding.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(corpus.Map3StrInt()))))
	seldefn := must.Node(encoding.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(`{"a":{">":{".":{}}}}`))))
	sel, err := selector.ParseSelector(seldefn)
	must.NotError(err)
	b.ResetTimer()

	var visitCountSanityCheck int
	for i := 0; i < b.N; i++ {
		visitCountSanityCheck = 0
		traversal.WalkMatching(n, sel, func(tp traversal.Progress, n ipld.Node) error {
			visitCountSanityCheck++ // this sanity check is sufficiently cheap to be worth it
			return nil              // no need to do anything here; just care about exercising the walk internals.
		})
	}
	if visitCountSanityCheck != 3 {
		b.Fatalf("visitCountSanityCheck should be 3, got %d", visitCountSanityCheck)
	}
}
