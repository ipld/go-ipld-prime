package ipld

var Null Node = nullNode{}
var Undef Node = undefNode{}

type nullNode struct{}

func (nullNode) ReprKind() ReprKind {
	return ReprKind_Null
}
func (nullNode) LookupString(key string) (Node, error) {
	return nil, ErrWrongKind{MethodName: "<null>.LookupString", AppropriateKind: ReprKindSet_JustMap, ActualKind: ReprKind_Null}
}
func (nullNode) Lookup(key Node) (Node, error) {
	return nil, ErrWrongKind{MethodName: "<null>.Lookup", AppropriateKind: ReprKindSet_JustMap, ActualKind: ReprKind_Null}
}
func (nullNode) LookupIndex(idx int) (Node, error) {
	return nil, ErrWrongKind{MethodName: "<null>.LookupIndex", AppropriateKind: ReprKindSet_JustList, ActualKind: ReprKind_Null}
}
func (nullNode) MapIterator() MapIterator {
	return mapIteratorReject{ErrWrongKind{MethodName: "<null>.MapIterator", AppropriateKind: ReprKindSet_JustMap, ActualKind: ReprKind_Null}}
}
func (nullNode) ListIterator() ListIterator {
	return listIteratorReject{ErrWrongKind{MethodName: "<null>.ListIterator", AppropriateKind: ReprKindSet_JustList, ActualKind: ReprKind_Null}}
}
func (nullNode) Length() int {
	return -1
}
func (nullNode) IsUndefined() bool {
	return false
}
func (nullNode) IsNull() bool {
	return true
}
func (nullNode) AsBool() (bool, error) {
	return false, ErrWrongKind{MethodName: "<null>.AsBool", AppropriateKind: ReprKindSet_JustBool, ActualKind: ReprKind_Null}
}
func (nullNode) AsInt() (int, error) {
	return 0, ErrWrongKind{MethodName: "<null>.AsInt", AppropriateKind: ReprKindSet_JustInt, ActualKind: ReprKind_Null}
}
func (nullNode) AsFloat() (float64, error) {
	return 0, ErrWrongKind{MethodName: "<null>.AsFloat", AppropriateKind: ReprKindSet_JustFloat, ActualKind: ReprKind_Null}
}
func (nullNode) AsString() (string, error) {
	return "", ErrWrongKind{MethodName: "<null>.AsString", AppropriateKind: ReprKindSet_JustString, ActualKind: ReprKind_Null}
}
func (nullNode) AsBytes() ([]byte, error) {
	return nil, ErrWrongKind{MethodName: "<null>.AsBytes", AppropriateKind: ReprKindSet_JustBytes, ActualKind: ReprKind_Null}
}
func (nullNode) AsLink() (Link, error) {
	return nil, ErrWrongKind{MethodName: "<null>.AsLink", AppropriateKind: ReprKindSet_JustLink, ActualKind: ReprKind_Null}
}
func (nullNode) NodeBuilder() NodeBuilder {
	panic("cannot build null nodes")
}

type undefNode struct{ nullNode }

func (undefNode) IsUndefined() bool {
	return true
}
func (undefNode) IsNull() bool {
	return false
}

type mapIteratorReject struct{ err error }
type listIteratorReject struct{ err error }

func (itr mapIteratorReject) Next() (Node, Node, error) { return nil, nil, itr.err }
func (itr mapIteratorReject) Done() bool                { return false }

func (itr listIteratorReject) Next() (int, Node, error) { return -1, nil, itr.err }
func (itr listIteratorReject) Done() bool               { return false }
