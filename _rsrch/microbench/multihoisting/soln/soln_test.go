package solution

import (
	"fmt"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(1) // necessary if we want to do precise accounting on runtime.ReadMemStats.
}

var sink interface{}

func TestAllocCount(t *testing.T) {
	memUsage := func(m1, m2 *runtime.MemStats) {
		fmt.Println(
			"Alloc:", m2.Alloc-m1.Alloc,
			"TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc,
			"HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc,
			"Mallocs:", m2.Mallocs-m1.Mallocs,
			"Frees:", m2.Frees-m1.Frees,
		)
	}
	var m [99]runtime.MemStats
	runtime.GC()
	runtime.GC() // i know not why, but as of go-1.13.3, and not in go-1.12.5, i have to call this twice before we start to get consistent numbers.
	runtime.ReadMemStats(&m[0])

	var x Node
	x = &Stroct{}
	runtime.GC()
	runtime.ReadMemStats(&m[1])

	x = x.LookupByString("foo")
	runtime.GC()
	runtime.ReadMemStats(&m[2])

	sink = x
	runtime.GC()
	runtime.ReadMemStats(&m[3])

	sink = nil
	runtime.GC()
	runtime.ReadMemStats(&m[4])
	memUsage(&m[0], &m[1])
	memUsage(&m[0], &m[2])
	memUsage(&m[0], &m[3])
	memUsage(&m[0], &m[4])
}
