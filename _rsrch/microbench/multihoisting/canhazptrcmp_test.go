package multihoisting

import (
	"testing"
)

func BenchmarkCmpPtrs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := hoisty{}
		e := returnE(&v.b)
		if e == &v.b.e {
			sink = e
		} else {
			b.Fail()
		}
		sink = e
	}
}

func BenchmarkCmpVal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := hoisty{}
		e := returnE(&v.b)
		if *e == v.b.e {
			sink = e
		} else {
			b.Fail()
		}
		sink = e
	}
}

// not sure these are nailing it because no interfaces were harmed in the making
// also the sink is... gonna cost different things for these; need something else for the inside of the 'if'.
//  ^ jk that last, actually.  it does one alloc either way... just shuffles which line technically triggers it.
// the value compare *is* slower, but on the order of two nanoseconds.
