package tests

import (
	"bytes"
	"testing"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/tests/corpus"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func SpecBenchmarkWalkMapStrInt_3n(b *testing.B, nb ipld.NodeBuilder) {
	n, err := encoding.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(corpus.Map3StrInt())))
	if err != nil {
		panic(err)
	}
	seldefn, err := encoding.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(`{"a":{">":{".":{}}}}`)))
	if err != nil {
		panic(err)
	}
	sel, err := selector.ParseSelector(seldefn)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		traversal.WalkMatching(n, sel, func(tp traversal.Progress, n ipld.Node) error {
			return nil // no need to do anything here; just care about exercising the walk internals.
		})
	}
}
