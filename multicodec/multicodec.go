package multicodec

import (
	"github.com/ipld/go-ipld-prime"
)

// EncoderRegistry is a simple map which maps a multicodec indicator number
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
var EncoderRegistry = make(map[uint64]ipld.Encoder)

// DecoderRegistry is a simple map which maps a multicodec indicator number
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
var DecoderRegistry = make(map[uint64]ipld.Decoder)
