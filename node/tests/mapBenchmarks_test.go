package tests

import (
	"encoding/json"
	"testing"

	"github.com/ipld/go-ipld-prime/must"
)

// This is analogous to the 'MapStrInt_3n' suite of benchmarks,
// but against a golang native map in regular go code,
// for getting a baseline impression to compare other things against.
func BenchmarkMapStrInt_3n_BaselineNativeMapAssignSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := make(map[string]int, 3)
		x["whee"] = 1
		x["woot"] = 2
		x["waga"] = 3
		sink = x
	}
}

func BenchmarkMapStrInt_3n_BaselineJsonUnmarshalMapSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := make(map[string]int, 3)
		must.NotError(json.Unmarshal([]byte(`{"whee":1,"woot":2,"waga":3}`), &x))
		sink = x
	}
}

func BenchmarkMapStrInt_3n_BaselineJsonMarshalMapSimpleKeys(b *testing.B) {
	x := map[string]int{"whee": 1, "woot": 2, "waga": 3}
	for i := 0; i < b.N; i++ {
		bs, err := json.Marshal(x)
		must.NotError(err)
		sink = bs
	}
}

var (
	sink_s string
	sink_i int64
)

func BenchmarkMapStrInt_3n_BaselineNativeMapIterationSimpleKeys(b *testing.B) {
	x := make(map[string]int64, 3)
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

// n25 -->

func BenchmarkMapStrInt_25n_BaselineNativeMapAssignSimpleKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := make(map[string]int64, 25)
		for i := 1; i <= 25; i++ {
			x[tableStrInt[i-1].s] = tableStrInt[i-1].i
		}
		sink = x
	}
}

func BenchmarkMapStrInt_25n_BaselineNativeMapIterationSimpleKeys(b *testing.B) {
	x := make(map[string]int64, 25)
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
