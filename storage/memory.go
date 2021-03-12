package storage

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
)

// Memory is a simple in-memory storage for data indexed by ipld.Link.
// (It's little more than a map -- in fact, the map is exported,
// and you can poke it directly.)
//
// The OpenRead method conforms to ipld.BlockReadOpener,
// and the OpenWrite method conforms to ipld.BlockWriteOpener.
// Therefore it's easy to use in a LinkSystem like this:
//
//		store := storage.Memory{}
//		lsys.StorageReadOpener = (&store).OpenRead
//		lsys.StorageWriteOpener = (&store).OpenWrite
//
// This storage is mostly expected to be used for testing and demos,
// and as an example of how you can implement and integrate your own storage systems.
type Memory struct {
	Bag map[ipld.Link][]byte
}

func (store *Memory) beInitialized() {
	if store.Bag != nil {
		return
	}
	store.Bag = make(map[ipld.Link][]byte)
}

func (store *Memory) OpenRead(lnkCtx ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {
	store.beInitialized()
	data, exists := store.Bag[lnk]
	if !exists {
		return nil, fmt.Errorf("404") // FIXME this needs a standard error type
	}
	return bytes.NewReader(data), nil
}

func (store *Memory) OpenWrite(lnkCtx ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
	store.beInitialized()
	buf := bytes.Buffer{}
	return &buf, func(lnk ipld.Link) error {
		store.Bag[lnk] = buf.Bytes()
		return nil
	}, nil
}
