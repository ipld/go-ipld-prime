package tests

import (
	"context"
	"testing"

	"github.com/ipld/go-ipld-prime/storage"
)

/*
	General note:
	It's important to be careful to benchmark the cost *per op* --
	and not mistake b.N for the scale of corpus to work on.
	The corpus size should be a parameter that you supply,
	and your benchmark table should have a column for them!
*/

func BenchPut(b *testing.B, store storage.WritableStorage, gen Gen, scale int) {
	b.ReportAllocs()
	b.Logf("benchmarking with b.N=%d", b.N)

	// Use a fixed context throughout; it's not really relevant.
	ctx := context.Background()

	// Setup phase: create data up to the scale provided.
	// Reset the timer afterwards.
	b.Logf("prepopulating %d entries into storage...", scale)
	for n := 0; n < scale; n++ {
		key, content := gen()
		err := store.Put(ctx, key, content)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.Logf("prepopulating %d entries into storage: done.", scale)
	b.ResetTimer()

	// Now continue doing puts in the benchmark loop.
	// Note that if 'scale' was initially small, and b.N is big, results may be skewed,
	//  because the last put of the series will actually be working at scale+b.N-1.
	for n := 0; n < b.N; n++ {
		// Attempt to avoid counting any time spent by the gen func.
		//  ... except don't, because the overhead of starting and stopping is actually really high compared to a likely gen function;
		//   in practice, starting and stopping this frequently causes:
		//    - alloc count to be reported *correctly* (which is nice)
		//    - but reported ns/op to become erratic, and inflated (not at all nice)
		//    - and actuall wall-clock run time to increase drastically (~22x!) (not deadly, but certainly unpleasant)
		//   It may be best to write a synthetic dummy benchmark to see how much the gen function costs, and subtract that from the other results.
		//b.StopTimer()
		key, content := gen()
		//b.StartTimer()

		// Do the put.
		err := store.Put(ctx, key, content)
		if err != nil {
			b.Fatal(err)
		}
	}
}
