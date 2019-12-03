package multihoisting

import (
	"fmt"
	"runtime"
	"testing"
)

// not sure how to test this
// benchmem is going to report allocation costs... which we know.
// the question is if the memory usage goes *down* after full gc.

// okay here we go: `runtime.ReadMemStats` has consistency forcers.

func init() {
	runtime.GOMAXPROCS(1)
}

func BenchmarkReextentingGC(b *testing.B) {
	memUsage := func(m1, m2 *runtime.MemStats) {
		fmt.Println(
			"Alloc:", m2.Alloc-m1.Alloc,
			"TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc,
			"HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc,
			"Mallocs:", m2.Mallocs-m1.Mallocs,
			"Frees:", m2.Frees-m1.Frees,
		)
	}
	var m1, m2, m3, m4, m5 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	sink1 = &hoisty{}
	runtime.GC()
	runtime.ReadMemStats(&m2)
	sink2 = &sink1.b
	runtime.GC()
	runtime.ReadMemStats(&m3)
	sink1 = nil
	runtime.GC()
	runtime.ReadMemStats(&m4)
	sink2 = nil
	runtime.GC()
	runtime.ReadMemStats(&m5)
	fmt.Println("first extent size, size to get ref inside, size after dropping enclosing, size after dropping both")
	memUsage(&m1, &m2)
	memUsage(&m1, &m3)
	memUsage(&m1, &m4)
	memUsage(&m1, &m5)
}

var sink1 *hoisty
var sink2 *hoisty2

// results:
// yeah, one giant alloc occurs in the first move.
// subsequent pointer-getting causes no new memory usage.
// nilling the top level pointer does *not* let *anything* be collected.
// nilling both allows collection of the whole thing.

// so in practice:
// we can use internal pointers heavily without consequence in alloc count...
// but it's desirable not to combine items with different lifetimes into
//  a single extent in memory, because the longest living thing will extend
//   the life of everything it was allocated with.
