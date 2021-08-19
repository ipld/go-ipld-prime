package tests

import (
	"fmt"
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/tests/corpus"
	"github.com/ipld/go-ipld-prime/traversal"
)

func BenchmarkSpec_Walk_Map3StrInt(b *testing.B, np datamodel.NodePrototype) {
	node := mustNodeFromJsonString(np, corpus.Map3StrInt())
	sel := mustSelectorFromJsonString(np, `{"a":{">":{".":{}}}}`)
	b.ResetTimer()

	var visitCountSanityCheck int
	for i := 0; i < b.N; i++ {
		visitCountSanityCheck = 0
		traversal.WalkMatching(node, sel, func(tp traversal.Progress, n datamodel.Node) error {
			visitCountSanityCheck++ // this sanity check is sufficiently cheap to be worth it
			return nil              // no need to do anything here; just care about exercising the walk internals.
		})
	}
	if visitCountSanityCheck != 3 {
		b.Fatalf("visitCountSanityCheck should be 3, got %d", visitCountSanityCheck)
	}
}

func BenchmarkSpec_Walk_MapNStrMap3StrInt(b *testing.B, np datamodel.NodePrototype) {
	sel := mustSelectorFromJsonString(np, `{"a":{">":{"a":{">":{".":{}}}}}}`)

	for _, n := range []int{0, 1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			node := mustNodeFromJsonString(np, corpus.MapNStrMap3StrInt(n))
			b.ResetTimer()

			var visitCountSanityCheck int
			for i := 0; i < b.N; i++ {
				visitCountSanityCheck = 0
				traversal.WalkMatching(node, sel, func(tp traversal.Progress, n datamodel.Node) error {
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
