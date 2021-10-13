package cidlink

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
)

// Memory is a simple in-memory storage for cidlinks. It's the same as `storage.Memory`
// but uses typical multihash semantics used when reading/writing cidlinks.
type Memory struct {
	Bag map[string][]byte
}

func (store *Memory) beInitialized() {
	if store.Bag != nil {
		return
	}
	store.Bag = make(map[string][]byte)
}

func (store *Memory) OpenRead(lnkCtx linking.LinkContext, lnk datamodel.Link) (io.Reader, error) {
	store.beInitialized()
	cl, ok := lnk.(Link)
	if !ok {
		return nil, fmt.Errorf("incompatible link type: %T", lnk)
	}
	data, exists := store.Bag[string(cl.Hash())]
	if !exists {
		return nil, os.ErrNotExist
	}
	return bytes.NewReader(data), nil
}

func (store *Memory) OpenWrite(lnkCtx linking.LinkContext) (io.Writer, linking.BlockWriteCommitter, error) {
	store.beInitialized()
	buf := bytes.Buffer{}
	return &buf, func(lnk datamodel.Link) error {
		cl, ok := lnk.(Link)
		if !ok {
			return fmt.Errorf("incompatible link type: %T", lnk)
		}

		store.Bag[string(cl.Hash())] = buf.Bytes()
		return nil
	}, nil
}
