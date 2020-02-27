package tests

import (
	"bytes"
	"fmt"
	"testing"

	refmtjson "github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/tests/corpus"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func mustNodeFromJsonString(nb ipld.NodeBuilder, str string) ipld.Node {
	return must.Node(encoding.Unmarshal(nb, refmtjson.NewDecoder(bytes.NewBufferString(str))))
}
func mustSelectorFromJsonString(nb ipld.NodeBuilder, str string) selector.Selector {
	// Needing an 'nb' parameter here is sort of off-topic, to be honest.
	//  Someday the selector package will probably contain codegen'd nodes of its own schema, and we'll use those unconditionally.
	//  For now... we'll just use whatever node you're already testing, because it oughta work
	//   (and because it avoids hardcoding any other implementation that might cause import cycle funtimes.).
	seldefn := mustNodeFromJsonString(nb, str)
	sel, err := selector.ParseSelector(seldefn)
	must.NotError(err)
	return sel
}

func BenchmarkSpec_Walk_Map3StrInt(b *testing.B, nb ipld.NodeBuilder) {
	node := mustNodeFromJsonString(nb, corpus.Map3StrInt())
	sel := mustSelectorFromJsonString(nb, `{"a":{">":{".":{}}}}`)
	b.ResetTimer()

	var visitCountSanityCheck int
	for i := 0; i < b.N; i++ {
		visitCountSanityCheck = 0
		traversal.WalkMatching(node, sel, func(tp traversal.Progress, n ipld.Node) error {
			visitCountSanityCheck++ // this sanity check is sufficiently cheap to be worth it
			return nil              // no need to do anything here; just care about exercising the walk internals.
		})
	}
	if visitCountSanityCheck != 3 {
		b.Fatalf("visitCountSanityCheck should be 3, got %d", visitCountSanityCheck)
	}
}

func BenchmarkSpec_Walk_MapNStrMap3StrInt(b *testing.B, nb ipld.NodeBuilder) {
	sel := mustSelectorFromJsonString(nb, `{"a":{">":{"a":{">":{".":{}}}}}}`)

	for _, n := range []int{0, 1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			node := mustNodeFromJsonString(nb, corpus.MapNStrMap3StrInt(n))
			b.ResetTimer()

			var visitCountSanityCheck int
			for i := 0; i < b.N; i++ {
				visitCountSanityCheck = 0
				traversal.WalkMatching(node, sel, func(tp traversal.Progress, n ipld.Node) error {
					visitCountSanityCheck++ // this sanity check is sufficiently cheap to be worth it
					return nil              // no need to do anything here; just care about exercising the walk internals.
				})
			}
			if visitCountSanityCheck != 3*n {
				b.Fatalf("visitCountSanityCheck should be %d, got %d", n*3, visitCountSanityCheck)
			}
		})
	}
}
