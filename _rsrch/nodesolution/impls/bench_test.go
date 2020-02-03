package impls

import (
	"encoding/json"
	"testing"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
	"github.com/ipld/go-ipld-prime/must"
)

var sink interface{}

func BenchmarkMap3nBaselineNativeMapAssignSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x = make(map[string]int, 3)
		x["whee"] = 1
		x["woot"] = 2
		x["waga"] = 3
		sink = x
	}
}

func BenchmarkMap3nBaselineJsonUnmarshalMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x = make(map[string]int, 3)
		must.NotError(json.Unmarshal([]byte(`{"whee":1,"woot":2,"waga":3}`), &x))
		sink = x
	}
}

func BenchmarkMap3nFeedGenericMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN3(Style__Map{})
	}
}

func BenchmarkMap3nFeedGennedMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN3(Type__Map_K_T{})
	}
}

func BenchmarkMap3nFeedGennedMapSimpleKeysDirectly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var nb ipld.NodeBuilder
		nb = Type__Map_K_T{}.NewBuilder()
		ma, err := nb.BeginMap(3)
		if err != nil {
			panic(err)
		}
		if va, err := ma.AssembleDirectly("whee"); err != nil {
			panic(err)
		} else {
			must.NotError(va.AssignInt(1))
		}
		if va, err := ma.AssembleDirectly("woot"); err != nil {
			panic(err)
		} else {
			must.NotError(va.AssignInt(2))
		}
		if va, err := ma.AssembleDirectly("waga"); err != nil {
			panic(err)
		} else {
			must.NotError(va.AssignInt(3))
		}
		must.NotError(ma.Finish())
		sink = nb.Build()
	}
}

var sink_s string
var sink_i int

func BenchmarkMap3nBaselineNativeMapIterationSimpleKeys(b *testing.B) {
	var x = make(map[string]int, 3)
	x["whee"] = 1
	x["woot"] = 2
	x["waga"] = 3
	sink = x
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k, v := range x {
			sink_s = k
			sink_i = v
		}
	}
}

func BenchmarkMap3nGenericMapIterationSimpleKeys(b *testing.B) {
	n := buildMapStrIntN3(Style__Map{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}

func BenchmarkMap3nGennedMapIterationSimpleKeys(b *testing.B) {
	n := buildMapStrIntN3(Type__Map_K_T{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}

// n25 -->

func BenchmarkMap25nBaselineNativeMapAssignSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x = make(map[string]int, 25)
		for i := 1; i <= 25; i++ {
			x[tableStrInt[i-1].s] = tableStrInt[i-1].i
		}
		sink = x
	}
}

func BenchmarkMap25nFeedGenericMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN25(Style__Map{})
	}
}

func BenchmarkMap25nFeedGennedMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = buildMapStrIntN25(Type__Map_K_T{})
	}
}

func BenchmarkMap25nFeedGennedMapSimpleKeysDirectly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var nb ipld.NodeBuilder
		nb = Type__Map_K_T{}.NewBuilder()
		ma, err := nb.BeginMap(25)
		if err != nil {
			panic(err)
		}
		for i := 1; i <= 25; i++ {
			if va, err := ma.AssembleDirectly(tableStrInt[i-1].s); err != nil {
				panic(err)
			} else {
				must.NotError(va.AssignInt(tableStrInt[i-1].i))
			}
		}
		must.NotError(ma.Finish())
		sink = nb.Build()
	}
}

func BenchmarkMap25nBaselineNativeMapIterationSimpleKeys(b *testing.B) {
	var x = make(map[string]int, 25)
	for i := 1; i <= 25; i++ {
		x[tableStrInt[i-1].s] = tableStrInt[i-1].i
	}
	sink = x
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k, v := range x {
			sink_s = k
			sink_i = v
		}
	}
}

func BenchmarkMap25nGenericMapIterationSimpleKeys(b *testing.B) {
	n := buildMapStrIntN25(Style__Map{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}

func BenchmarkMap25nGennedMapIterationSimpleKeys(b *testing.B) {
	n := buildMapStrIntN25(Type__Map_K_T{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itr := n.MapIterator()
		for k, v, _ := itr.Next(); !itr.Done(); k, v, _ = itr.Next() {
			sink = k
			sink = v
		}
	}
}
