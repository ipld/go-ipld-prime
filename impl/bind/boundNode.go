package ipldbind

import (
	"reflect"

	"github.com/ipld/go-ipld-prime"
	"github.com/polydawn/refmt/obj/atlas"
)

var (
	_ ipld.Node = &Node{}
)

/*
	Node binds to some Go object in memory, using the definitions provided
	by refmt's object atlasing tools.

	This binding does not provide a serialization valid for hashing; to
	compute a CID, you'll have to convert to another kind of node.
	If you're not sure which kind serializable node to use, try `ipldcbor.Node`.
*/
type Node struct {
	bound reflect.Value
	atlas atlas.Atlas
}
