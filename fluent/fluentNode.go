package fluent

import (
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

// fluent.Node is an interface with all the same methods names as ipld.Node,
// but all of the methods return only the thing itself, and not error values --
// which makes chaining calls easier.
//
// The very first error value encountered will be stored, and can be viewed later.
// After an error is encountered, all subsequent traversal methods will
// silently return the same error-storing node.
// Any of the terminal scalar-returning methods will panic if an error is stored.
// (The fluent.Recover function can be used to nicely gather these panics.)
type Node interface {
	TraverseField(path string) Node
	TraverseIndex(idx int) Node
	AsBool() bool
	AsString() string
	AsInt() int
	AsLink() cid.Cid
	GetError() error
}

func WrapNode(n ipld.Node) Node {
	return node{n, nil}
}

type node struct {
	n   ipld.Node
	err error
}

type Error struct {
	Err error
}

func (n node) GetError() error {
	return n.err
}
func (n node) TraverseField(path string) Node {
	if n.err != nil {
		return n
	}
	v, err := n.n.TraverseField(path)
	if err != nil {
		return node{nil, err}
	}
	return node{v, nil}
}
func (n node) TraverseIndex(idx int) Node {
	if n.err != nil {
		return n
	}
	v, err := n.n.TraverseIndex(idx)
	if err != nil {
		return node{nil, err}
	}
	return node{v, nil}
}
func (n node) AsBool() bool {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.AsBool()
	if err != nil {
		panic(Error{err})
	}
	return v
}
func (n node) AsString() string {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.AsString()
	if err != nil {
		panic(Error{err})
	}
	return v
}
func (n node) AsInt() int {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.AsInt()
	if err != nil {
		panic(Error{err})
	}
	return v
}
func (n node) AsLink() cid.Cid {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.AsLink()
	if err != nil {
		panic(Error{err})
	}
	return v
}
