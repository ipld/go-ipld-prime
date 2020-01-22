package impls

import (
	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

func Int(value int) ipld.Node {
	return plainInt(value)
}

// plainInt is a simple boxed int that complies with ipld.Node.
type plainInt int

func (plainInt) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Int
}
func (plainInt) LookupString(string) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) Lookup(key ipld.Node) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) MapIterator() ipld.MapIterator {
	panic("no")
}
func (plainInt) ListIterator() ipld.ListIterator {
	panic("no")
}
func (plainInt) Length() int {
	return -1
}
func (plainInt) IsUndefined() bool {
	return false
}
func (plainInt) IsNull() bool {
	return false
}
func (plainInt) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Int}
}
func (n plainInt) AsInt() (int, error) {
	return int(n), nil
}
func (plainInt) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) AsString() (string, error) {
	return "", ipld.ErrWrongKind{MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) Style() ipld.NodeStyle {
	panic("todo")
}
