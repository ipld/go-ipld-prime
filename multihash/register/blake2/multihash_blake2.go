/*
	This package has no purpose except to perform registration of multihashes.

	It is meant to be used as a side-effecting import, e.g.

		import (
			_ "github.com/ipld/go-ipld-prime/mulithash/register/blake2"
		)

	This package registers several multihashes for the blake2 family
	(both the 's' and the 'b' variants, and in a variety of sizes).
*/
package blake2

import (
	"hash"

	"github.com/minio/blake2b-simd"
	"golang.org/x/crypto/blake2s"

	"github.com/ipld/go-ipld-prime/multihash"
)

const (
	BLAKE2B_MIN = 0xb201
	BLAKE2B_MAX = 0xb240
	BLAKE2S_MIN = 0xb241
	BLAKE2S_MAX = 0xb260
)

func init() {
	// BLAKE2S
	// This package only enables support for 32byte (256 bit) blake2s.
	multihash.Registry[BLAKE2S_MIN+31] = func() hash.Hash { h, _ := blake2s.New256(nil); return h }

	// BLAKE2B
	// There's a whole range of these.
	for c := uint64(BLAKE2B_MIN); c <= BLAKE2B_MAX; c++ {
		size := int(c - BLAKE2B_MIN + 1)
		multihash.Registry[c] = func() hash.Hash {
			hasher, err := blake2b.New(&blake2b.Config{Size: uint8(size)})
			if err != nil {
				panic(err)
			}
			return hasher
		}
	}
}
