package multihoisting

import (
	"testing"
)

var sink interface{}

type nest struct {
	a, b, c nest2
}

type nest2 struct {
	d nest3
}

type nest3 string

type fizz interface {
	buzz(string) fizz
}

func (x nest) buzz(k string) fizz {
	switch k {
	case "a":
		return x.a
	case "b":
		return x.b
	case "c":
		return x.c
	default:
		return nil
	}
}
func (x nest2) buzz(k string) fizz {
	switch k {
	case "d":
		return x.d
	default:
		return nil
	}
}
func (x nest3) buzz(k string) fizz {
	return x
}

// b.Log(unsafe.Sizeof(nest{}))

// These all three score:
//		BenchmarkWot1-8         50000000                31.6 ns/op            16 B/op          1 allocs/op
//		BenchmarkWot2-8         30000000                35.3 ns/op            16 B/op          1 allocs/op
//		BenchmarkWot3-8         50000000                37.1 ns/op            16 B/op          1 allocs/op
//
// Don't entirely get it (namely the 16).
// Something fancy that's enabled by inlining, for sure.
// The one alloc makes enough sense: it's just what's going into 'sink'.
func BenchmarkWot1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := nest{}
		sink = v.buzz("a")
	}
}
func BenchmarkWot2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := nest{}
		fizzer := v.buzz("a")
		sink = fizzer.buzz("d")
	}
}
func BenchmarkWot3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := nest{}
		fizzer := v.buzz("a")
		fizzer = fizzer.buzz("d")
		sink = fizzer.buzz(".")
	}
}

// This comes out pretty much like you'd expect:
//		BenchmarkZot1-8         20000000                80.6 ns/op            64 B/op          2 allocs/op
//		BenchmarkZot2-8         20000000                88.1 ns/op            64 B/op          2 allocs/op
//		BenchmarkZot3-8         20000000                84.9 ns/op            64 B/op          2 allocs/op
//
// The `nest` type gets moved to the heap -- 1 alloc, 48 bytes.
// Then the `nest2` type also gets moved to the heap when returned -- 1 alloc, 16 bytes.
// Wait, where's 3?  one of either 2 or 3 is getting magiced, but which and why
// is it because the addr of 3 is literally 2, or
func BenchmarkZot1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := nest{}
		sink = buzzit(v, "a")
	}
}
func BenchmarkZot2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := nest{}
		fizzer := buzzit(v, "a")
		sink = buzzit(fizzer, "d")
	}
}
func BenchmarkZot3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := nest{}
		fizzer := buzzit(v, "a")
		fizzer = buzzit(fizzer, "d")
		sink = buzzit(fizzer, ".")
	}
}

// This is a function that bamboozles inlining.
// (Note you need to take a ptr to it for that to work.)
// (FIXME DOC wait what??? the above is not true, why)
func buzzit(fizzer fizz, k string) fizz {
	return fizzer.buzz(k)
}

type wider struct {
	z, y, x nest
}

func (x wider) buzz(k string) fizz {
	switch k {
	case "z":
		return x.z
	case "y":
		return x.y
	case "x":
		return x.x
	default:
		return nil
	}
}
func BenchmarkZot4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := wider{}
		fizzer := buzzit(v, "z")
		fizzer = buzzit(fizzer, "a")
		fizzer = buzzit(fizzer, "d")
		sink = buzzit(fizzer, ".")
	}
}

// So the big question is, can we get multiple heap pointers out of a single move,
// and is that something we can do with a choice in one place in advance
// (rather than requiring different code in a fractal of use sites to agree)?

func BenchmarkMultiStartingStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := wider{}
		fizzer1 := buzzit(v, "z")
		fizzer2 := buzzit(fizzer1, "a")
		fizzer2 = buzzit(fizzer2, "d")
		sink = buzzit(fizzer2, ".")
		fizzer2 = buzzit(fizzer1, "b")
		fizzer2 = buzzit(fizzer2, "d")
		sink = buzzit(fizzer2, ".")
	}
}
func BenchmarkMultiStartingHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := &wider{}
		fizzer1 := buzzit(v, "z")
		fizzer2 := buzzit(fizzer1, "a")
		fizzer2 = buzzit(fizzer2, "d")
		sink = buzzit(fizzer2, ".")
		fizzer2 = buzzit(fizzer1, "b")
		fizzer2 = buzzit(fizzer2, "d")
		sink = buzzit(fizzer2, ".")
	}
}

func escape(root *nest) fizz {
	confound := 8
	var fizzer fizz
	fizzer = root
	if confound%2 == 0 {
		fizzer = &root.a
	}
	if confound%4 == 0 {
		fizzer = &root.b
	}
	if confound%8 == 0 {
		fizzer = &root.c
	}
	return fizzer
}
