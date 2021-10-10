package storage

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
)

// TODO: move me

// Memory is a simple in-memory storage for data indexed by datamodel.Link.
// (It's little more than a map -- in fact, the map is exported,
// and you can poke it directly.)
//
// The OpenRead method conforms to linking.BlockReadOpener,
// and the OpenWrite method conforms to linking.BlockWriteOpener.
// Therefore it's easy to use in a LinkSystem like this:
//
//		store := storage.Memory{}
//		lsys.StorageReadOpener = (&store).OpenRead
//		lsys.StorageWriteOpener = (&store).OpenWrite
//
// This storage is mostly expected to be used for testing and demos,
// and as an example of how you can implement and integrate your own storage systems.
type Memory struct {
	Bag map[datamodel.Link][]byte
}

func (store *Memory) beInitialized() {
	if store.Bag != nil {
		return
	}
	store.Bag = make(map[datamodel.Link][]byte)
}

func (store *Memory) OpenRead(lnkCtx linking.LinkContext, lnk datamodel.Link) (io.Reader, error) {
	store.beInitialized()
	data, exists := store.Bag[lnk]
	if !exists {
		return nil, fmt.Errorf("404") // FIXME this needs a standard error type
	}
	return bytes.NewReader(data), nil
}

func (store *Memory) OpenWrite(lnkCtx linking.LinkContext) (io.Writer, linking.BlockWriteCommitter, error) {
	store.beInitialized()
	buf := bytes.Buffer{}
	return &buf, func(lnk datamodel.Link) error {
		store.Bag[lnk] = buf.Bytes()
		return nil
	}, nil
}
