package amend

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

var (
	_ datamodel.Node = &anyAmender{}
	_ Amender        = &anyAmender{}
)

type anyAmender struct {
	base    datamodel.Node
	parent  Amender
	created bool
}

func newAnyAmender(base datamodel.Node, parent Amender, create bool) Amender {
	return &anyAmender{base, parent, create}
}

func (a *anyAmender) Build() datamodel.Node {
	// `anyAmender` is also a `Node`.
	return (datamodel.Node)(a)
}

func (a *anyAmender) Kind() datamodel.Kind {
	return a.base.Kind()
}

func (a *anyAmender) LookupByString(key string) (datamodel.Node, error) {
	return a.base.LookupByString(key)
}

func (a *anyAmender) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	return a.base.LookupByNode(key)
}

func (a *anyAmender) LookupByIndex(idx int64) (datamodel.Node, error) {
	return a.base.LookupByIndex(idx)
}

func (a *anyAmender) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	return a.base.LookupBySegment(seg)
}

func (a *anyAmender) MapIterator() datamodel.MapIterator {
	return a.base.MapIterator()
}

func (a *anyAmender) ListIterator() datamodel.ListIterator {
	return a.base.ListIterator()
}

func (a *anyAmender) Length() int64 {
	return a.base.Length()
}

func (a *anyAmender) IsAbsent() bool {
	return a.base.IsAbsent()
}

func (a *anyAmender) IsNull() bool {
	return a.base.IsNull()
}

func (a *anyAmender) AsBool() (bool, error) {
	return a.base.AsBool()
}

func (a *anyAmender) AsInt() (int64, error) {
	return a.base.AsInt()
}

func (a *anyAmender) AsFloat() (float64, error) {
	return a.base.AsFloat()
}

func (a *anyAmender) AsString() (string, error) {
	return a.base.AsString()
}

func (a *anyAmender) AsBytes() ([]byte, error) {
	return a.base.AsBytes()
}

func (a *anyAmender) AsLink() (datamodel.Link, error) {
	return a.base.AsLink()
}

func (a *anyAmender) Prototype() datamodel.NodePrototype {
	return a.base.Prototype()
}

func (a *anyAmender) Get(path datamodel.Path) (datamodel.Node, error) {
	// If the base node is an amender, use it, otherwise panic.
	if amd, castOk := a.base.(Amender); castOk {
		return amd.Get(path)
	}
	panic("misuse")
}

func (a *anyAmender) Add(path datamodel.Path, value datamodel.Node, createParents bool) error {
	// If the base node is an amender, use it, otherwise panic.
	if amd, castOk := a.base.(Amender); castOk {
		return amd.Add(path, value, createParents)
	}
	panic("misuse")
}

func (a *anyAmender) Remove(path datamodel.Path) (datamodel.Node, error) {
	// If the base node is an amender, use it, otherwise panic.
	if amd, castOk := a.base.(Amender); castOk {
		return amd.Remove(path)
	}
	panic("misuse")
}

func (a *anyAmender) Replace(path datamodel.Path, value datamodel.Node) (datamodel.Node, error) {
	// If the base node is an amender, use it, otherwise panic.
	if amd, castOk := a.base.(Amender); castOk {
		return amd.Replace(path, value)
	}
	panic("misuse")
}
