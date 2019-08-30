package fluent

import (
	"github.com/ipld/go-ipld-prime"
)

// fluent.Node is an interface with all the same methods names as ipld.Node,
// but all of the methods return only the thing itself, and not error values --
// which makes chaining calls easier.
//
// The very first error value encountered will be stored, and can be viewed later.
// After an error is encountered, all subsequent lookup methods will
// silently return the same error-storing node.
// Any of the terminal scalar-returning methods will panic if an error is stored.
// (The fluent.Recover function can be used to nicely gather these panics.)
type Node interface {
	ReprKind() ipld.ReprKind
	LookupString(path string) Node
	Lookup(key Node) Node
	LookupIndex(idx int) Node
	MapIterator() MapIterator
	ListIterator() ListIterator
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
func (n node) ReprKind() ipld.ReprKind {
	if n.err != nil {
		panic(Error{n.err})
	}
	return n.n.ReprKind()
}
func (n node) LookupString(path string) Node {
	if n.err != nil {
		return n
	}
	v, err := n.n.LookupString(path)
	if err != nil {
		return node{nil, err}
	}
	return node{v, nil}
}
func (n node) Lookup(key Node) Node {
	if n.err != nil {
		return n
	}
	v, err := n.n.Lookup(key.(node).n) // hacky.  needs fluent.Node needs unbox method.
	if err != nil {
		return node{nil, err}
	}
	return node{v, nil}
}
func (n node) LookupIndex(idx int) Node {
	if n.err != nil {
		return n
	}
	v, err := n.n.LookupIndex(idx)
	if err != nil {
		return node{nil, err}
	}
	return node{v, nil}
}
func (n node) LookupSegment(seg ipld.PathSegment) Node {
	if n.err != nil {
		return n
	}
	v, err := n.n.LookupSegment(seg)
	if err != nil {
		return node{nil, err}
	}
	return node{v, nil}
}
func (n node) MapIterator() MapIterator {
	if n.err != nil {
		panic(Error{n.err})
	}
	return &mapIterator{n.n.MapIterator()}
}
func (n node) ListIterator() ListIterator {
	if n.err != nil {
		panic(Error{n.err})
	}
	return &listIterator{n.n.ListIterator()}
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

type MapIterator interface {
	Next() (key Node, value Node)
	Done() bool
}

type mapIterator struct {
	d ipld.MapIterator
}

func (itr *mapIterator) Next() (Node, Node) {
	k, v, err := itr.d.Next()
	if err != nil {
		panic(Error{err})
	}
	return node{k, nil}, node{v, nil}
}
func (itr *mapIterator) Done() bool {
	return itr.d.Done()
}

type ListIterator interface {
	Next() (idx int, value Node)
	Done() bool
}

type listIterator struct {
	d ipld.ListIterator
}

func (itr *listIterator) Next() (int, Node) {
	idx, v, err := itr.d.Next()
	if err != nil {
		panic(Error{err})
	}
	return idx, node{v, nil}
}
func (itr *listIterator) Done() bool {
	return itr.d.Done()
}
