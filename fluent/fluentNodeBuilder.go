package fluent

import (
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

type NodeBuilder interface {
	CreateMap() MapBuilder
	AmendMap() MapBuilder
	CreateList() ListBuilder
	AmendList() ListBuilder
	CreateNull() ipld.Node
	CreateBool(bool) ipld.Node
	CreateInt(int) ipld.Node
	CreateFloat(float64) ipld.Node
	CreateString(string) ipld.Node
	CreateBytes([]byte) ipld.Node
	CreateLink(cid.Cid) ipld.Node
}

type MapBuilder interface {
	InsertAll(map[ipld.Node]ipld.Node) MapBuilder
	Insert(k, v ipld.Node) MapBuilder
	Delete(k ipld.Node) MapBuilder
	Build() ipld.Node
}

type ListBuilder interface {
	AppendAll([]ipld.Node) ListBuilder
	Append(v ipld.Node) ListBuilder
	Set(idx int, v ipld.Node) ListBuilder
	Build() ipld.Node
}

func WrapNodeBuilder(nb ipld.NodeBuilder) NodeBuilder {
	return &nodeBuilder{nb, nil}
}

type nodeBuilder struct {
	nb  ipld.NodeBuilder
	err error
}

func (nb *nodeBuilder) CreateMap() MapBuilder {
	return mapBuilder{nb.nb.CreateMap()}
}
func (nb *nodeBuilder) AmendMap() MapBuilder {
	return mapBuilder{nb.nb.AmendMap()}
}
func (nb *nodeBuilder) CreateList() ListBuilder {
	return listBuilder{nb.nb.CreateList()}
}
func (nb *nodeBuilder) AmendList() ListBuilder {
	return listBuilder{nb.nb.AmendList()}
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
func (nb *nodeBuilder) CreateLink(v cid.Cid) ipld.Node {
	n, err := nb.nb.CreateLink(v)
	if err != nil {
		panic(Error{err})
	}
	return n
}

type mapBuilder struct {
	ipld.MapBuilder
}

func (mb mapBuilder) InsertAll(vs map[ipld.Node]ipld.Node) MapBuilder {
	mb.MapBuilder.InsertAll(vs)
	return mb
}
func (mb mapBuilder) Insert(k, v ipld.Node) MapBuilder {
	mb.MapBuilder.Insert(k, v)
	return mb
}
func (mb mapBuilder) Delete(k ipld.Node) MapBuilder {
	mb.MapBuilder.Delete(k)
	return mb
}
func (mb mapBuilder) Build() ipld.Node {
	n, err := mb.MapBuilder.Build()
	if err != nil {
		panic(Error{err})
	}
	return n
}

type listBuilder struct {
	ipld.ListBuilder
}

func (lb listBuilder) AppendAll(vs []ipld.Node) ListBuilder {
	lb.ListBuilder.AppendAll(vs)
	return lb
}
func (lb listBuilder) Append(v ipld.Node) ListBuilder {
	lb.ListBuilder.Append(v)
	return lb
}
func (lb listBuilder) Set(idx int, v ipld.Node) ListBuilder {
	lb.ListBuilder.Set(idx, v)
	return lb
}
func (lb listBuilder) Build() ipld.Node {
	n, err := lb.ListBuilder.Build()
	if err != nil {
		panic(Error{err})
	}
	return n
}
