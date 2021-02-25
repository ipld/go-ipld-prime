/*
	This package has no purpose except to perform registration of multihashes.

	It is meant to be used as a side-effecting import, e.g.

		import (
			_ "github.com/ipld/go-ipld-prime/mulithash/register/murmur3"
		)

	This package registers multihashes for the murmur3 family.
*/
package murmur3

// import (
// 	"github.com/gxed/hashland/murmur3"
//
// 	"github.com/ipld/go-ipld-prime/multihash"
// )

func init() {
	// REVIEW: what go-multihash has done historically is New32, but this doesn't match what the multihash table says, which is 128!
	// These are also very clearly noncryptographic functions and not suitable for content-addressing use (and would require writing adapters to qualify for hash.Hash), so I'm opting to... not.
	// multihash.Registry[0x22] = murmur3.New32
}
