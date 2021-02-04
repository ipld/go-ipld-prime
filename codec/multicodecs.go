package codec

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"golang.org/x/crypto/sha3"

	"github.com/ipld/go-ipld-prime"
)

// MulticodecEncoderRegistry is a simple map which maps a multicodec indicator number
// to an ipld.Encoder function.
//
// Packages which implement an IPLD codec and have a multicodec number reserved in
// https://github.com/multiformats/multicodec/blob/master/table.csv
// are encouraged to register themselves in this map at package init time.
// (Doing this at package init time ensures this map can be accessed without race conditions.)
//
// The linking/cid.DefaultLinkSystem will use this map to find encoders
// to use when serializing data for linking and storage.
//
// This registry map is only used for default behaviors.
// If you don't want to rely on it, you can always construct your own LinkSystem.
// (For this reason, there's no special effort made to detect conflicting registrations in this map.
// If more than one package registers for the same multicodec indicator, and
// you somehow end up with both in your import tree, and yet care about which wins:
// then just don't use this registry anymore: make a LinkSystem that does what you need.)
var MulticodecEncoderRegistry = make(map[uint64]ipld.Encoder)

// MulticodecDecoderRegistry is a simple map which maps a multicodec indicator number
// to an ipld.Decoder function.
//
// Packages which implement an IPLD codec and have a multicodec number reserved in
// https://github.com/multiformats/multicodec/blob/master/table.csv
// are encouraged to register themselves in this map at package init time.
// (Doing this at package init time ensures this map can be accessed without race conditions.)
//
// The linking/cid.DefaultLinkSystem will use this map to find decoders
// to use when deserializing data from storage.
//
// This registry map is only used for default behaviors.
// If you don't want to rely on it, you can always construct your own LinkSystem.
// (For this reason, there's no special effort made to detect conflicting registrations in this map.
// If more than one package registers for the same multicodec indicator, and
// you somehow end up with both in your import tree, and yet care about which wins:
// then just don't use this registry anymore: make a LinkSystem that does what you need.)
var MulticodecDecoderRegistry = make(map[uint64]ipld.Decoder)

// MultihashRegistry is a simple map which maps a multihash indicator number
// to a standard golang Hash interface.
//
// Hashers which are available in the golang stdlib are registered here automatically.
// Some hashes from x/crypto are also included out-of-the-box.
//
// Packages which want to register more hashing functions and have a multihash number reserved in
// https://github.com/multiformats/multicodec/blob/master/table.csv
// are encouraged to do so at package init time.
// (Doing this at package init time ensures this map can be accessed without race conditions.)
//
// The linking/cid.DefaultLinkSystem will use this map to find decoders
// to use when deserializing data from storage.
//
// This registry map is only used for default behaviors.
// If you don't want to rely on it, you can always construct your own LinkSystem.
// (For this reason, there's no special effort made to detect conflicting registrations in this map.
// If more than one package registers for the same multicodec indicator, and
// you somehow end up with both in your import tree, and yet care about which wins:
// then just don't use this registry anymore: make a LinkSystem that does what you need.)
var MultihashRegistry = make(map[uint64]func() hash.Hash)

func init() {
	MultihashRegistry[0xd5] = md5.New
	MultihashRegistry[0x11] = sha1.New
	MultihashRegistry[0x12] = sha256.New
	MultihashRegistry[0x13] = sha512.New
	MultihashRegistry[0x14] = sha3.New512
	MultihashRegistry[0x15] = sha3.New384
	MultihashRegistry[0x16] = sha3.New256
	MultihashRegistry[0x17] = sha3.New224
}
