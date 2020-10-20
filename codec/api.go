package codec

import (
	"io"

	"github.com/ipld/go-ipld-prime"
)

// TODO: i still can't decide between marshaller vs encoder terminology.
//   i liked the defn i had in refmt: encoder handles tokens-to-bytes; marshaller handles trees to tokens.
//   but that distinction doesn't exist here.
//      if it did, we'd need a token type.
//      and in that case the encoder parts would be an internal codec code reuse choice, not necessary to expose.
//          if this was the case, it would suggest these functions should be called marshaller.
//   an alternate definition is: marshallers are things that twist a structure to a tokenizable form;
//    but going from tree (already trivially tokenizable) to serial is still considered an encoder.
//        i could also see this definition holding water, and it appears to be what i'm rolling with at the moment.
//
// maybe we really should make a TokenWalker thing.  Put it in codectools.
//   i still really don't know how we'd describe links in that, though.  it's really hard to claim links are a token.
//   maybe we can cram links into some sort of "extra" field in the token union.

// Encoder is the essential definition of a function that takes IPLD Data Model data in memory and serializes it.
// IPLD Codecs are written by implementing this function interface (as well as (typically) a matched Decoder).
//
// Encoder functions can be composed into an ipld.LinkSystem to provide
// a "one stop shop" API for handling content addressable storage.
// Encoder functions can also be used directly if you want to handle serial data streams.
//
// Most codec packages will have a ReusableEncoder type
// (which contains any working memory needed by the encoder implementation,
// as well as any configuration options),
// and that type will have an Encode function matching this interface.
//
// By convention, codec packages that have a multicodec contract will also have
// a package-scope exported function called Encode which also matches this interface,
// and is the equivalent of creating a zero-value ReusableEncoder (aka, default config)
// and using its Encode method.
// This package-scope function will typically also internally use a sync.Pool
// to keep some ReusableEncoder values on hand to avoid unnecesary allocations.
//
// Note that a ReusableEncoder type that supports configuration options
// does not functionally expose those options when invoked by the multicodec system --
// multicodec indicators do not provide room for extended configuration info.
// Codecs that expose configuration options are doing so for library users to enjoy;
// it does not mean those non-default configurations will necessarly be available
// in all scenarios that use codecs indirectly.
// There is also no standard interface for such configurations: by nature,
// if they exist at all, they vary per codec.
type Encoder func(data ipld.Node, output io.Writer) error

// Decoder is the essential definiton of a function that consumes serial data and unfurls it into IPLD Data Model-compatible in-memory representations.
// IPLD Codecs are written by implementing this function interface (as well as (typically) a matched Encoder).
//
// Decoder is the dual of Encoder.
// Most of the documentation for the Encoder function interface
// also applies wholesale to the Decoder interface.
type Decoder func(into ipld.NodeAssembler, input io.Reader) error

type ErrBudgetExhausted struct{}

func (e ErrBudgetExhausted) Error() string {
	return "decoder resource budget exhausted (message too long or too complex)"
}
