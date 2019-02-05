package ipld

import (
	"github.com/polydawn/refmt/shared"
)

// future consideration: TokenizableNode.PushTokens is one of the easiest
//  features in the whole area to code... but the least useful for composition.
//  we might drop it outright at some point.

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

type NodeUnmarshaller func(src shared.TokenSource) (Node, error)
