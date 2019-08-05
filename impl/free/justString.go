package ipldfree

import (
	ipld "github.com/ipld/go-ipld-prime"
	nodeutil "github.com/ipld/go-ipld-prime/impl/util"
)

func String(value string) ipld.Node {
	return justString{value}
}

// justString is a simple boxed string that complies with ipld.Node.
// It doesn't actually contain type info or comply with typed.Node
// (which makes it cheaper: this struct doesn't trigger 'convt2e').
// justString is particularly useful for boxing things like struct keys.
type justString struct {
	x string
}

// FUTURE: we'll also want a typed string, of course.
//  Looking forward to benchmarking how that shakes out: it will almost
//   certainly add cost in the form of 'convt2e', but we'll see how much.
//    It'll also be particularly interesting to find out if common patterns of
//     usage around map iterators will get the compiler to skip that cost if
//      the key is unused by the caller.

func (justString) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (justString) TraverseField(string) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "TraverseField", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (justString) TraverseIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "TraverseIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
}
func (justString) MapIterator() ipld.MapIterator {
	return nodeutil.MapIteratorErrorThunk(ipld.ErrWrongKind{MethodName: "MapIterator", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String})
}
func (justString) ListIterator() ipld.ListIterator {
	return nodeutil.ListIteratorErrorThunk(ipld.ErrWrongKind{MethodName: "ListIterator", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String})
}
func (justString) Length() int {
	return -1
}
func (justString) IsUndefined() bool {
	return false
}
func (justString) IsNull() bool {
	return false
}
func (justString) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_String}
}
func (justString) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_String}
}
func (justString) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_String}
}
func (x justString) AsString() (string, error) {
	return x.x, nil
}
func (justString) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_String}
}
func (justString) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_String}
}
func (justString) NodeBuilder() ipld.NodeBuilder {
	return justStringNodeBuilder{}
}

type justStringNodeBuilder struct{}

func (nb justStringNodeBuilder) CreateMap() (ipld.MapBuilder, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) AmendMap() (ipld.MapBuilder, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AmendMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateList() (ipld.ListBuilder, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) AmendList() (ipld.ListBuilder, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AmendList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateNull() (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateNull", AppropriateKind: ipld.ReprKindSet_JustNull, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateBool(v bool) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateInt(v int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateFloat(v float64) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateString(v string) (ipld.Node, error) {
	return justString{v}, nil
}
func (nb justStringNodeBuilder) CreateBytes(v []byte) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateLink(v ipld.Link) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "CreateLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_String}
}
