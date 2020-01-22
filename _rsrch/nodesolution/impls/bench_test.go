package impls

import (
	"encoding/json"
	"testing"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
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
		if err := json.Unmarshal([]byte(`{"whee":1,"woot":2,"waga":3}`), &x); err != nil {
			panic(err)
		}
		sink = x
	}
}

func BenchmarkMap3nFeedGenericMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var nb ipld.NodeBuilder
		nb = Style__Map{}.NewBuilder()
		ma, err := nb.BeginMap(3)
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
		if err := ma.AssembleKey().AssignString("waga"); err != nil {
			panic(err)
		}
		if err := ma.AssembleValue().AssignInt(3); err != nil {
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

func BenchmarkMap3nFeedGennedMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var nb ipld.NodeBuilder
		nb = Type__Map_K_T{}.NewBuilder()
		ma, err := nb.BeginMap(3)
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
		if err := ma.AssembleKey().AssignString("waga"); err != nil {
			panic(err)
		}
		if err := ma.AssembleValue().AssignInt(3); err != nil {
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
			if err := va.AssignInt(1); err != nil {
				panic(err)
			}
		}
		if va, err := ma.AssembleDirectly("woot"); err != nil {
			panic(err)
		} else {
			if err := va.AssignInt(2); err != nil {
				panic(err)
			}
		}
		if va, err := ma.AssembleDirectly("waga"); err != nil {
			panic(err)
		} else {
			if err := va.AssignInt(3); err != nil {
				panic(err)
			}
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
