package repose

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type MulticodecDecodeTable struct {
	Table map[uint64]MulticodecDecoder
}

type MulticodecEncodeTable struct {
	Table map[uint64]MulticodecEncoder
}

// Almost certainly curries an ipld.NodeBuilder inside itself; or some other
// more specialized thing of the same purpose.
type MulticodecDecoder func(io.Reader) (ipld.Node, error)

// Tends to be implemented by probing the node to see if it matches a special
// interface that we know can do this particular kind of encoding
// (e.g. if you're using ipldgit.Node and making a MulticodecEncoder to register
// as the rawgit multicodec, you'll probe for that specific thing, since it's
// implemented on the node itself),
// but may also be able to work based on the ipld.Node interface alone
// (e.g. you can do dag-cbor to any kind of Node).
type MulticodecEncoder func(ipld.Node, io.Writer) error
