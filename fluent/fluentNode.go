package fluent

import (
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
	Kind() ipld.ReprKind
	TraverseField(path string) Node
	TraverseIndex(idx int) Node
	Keys() KeyIterator
	KeysImmediate() []string
	Length() int
	IsNull() bool
	AsBool() bool
	AsInt() int
	AsFloat() float64
	AsString() string
	AsBytes() []byte
	AsLink() ipld.Link
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

func (e Error) Error() string {
	return e.Err.Error()
}

func (n node) GetError() error {
	return n.err
}
func (n node) Kind() ipld.ReprKind {
	if n.err != nil {
		panic(Error{n.err})
	}
	return n.n.Kind()
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
func (n node) Keys() KeyIterator {
	if n.err != nil {
		panic(Error{n.err})
	}
	return &keyIterator{n.n.Keys()}
}
func (n node) KeysImmediate() []string {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.KeysImmediate()
	if err != nil {
		panic(Error{err})
	}
	return v
}
func (n node) Length() int {
	if n.err != nil {
		panic(Error{n.err})
	}
	return n.n.Length()
}
func (n node) IsNull() bool {
	if n.err != nil {
		panic(Error{n.err})
	}
	return n.n.IsNull()
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
func (n node) AsFloat() float64 {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.AsFloat()
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
func (n node) AsBytes() []byte {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.AsBytes()
	if err != nil {
		panic(Error{err})
	}
	return v
}
func (n node) AsLink() ipld.Link {
	if n.err != nil {
		panic(Error{n.err})
	}
	v, err := n.n.AsLink()
	if err != nil {
		panic(Error{err})
	}
	return v
}

type KeyIterator interface {
	Next() string
	HasNext() bool
}

type keyIterator struct {
	d ipld.KeyIterator
}

func (ki *keyIterator) Next() string {
	v, err := ki.d.Next()
	if err != nil {
		panic(Error{err})
	}
	return v
}
func (ki *keyIterator) HasNext() bool {
	return ki.d.HasNext()
}
