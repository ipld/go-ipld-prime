package dagcbor

import (
	"io"

	"github.com/polydawn/refmt/cbor"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

var (
	_ cidlink.MulticodecDecoder = Decoder
	_ cidlink.MulticodecEncoder = Encoder
)

func init() {
	cidlink.RegisterMulticodecDecoder(0x71, Decoder)
	cidlink.RegisterMulticodecEncoder(0x71, Encoder)
}

func Decoder(nb ipld.NodeBuilder, r io.Reader) (ipld.Node, error) {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.NodeBuilder supports this, for example.)
	type detectFastPath interface {
		DecodeDagCbor(io.Reader) (ipld.Node, error)
	}
	if nb2, ok := nb.(detectFastPath); ok {
		return nb2.DecodeDagCbor(r)
	}
	// Okay, generic builder path.
	return Unmarshal(nb, cbor.NewDecoder(cbor.DecodeOptions{}, r))
}

func Encoder(n ipld.Node, w io.Writer) error {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.Node supports this, for example.)
	type detectFastPath interface {
		EncodeDagCbor(io.Writer) error
	}
	if n2, ok := n.(detectFastPath); ok {
		return n2.EncodeDagCbor(w)
	}
	// Okay, generic inspection path.
	return Marshal(n, cbor.NewEncoder(w))
}
