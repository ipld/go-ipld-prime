package ipldfree

import (
	"fmt"
	"runtime"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
)

func BenchmarkJustString(b *testing.B) {
	var node ipld.Node
	for i := 0; i < b.N; i++ {
		node = String("boxme")
	}
	_ = node
}

func BenchmarkJustStringUse(b *testing.B) {
	var node ipld.Node
	for i := 0; i < b.N; i++ {
		node = String("boxme")
		s, err := node.AsString()
		_ = s
		_ = err
	}
}

func BenchmarkJustStringLogAllocs(b *testing.B) {
	memUsage := func(m1, m2 *runtime.MemStats) {
		fmt.Println(
			"Alloc:", m2.Alloc-m1.Alloc,
			"TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc,
			"HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc,
		)
	}
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	var node ipld.Node = String("boxme")
	runtime.ReadMemStats(&m2)
	memUsage(&m1, &m2)
	sinkNode = node // necessary to avoid clever elision.
}

var sinkNode ipld.Node
