package dsadapter

import (
	"context"

	"github.com/ipfs/go-datastore"
)

// Adapter implements go-ipld-prime/storage.ReadableStorage
// and go-ipld-prime/storage.WritableStorage
// backed by a go-datastore.Datastore.
//
// Optionally, an EscapingFunc may also be set,
// which transforms the (possibly binary) keys considered acceptable
// by the go-ipld-prime/storage APIs into a subset that
// the go-datastore can accept.
// (Be careful to use any escaping consistently,
// and be wary of potential unexpected behavior if the escaping function might
// collapse two distinct keys into the same "escaped" key.)
//
// The go-datastore.Datastore may internally have other configuration,
// such as key sharding functions, etc, and we don't interfere with that here;
// such configuration should be handled when creating the go-datastore value.
//
// Contexts given to this system are checked for errors at the beginning of an operation,
// but otherwise have no effect, because the Datastore API doesn't accept context parameters.
type Adapter struct {
	Wrapped      datastore.Datastore
	EscapingFunc func(string) string
}

// Has implements go-ipld-prime/storage.Storage.Has.
func (a *Adapter) Has(ctx context.Context, key string) (bool, error) {
	// Return early if the context is already closed.
	// This is also the last time we'll check the context,
	// since go-datastore doesn't take them.
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// If we have an EscapingFunc, apply it.
	if a.EscapingFunc != nil {
		key = a.EscapingFunc(key)
	}

	// Wrap the key into go-datastore's concrete type that it requires.
	// Note that this does a bunch of actual work, which may be surprising.
	// The key may be transformed (as per path.Clean).
	// There will also be an allocation, if the key doesn't start with "/".
	// (Avoiding these performance drags is part of why we started
	// new interfaces in go-ipld-prime/storage.)
	k := datastore.NewKey(key)

	// Delegate the has call.
	// Note that for some datastore implementations, this will do *yet more*
	// validation on the key, and may return errors from that.
	return a.Wrapped.Has(ctx, k)
}

// Get implements go-ipld-prime/storage.ReadableStorage.Get.
func (a *Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	// Return early if the context is already closed.
	// This is also the last time we'll check the context,
	// since go-datastore doesn't take them.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// If we have an EscapingFunc, apply it.
	if a.EscapingFunc != nil {
		key = a.EscapingFunc(key)
	}

	// Wrap the key into go-datastore's concrete type that it requires.
	// Note that this does a bunch of actual work, which may be surprising.
	// The key may be transformed (as per path.Clean).
	// There will also be an allocation, if the key doesn't start with "/".
	// (Avoiding these performance drags is part of why we started
	// new interfaces in go-ipld-prime/storage.)
	k := datastore.NewKey(key)

	// Delegate the get call.
	// Note that for some datastore implementations, this will do *yet more*
	// validation on the key, and may return errors from that.
	return a.Wrapped.Get(ctx, k)
}

// Put implements go-ipld-prime/storage.WritableStorage.Put.
func (a *Adapter) Put(ctx context.Context, key string, content []byte) error {
	// Return early if the context is already closed.
	// This is also the last time we'll check the context,
	// since go-datastore doesn't take them.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// If we have an EscapingFunc, apply it.
	if a.EscapingFunc != nil {
		key = a.EscapingFunc(key)
	}

	// Wrap the key into go-datastore's concrete type that it requires.
	// Note that this does a bunch of actual work, which may be surprising.
	// The key may be transformed (as per path.Clean).
	// There will also be an allocation, if the key doesn't start with "/".
	// (Avoiding these performance drags is part of why we started
	// new interfaces in go-ipld-prime/storage.)
	k := datastore.NewKey(key)

	// Delegate the put call.
	// Note that for some datastore implementations, this will do *yet more*
	// validation on the key, and may return errors from that.
	return a.Wrapped.Put(ctx, k, content)
}
