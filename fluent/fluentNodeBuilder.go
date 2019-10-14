package fluent

import (
	"github.com/ipld/go-ipld-prime"
)

type NodeBuilder interface {
	CreateMap(MapBuildingClosure) ipld.Node
	AmendMap(MapBuildingClosure) ipld.Node
	CreateList(ListBuildingClosure) ipld.Node
	AmendList(ListBuildingClosure) ipld.Node
	CreateNull() ipld.Node
	CreateBool(bool) ipld.Node
	CreateInt(int) ipld.Node
	CreateFloat(float64) ipld.Node
	CreateString(string) ipld.Node
	CreateBytes([]byte) ipld.Node
	CreateLink(ipld.Link) ipld.Node
}

type MapBuilder interface {
	Insert(k, v ipld.Node)
	Delete(k ipld.Node)
}

type ListBuilder interface {
	AppendAll([]ipld.Node)
	Append(v ipld.Node)
	Set(idx int, v ipld.Node)
}

// MapBuildingClosure is the signiture of a function which builds a Node of kind map.
//
// The MapBuilder parameter is used to accumulate the new Node for the
// duration of the function; and when the function returns, that builder
// will be invoked.  (In other words, there's no need to call `Build` within
// the closure itself -- and correspondingly, note the lack of return value.)
//
// Additional NodeBuilder handles are provided for building keys and values.
// These are used when handling typed Node implementations, since in that
// case they may contain behavior related to the type contracts.
// (For untyped nodes, this is degenerate: these builders are not
// distinct from the parent builder driving this closure.)
//
// REVIEW : whether 'knb' is needed.  Not sure, and there are other pending
// discussions on this.  (It's mostly a concern about having a place to do
// validation on construction; but it's possible we can solve this without
// additional Nodes and Builders by making that validation the responsibility
// of the inside of the mb.Insert method; but will this compose well, and
// will it convey the appropriate times to do the validations correctly?
// Is 'knb' relevant even if that last question is 'no'?  If a concern is
// to avoid double-validations, that argues for `mb.Insert(Node, Node)` over
// `mb.Insert(string, Node)`, but if avoiding double-validations, that means
// we already have a Node and don't need 'knb' to get one.  ... Design!)
type MapBuildingClosure func(mb MapBuilder, knb NodeBuilder, vnb NodeBuilder)

// ListBuildingClosure is the signiture of a function which builds a Node of kind list.
//
// The ListBuilder parameter is used to accumulate the new Node for the
// duration of the function; and when the function returns, that builder
// will be invoked.  (In other words, there's no need to call `Build` within
// the closure itself -- and correspondingly, note the lack of return value.)
//
// Additional NodeBuilder handles are provided for building the values.
// These are used when handling typed Node implementations, since in that
// case they may contain behavior related to the type contracts.
// (For untyped nodes, this is degenerate: the 'vnb' builder is not
// distinct from the parent builder driving this closure.)
type ListBuildingClosure func(lb ListBuilder, vnb NodeBuilder)

func WrapNodeBuilder(nb ipld.NodeBuilder) NodeBuilder {
	return &nodeBuilder{nb}
}

type nodeBuilder struct {
	nb ipld.NodeBuilder
}

func (nb *nodeBuilder) CreateMap(fn MapBuildingClosure) ipld.Node {
	mb, err := nb.nb.CreateMap()
	if err != nil {
		panic(Error{err})
	}
	fn(mapBuilder{mb}, nb, nb) // FUTURE: check for typed.NodeBuilder; need to specialize latter params before calling down if so.
	n, err := mb.Build()
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) AmendMap(fn MapBuildingClosure) ipld.Node {
	mb, err := nb.nb.AmendMap()
	if err != nil {
		panic(Error{err})
	}
	fn(mapBuilder{mb}, nb, nb) // FUTURE: check for typed.NodeBuilder; need to specialize latter params before calling down if so.
	n, err := mb.Build()
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateList(fn ListBuildingClosure) ipld.Node {
	lb, err := nb.nb.CreateList()
	if err != nil {
		panic(Error{err})
	}
	fn(listBuilder{lb}, nb) // FUTURE: check for typed.NodeBuilder; need to specialize latter params before calling down if so.
	n, err := lb.Build()
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) AmendList(fn ListBuildingClosure) ipld.Node {
	lb, err := nb.nb.AmendList()
	if err != nil {
		panic(Error{err})
	}
	fn(listBuilder{lb}, nb) // FUTURE: check for typed.NodeBuilder; need to specialize latter params before calling down if so.
	n, err := lb.Build()
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateNull() ipld.Node {
	n, err := nb.nb.CreateNull()
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateBool(v bool) ipld.Node {
	n, err := nb.nb.CreateBool(v)
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateInt(v int) ipld.Node {
	n, err := nb.nb.CreateInt(v)
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateFloat(v float64) ipld.Node {
	n, err := nb.nb.CreateFloat(v)
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateString(v string) ipld.Node {
	n, err := nb.nb.CreateString(v)
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateBytes(v []byte) ipld.Node {
	n, err := nb.nb.CreateBytes(v)
	if err != nil {
		panic(Error{err})
	}
	return n
}
func (nb *nodeBuilder) CreateLink(v ipld.Link) ipld.Node {
	n, err := nb.nb.CreateLink(v)
	if err != nil {
		panic(Error{err})
	}
	return n
}

type mapBuilder struct {
	ipld.MapBuilder
}

func (mb mapBuilder) Insert(k, v ipld.Node) {
	if err := mb.MapBuilder.Insert(k, v); err != nil {
		if err != nil {
			panic(Error{err})
		}
	}
}
func (mb mapBuilder) Delete(k ipld.Node) {
	if err := mb.MapBuilder.Delete(k); err != nil {
		if err != nil {
			panic(Error{err})
		}
	}
}

type listBuilder struct {
	ipld.ListBuilder
}

func (lb listBuilder) AppendAll(vs []ipld.Node) {
	if err := lb.ListBuilder.AppendAll(vs); err != nil {
		if err != nil {
			panic(Error{err})
		}
	}
}
func (lb listBuilder) Append(v ipld.Node) {
	if err := lb.ListBuilder.Append(v); err != nil {
		if err != nil {
			panic(Error{err})
		}
	}
}
func (lb listBuilder) Set(idx int, v ipld.Node) {
	if err := lb.ListBuilder.Set(idx, v); err != nil {
		if err != nil {
			panic(Error{err})
		}
	}
}
