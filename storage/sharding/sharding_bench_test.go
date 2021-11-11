package sharding

import (
	"testing"
)

var globalSink interface{}

// This doesn't benchmark each of the sharding functions because... they're all roughly the same, really.
// It's mainly to make sure that our documentation's claim about zero-alloc operation is true.
func Benchmark(b *testing.B) {
	b.ReportAllocs()
	k := "abcdefgh"
	v := make([]string, 0, 3)
	var sink string
	for n := 0; n < b.N; n++ {
		v = v[0:0]
		Shard_r133(k, &v)
		sink = v[1]
	}
	globalSink = sink // make very very sure the compiler can't optimize our 'v' into oblivion.
}
