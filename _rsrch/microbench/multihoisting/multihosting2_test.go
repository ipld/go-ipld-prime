package multihoisting

import (
	"testing"
)

type hoisty struct {
	a, b, c hoisty2
}
type hoisty2 struct {
	d, e hoisty3
}
type hoisty3 struct {
	f hoisty4
}
type hoisty4 string

func BenchmarkHoist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := hoisty{}
		sink = returnE(&v.a)
		sink = returnE(&v.b)
		sink = returnE(&v.c)
	}
}

func BenchmarkHoist2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := &hoisty{}
		sink = returnE(&v.a)
		sink = returnE(&v.b)
		sink = returnE(&v.c)
	}
}

func returnE(x *hoisty2) *hoisty3 {
	return &x.e
}

// okay, am now confident the above two are the same.  escape analysis rocks em.

func BenchmarkHoistBranched(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := &hoisty{}
		sink = returnE(&v.a)
		oi := returnE(&v.b)
		sink = returnE(&v.c)
		sink = returnF(oi)
	}
}

func returnF(x *hoisty3) *hoisty4 {
	return &x.f
}

// now let's see if interfaces make this worse, somehow.
// (i'm *hoping* all the buzz methods being pointery will fail escape clearly,
//  and then returning more pointers from mid already-deffo-heap structures will be cheap.)

func BenchmarkHoistBranchedInterfacey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := &hoisty{}
		sink = v.buzz("a")
		oi := v.buzz("b")
		sink = v.buzz("c")
		sink = oi.buzz("f")
	}
}

func (x *hoisty) buzz(k string) fizz {
	switch k {
	case "a":
		return &x.a
	case "b":
		return &x.b
	case "c":
		return &x.c
	default:
		return nil
	}
}
func (x *hoisty2) buzz(k string) fizz {
	switch k {
	case "d":
		return &x.d
	case "e":
		return &x.e
	default:
		return nil
	}
}
func (x *hoisty3) buzz(k string) fizz {
	switch k {
	case "f":
		return &x.f
	default:
		return nil
	}
}
func (x *hoisty4) buzz(k string) fizz {
	return x
}

// also important: can i assign into this?
// ja.  it's fine.  no addntl allocs, with or without inlining.

func BenchmarkHoistAssign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := &hoisty{}
		oi := returnE(&v.b)
		oi.f = "yoi"
		sink = returnF(oi)
	}
}
