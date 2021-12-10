package bsadapter

import (
	"context"
	"fmt"

	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-blockstore"
)

// Adapter implements go-ipld-prime/storage.ReadableStorage
// and go-ipld-prime/storage.WritableStorage
// backed by a go-ipfs-blockstore.Blockstore.
//
// The go-ipfs-blockstore.Blockstore may internally have other configuration.
// We don't interfere with that here;
// such configuration should be handled when creating the go-ipfs-blockstore value.
//
// Note that this system will only work for certain structures of keys --
// this is because the blockstore API works on the level of CIDs.
// As long as your key string is the binary form of a CID, it will work correctly.
// Other keys are not possible to support with this adapter.
//
// Contexts given to this system are checked for errors at the beginning of an operation,
// but otherwise have no effect, because the Blockstore API doesn't accept context parameters.
type Adapter struct {
	Wrapped blockstore.Blockstore
}

// Has implements go-ipld-prime/storage.Storage.Has.
func (a *Adapter) Has(ctx context.Context, key string) (bool, error) {
	// Return early if the context is already closed.
	// This is also the last time we'll check the context,
	// since the Has method on Blockstore doesn't take them.
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// Do the inverse of cid.KeyString(),
	// which is how a valid key for this adapter must've been produced.
	k, err := cidFromBinString(key)
	if err != nil {
		return false, err
	}

	// Delegate the Has call.
	return a.Wrapped.Has(ctx, k)
}

// Get implements go-ipld-prime/storage.ReadableStorage.Get.
func (a *Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	// Return early if the context is already closed.
	// This is also the last time we'll check the context,
	// since the Put method on Blockstore doesn't take them.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Do the inverse of cid.KeyString(),
	// which is how a valid key for this adapter must've been produced.
	k, err := cidFromBinString(key)
	if err != nil {
		return nil, err
	}

	// Delegate the Get call.
	block, err := a.Wrapped.Get(ctx, k)
	if err != nil {
		return nil, err
	}

	// Unwrap the actual raw data for return.
	// Discard the rest.  (It's a shame there was an alloc for that structure.)
	return block.RawData(), nil
}

// Put implements go-ipld-prime/storage.WritableStorage.Put.
func (a *Adapter) Put(ctx context.Context, key string, content []byte) error {
	// Return early if the context is already closed.
	// This is also the last time we'll check the context,
	// since the Put method on Blockstore doesn't take them.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Do the inverse of cid.KeyString(),
	// which is how a valid key for this adapter must've been produced.
	k, err := cidFromBinString(key)
	if err != nil {
		return err
	}

	// Create a structure that has the cid and the raw content together.
	// This is necessary because it's the format demanded by Blockstore.
	// (Unfortunately, it also provokes an allocation, because it uses interfaces;
	// but we can't avoid that without changing the code in go-ipfs-blockstore.)
	// The error is treated as a panic because it's only possible if a global debug var is set,
	// and is for behavior that is not meant to be part of the contract of the storage APIs.
	block, err := blocks.NewBlockWithCid(content, k)
	if err != nil {
		panic(err)
	}

	// Delegate the Put call.
	return a.Wrapped.Put(ctx, block)
}

// Do the inverse of cid.KeyString().
// (Unclear why go-cid doesn't offer a function for this itself.)
func cidFromBinString(key string) (cid.Cid, error) {
	l, k, err := cid.CidFromBytes([]byte(key))
	if err != nil {
		return cid.Undef, fmt.Errorf("bsrvadapter: key was not a cid: %w", err)
	}
	if l != len(key) {
		return cid.Undef, fmt.Errorf("bsrvadapter: key was not a cid: had %d bytes leftover", len(key)-l)
	}
	return k, nil
}
