package bsrvadapter

import (
	"context"
	"fmt"

	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
)

// Adapter implements go-ipld-prime/storage.ReadableStorage
// and go-ipld-prime/storage.WritableStorage
// backed by a go-blockservice.BlockService.
//
// The go-blockservice.BlockService may internally have other configuration,
// and contain whole other systems like Bitswap for transport.
// We don't interfere with that here;
// such configuration should be handled when creating the go-blockservice value.
//
// Note that this system will only work for certain structures of keys --
// this is because the blockservice API works on the level of CIDs.
// As long as your key string is the binary form of a CID, it will work correctly.
// Other keys are not possible to support with this adapter.
//
// Contexts given to this system are passed through where possible, but it is not possible in all cases.
// For operations where the underlying interface doesn't accept a context parameter,
// this adapter will check the context for errors before beginning an operation,
// but the context will otherwise have no effect.
// For operations where BlockService does accept a context, we pass it on.
type Adapter struct {
	Wrapped blockservice.BlockService
}

// Has implements go-ipld-prime/storage.Storage.Has.
//
// Note that for a BlockService, the Has operation has rather unusual semantics.
// Has may return false, and an immediately subsequent Get for the same key might return data!
// This is because the Has operation is defined as whether
// the Blockstore that the BlockService wraps has the requested key, immediately, locally;
// while the Get operation might use the BlockService to go _find_ the requested key
// and its content, even remotely!
func (a *Adapter) Has(ctx context.Context, key string) (bool, error) {
	// Return early if the context is already closed.
	// This is also the last time we'll check the context,
	// since the Has method is on Blockstore, which doesn't take them.
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
	return a.Wrapped.Blockstore().Has(ctx, k)
}

// Get implements go-ipld-prime/storage.ReadableStorage.Get.
func (a *Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	// No need to check the context proactively here --
	// the BlockService API actually accepts context.

	// Do the inverse of cid.KeyString(),
	// which is how a valid key for this adapter must've been produced.
	k, err := cidFromBinString(key)
	if err != nil {
		return nil, err
	}

	// Delegate the Get call.
	// It's called "GetBlock" in BlockService.
	block, err := a.Wrapped.GetBlock(ctx, k)
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
	// since the AddBlock method on BlockService that we'll eventually be delegating to doesn't take them.
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
	// This is necessary because it's the format demanded by BlockService.
	// (Unfortunately, it also provokes an allocation, because it uses interfaces;
	// but we can't avoid that without changing the code in go-blockservice.)
	// The error is treated as a panic because it's only possible if a global debug var is set,
	// and is for behavior that is not meant to be part of the contract of the storage APIs.
	block, err := blocks.NewBlockWithCid(content, k)
	if err != nil {
		panic(err)
	}

	// Delegate the Put call.
	// It's called "AddBlock" in BlockService.
	return a.Wrapped.AddBlock(ctx, block)
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
