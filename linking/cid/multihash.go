package cidlink

import (
	"bytes"

	mh "github.com/multiformats/go-multihash"
)

func newMHHash(mhType uint64, length int) (*mhhash, error) {
	return &mhhash{mhType, length, new(bytes.Buffer)}, nil
}

type mhhash struct {
	mhType uint64
	mhLen  int
	*bytes.Buffer
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (mhh *mhhash) Sum(b []byte) []byte {
	fullBytes := mhh.Buffer.Bytes()
	if len(fullBytes) > 0 {
		fullBytes = append(fullBytes, b...)
	} else {
		fullBytes = b
	}
	sum, _ := mh.Digest(fullBytes, mhh.mhType, mhh.mhLen)
	return sum
}
