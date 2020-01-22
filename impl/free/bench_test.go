package ipldfree

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
)

var sink interface{}

func buildMapStrIntN3() ipld.Node {
	nb := NodeBuilder()
	mb, err := nb.CreateMap()
	mustNotError(err)
	kb := mb.BuilderForKeys()
	vb := mb.BuilderForValue("")
	mustNotError(mb.Insert(mustNode(kb.CreateString("whee")), mustNode(vb.CreateInt(1))))
	mustNotError(mb.Insert(mustNode(kb.CreateString("woot")), mustNode(vb.CreateInt(2))))
	mustNotError(mb.Insert(mustNode(kb.CreateString("waga")), mustNode(vb.CreateInt(3))))
	return mustNode(mb.Build())
}

func BenchmarkMap3nFeedGenericMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN3()
	}
}

func BenchmarkMap3nGenericMapIterationSimpleKeys(b *testing.B) {
	n := buildMapStrIntN3()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}

// copy of helper functions from must package, because import cycles, sigh -->

func mustNotError(e error) {
	if e != nil {
		panic(e)
	}
}
func mustNode(n ipld.Node, e error) ipld.Node {
	if e != nil {
		panic(e)
	}
	return n
}
