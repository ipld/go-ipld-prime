package benchmarks

import (
	"encoding/base32"
	"fmt"
	"os"
	"testing"

	flatfs "github.com/ipfs/go-ds-flatfs"
	"github.com/ipld/go-ipld-prime/storage"
	"github.com/ipld/go-ipld-prime/storage/dsadapter"
	"github.com/ipld/go-ipld-prime/storage/fsstore"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/ipld/go-ipld-prime/storage/tests"
)

func BenchmarkPut(b *testing.B) {
	// - memstore
	// - dsadapter with flatfs
	// - bsrvadapter wrapped around that (todo)
	// - fsstore
	tt := []struct {
		storeName        string // used in test name
		storeConstructor func() storage.WritableStorage
	}{{
		storeName: "memstore",
		storeConstructor: func() storage.WritableStorage {
			return &memstore.Store{}
		},
	}, {
		storeName: "dsadapter-flatfs-base32-defaultshard",
		storeConstructor: func() storage.WritableStorage {
			shardFn, err := flatfs.ParseShardFunc("/repo/flatfs/shard/v1/next-to-last/2")
			if err != nil {
				panic(err)
			}
			ds, err := flatfs.CreateOrOpen(".", shardFn, false)
			if err != nil {
				panic(err)
			}
			return &dsadapter.Adapter{
				Wrapped: ds,
				EscapingFunc: func(raw string) string {
					return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(raw))
				},
			}
		},
	}, {
		storeName: "fsstore-base32-defaultshard",
		storeConstructor: func() storage.WritableStorage {
			store := &fsstore.Store{}
			if err := store.InitDefaults("."); err != nil {
				panic(err)
			}
			return store
		},
	}}
	for _, ttr := range tt {
		for _, scale := range []int{
			// 1 << 8, // probably too small to be useful; b.N will always be much bigger than this.
			1 << 12,
			1 << 16,
			// 1 << 20, // already getting too big to fit the setup phase into the default benchmark time windows when using disk storage.
			// 1 << 24,
		} {
			b.Run(fmt.Sprintf("%s/scale=%d", ttr.storeName, scale), func(b *testing.B) {
				// Make a tempdir.  Change cwd to it.
				//  We'll assume the storage system, if it needs filesystem, can use the cwd.
				//  Using b.TempDir means the cleanup happens handled by the test system (and critically, not on our clock).
				dir := b.TempDir()
				retreat, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				defer os.Chdir(retreat)
				if err := os.Chdir(dir); err != nil {
					panic(err)
				}

				// Create the store, and put it to work!
				store := ttr.storeConstructor()
				gen := tests.NewCounterGen(1000000) // Use a large enough number that any sharding function kicks in (e.g. the b10 string is >=7 chars).
				tests.BenchPut(b, store, gen, scale)
			})
		}
	}
}
