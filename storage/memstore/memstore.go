package memstore

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

// Store is a simple in-memory storage.
// (It's little more than a map -- in fact, the map is exported,
// and you can poke it directly.)
//
// Store conforms to the storage.ReadableStorage and storage.WritableStorage APIs.
// Additionally, it supports storage.PeekableStorage and storage.StreamingReadableStorage,
// because it can do so while provoking fewer copies.
//
// If you want to use this store with streaming APIs,
// you can still do so by using the functions in the storage package,
// such as storage.GetStream and storage.PutStream, which will synthesize the correct behavior.
//
// You can use this storage with a linking.LinkSystem easily,
// by using the LinkSystem.SetReadStorage and/or LinkSystem.SetWriteStorage methods.
//
// There are no construction parameters for sharding functions nor escaping functions.
// Any keys are acceptable.
//
// This storage is mostly expected to be used for testing and demos,
// and as an example of how you can implement and integrate your own storage systems.
// It does not provide persistence beyond memory.
type Store struct {
	Bag map[string][]byte
}

func (store *Store) beInitialized() {
	if store.Bag != nil {
		return
	}
	store.Bag = make(map[string][]byte)
}

// Has implements go-ipld-prime/storage.Storage.Has.
func (store *Store) Has(ctx context.Context, key string) (bool, error) {
	if store.Bag == nil {
		return false, nil
	}
	_, exists := store.Bag[key]
	return exists, nil
}

// Get implements go-ipld-prime/storage.ReadableStorage.Get.
//
// Note that this internally performs a defensive copy;
// use Peek for higher performance if you are certain you won't mutate the returned slice.
func (store *Store) Get(ctx context.Context, key string) ([]byte, error) {
	store.beInitialized()
	content, exists := store.Bag[key]
	if !exists {
		return nil, fmt.Errorf("404") // FIXME this needs a standard error type
	}
	cpy := make([]byte, len(content))
	copy(cpy, content)
	return cpy, nil
}

// Put implements go-ipld-prime/storage.WritableStorage.Put.
func (store *Store) Put(ctx context.Context, key string, content []byte) error {
	store.beInitialized()
	if _, exists := store.Bag[key]; exists {
		return nil
	}
	cpy := make([]byte, len(content))
	copy(cpy, content)
	store.Bag[key] = cpy
	return nil
}

// GetStream implements go-ipld-prime/storage.StreamingReadableStorage.GetStream.
//
// It's useful for this storage implementation to explicitly support this,
// because returning a reader gives us room to avoid needing a defensive copy.
func (store *Store) GetStream(ctx context.Context, key string) (io.ReadCloser, error) {
	content, exists := store.Bag[key]
	if !exists {
		return nil, fmt.Errorf("404") // FIXME this needs a standard error type
	}
	return noopCloser{bytes.NewReader(content)}, nil
}

// Peek implements go-ipld-prime/storage.PeekableStorage.Peek.
func (store *Store) Peek(ctx context.Context, key string) ([]byte, io.Closer, error) {
	content, exists := store.Bag[key]
	if !exists {
		return nil, nil, fmt.Errorf("404") // FIXME this needs a standard error type
	}
	return content, noopCloser{nil}, nil
}

type noopCloser struct {
	io.Reader
}

func (noopCloser) Close() error { return nil }
