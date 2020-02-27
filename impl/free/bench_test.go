package ipldfree

import (
	"strconv"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/tests"
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

var tableStrInt = [25]struct {
	s string
	i int
}{}

func init() {
	for i := 1; i <= 25; i++ {
		tableStrInt[i-1] = struct {
			s string
			i int
		}{"k" + strconv.Itoa(i), i}
	}
}

func buildMapStrIntN25() ipld.Node {
	nb := NodeBuilder()
	mb, err := nb.CreateMap()
	mustNotError(err)
	kb := mb.BuilderForKeys()
	vb := mb.BuilderForValue("")
	for i := 1; i <= 25; i++ {
		mustNotError(mb.Insert(mustNode(kb.CreateString(tableStrInt[i-1].s)), mustNode(vb.CreateInt(tableStrInt[i-1].i))))
	}
	return mustNode(mb.Build())
}

func BenchmarkMap25nFeedGenericMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN25()
	}
}

func BenchmarkMap25nGenericMapIterationSimpleKeys(b *testing.B) {
	n := buildMapStrIntN25()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}

// benchmarks covering encoding -->

func BenchmarkSpec_Marshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt(b, NodeBuilder())
}
func BenchmarkSpec_Marshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_MapNStrMap3StrInt(b, NodeBuilder())
}

func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_Map3StrInt(b, NodeBuilder())
}
func BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b, NodeBuilder())
}

// benchmarks covering traversal -->

func BenchmarkSpec_Walk_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_Map3StrInt(b, NodeBuilder())
}
func BenchmarkSpec_Walk_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_MapNStrMap3StrInt(b, NodeBuilder())
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
