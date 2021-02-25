package multihash

import (
	"bytes"
	"crypto/sha256"
	"hash"
)

type identityMultihash struct {
	bytes.Buffer
}

func (identityMultihash) BlockSize() int {
	return 32 // A prefered block size is nonsense for the "identity" "hash".  An arbitrary but unsurprising and positive nonzero number has been chosen to minimize the odds of fascinating bugs.
}

func (x identityMultihash) Size() int {
	return x.Len()
}

func (x identityMultihash) Sum(digest []byte) []byte {
	return x.Bytes()
}

type doubleSha256 struct {
	main hash.Hash
}

func (x doubleSha256) Write(body []byte) (int, error) {
	return x.main.Write(body)
}

func (doubleSha256) BlockSize() int {
	return sha256.BlockSize
}

func (doubleSha256) Size() int {
	return sha256.Size
}

func (x doubleSha256) Reset() {
	x.main.Reset()
}

func (x doubleSha256) Sum(digest []byte) []byte {
	intermediate := [sha256.Size]byte{}
	x.main.Sum(intermediate[:])
	h2 := sha256.New()
	h2.Write(intermediate[:])
	return h2.Sum(digest)
}
