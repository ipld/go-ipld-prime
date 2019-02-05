package ipld

import (
	"github.com/polydawn/refmt/shared"
)

// TokenizableNode is an optional interface which an ipld.Node may also
// implement to indicate that it has an efficient method for translating
// itself to a serial Token stream.
//
// Any ipld.Node is tokenizable via generic inspection, so providing
// this interface is optional and purely for performance improvement.
type TokenizableNode interface {
	// PushTokens converts this node and its children into a stream of refmt Tokens
	// and push them sequentially into the given TokenSink.
	// This is useful for serializing, or hashing, or feeding to any other
	// TokenSink (for example, converting to other ipld.Node implementations
	// which can construct themselves from a token stream).
	PushTokens(sink shared.TokenSink) error
}
