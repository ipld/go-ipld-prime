package ipldfree

import (
	ipld "github.com/ipld/go-ipld-prime"
	nodeutil "github.com/ipld/go-ipld-prime/impl/util"
)

func String(value string) ipld.Node {
	return justString(value)
}

// justString is a simple boxed string that complies with ipld.Node.
// It's useful for many things, such as boxing map keys.
//
// The implementation is a simple typedef of a string;
// handling it as a Node incurs 'runtime.convTstring',
// which is about the best we can do.
type justString string

func (justString) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (justString) LookupString(string) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (justString) Lookup(key ipld.Node) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (justString) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
}
func (justString) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: ipld.ReprKind_String}
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
	return string(x), nil
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
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) AmendMap() (ipld.MapBuilder, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.AmendMap", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateList() (ipld.ListBuilder, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) AmendList() (ipld.ListBuilder, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.AmendList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateNull() (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateNull", AppropriateKind: ipld.ReprKindSet_JustNull, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateBool(v bool) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateInt(v int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateFloat(v float64) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateString(v string) (ipld.Node, error) {
	return justString(v), nil
}
func (nb justStringNodeBuilder) CreateBytes(v []byte) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_String}
}
func (nb justStringNodeBuilder) CreateLink(v ipld.Link) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "NodeBuilder.CreateLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_String}
}
