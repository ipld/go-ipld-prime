package impls

import (
	"encoding/json"
	"testing"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

var sink interface{}

func BenchmarkBaselineNativeMapAssignSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x = make(map[string]int, 2)
		x["whee"] = 1
		x["woot"] = 2
		sink = x
	}
}

func BenchmarkBaselineJsonUnmarshalMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x = make(map[string]int, 2)
		if err := json.Unmarshal([]byte(`{"whee":1,"woot":2}`), &x); err != nil {
			panic(err)
		}
		sink = x
	}
}

func BenchmarkFeedGennedMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var nb ipld.NodeBuilder
		nb = NewBuilder_Map_K_T()
		ma, err := nb.BeginMap(2)
		if err != nil {
			panic(err)
		}
		if err := ma.AssembleKey().AssignString("whee"); err != nil {
			panic(err)
		}
		if err := ma.AssembleValue().AssignInt(1); err != nil {
			panic(err)
		}
		if err := ma.AssembleKey().AssignString("woot"); err != nil {
			panic(err)
		}
		if err := ma.AssembleValue().AssignInt(2); err != nil {
			panic(err)
		}
		if err := ma.Done(); err != nil {
			panic(err)
		}
		if n, err := nb.Build(); err != nil {
			panic(err)
		} else {
			sink = n
		}
	}
}
