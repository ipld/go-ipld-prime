package ipldcbor

import (
	"github.com/ipld/go-ipld-prime"
)

var (
	_ ipld.Node             = &Node{}
	_ ipld.SerializableNode = &Node{}
)

/*
	Node in ipldcbor is implemented // ... if you want to hash it or write it: both are better streamed (especially the former)
*/
type Node struct {
}
