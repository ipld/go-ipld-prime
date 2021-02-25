package multihash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// Registry is a simple map which maps a multihash indicator number
// to a standard golang Hash interface.
//
// Multihash indicator numbers are reserved and described in
// https://github.com/multiformats/multicodec/blob/master/table.csv .
// The keys used in this map must match those reservations.
//
// Hashers which are available in the golang stdlib are registered here automatically.
//
// Packages which want to register more hashing functions (and have a multihash number reserved!)
// are encouraged to do so at package init time.
// (Doing this at package init time ensures this map can be accessed without race conditions.)
//
// The linking/cid.DefaultLinkSystem will use this map to find hashers
// to use when serializing data and computing links,
// and when loading data from storage and verifying its integrity.
//
// This registry map is only used for default behaviors.
// If you don't want to rely on it, you can always construct your own LinkSystem.
// (For this reason, there's no special effort made to detect conflicting registrations in this map.
// If more than one package registers for the same multicodec indicator, and
// you somehow end up with both in your import tree, and yet care about which wins:
// then just don't use this registry anymore: make a LinkSystem that does what you need.)
// This should never be done to make behavior alterations
// (hash functions are well standardized and so is the multihash indicator table),
// but may be relevant if one is really itching to try out different hash implementations for performance reasons.
var Registry = make(map[uint64]func() hash.Hash)

func init() {
	Registry[0x00] = func() hash.Hash { return &identityMultihash{} }
	Registry[0xd5] = md5.New
	Registry[0x11] = sha1.New
	Registry[0x12] = sha256.New
	Registry[0x13] = sha512.New
	// Registry[0x1f] = sha256.New224 // SOON
	// Registry[0x20] = sha512.New384 // SOON
	Registry[0x56] = func() hash.Hash { return &doubleSha256{} }
}
