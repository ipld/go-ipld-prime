package repose

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/polydawn/refmt/cbor"
)

var (
	_ MulticodecDecoder = DecoderDagCbor
	_ MulticodecEncoder = EncoderDagCbor
)

func DecoderDagCbor(nb ipld.NodeBuilder, r io.Reader) (ipld.Node, error) {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.NodeBuilder supports this, for example.)
	type detectFastPath interface {
		DecodeCbor(io.Reader) (ipld.Node, error)
	}
	if nb2, ok := nb.(detectFastPath); ok {
		return nb2.DecodeCbor(r)
	}
	// Okay, generic builder path.
	return encoding.Unmarshal(nb, cbor.NewDecoder(r))
}

func EncoderDagCbor(n ipld.Node, w io.Writer) error {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.Node supports this, for example.)
	type detectFastPath interface {
		EncodeCbor(io.Writer) error
	}
	if n2, ok := n.(detectFastPath); ok {
		return n2.EncodeCbor(w)
	}
	// Okay, generic inspection path.
	return encoding.Marshal(n, cbor.NewEncoder(w))
}
