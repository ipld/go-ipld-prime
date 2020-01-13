package impls

// Map_K2_T2 and this file is how a codegen'd map type would work.  it's allowed to use concrete key and value types.
// In constrast with Map_K_T, this one has both complex keys and a struct for the value.

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

// --- we need some types to use for keys and values: --->
/*	ipldsch:
	type K2 struct { u string, i string } representation stringjoin (":")
	type T2 struct { a int, b int, c int, d int }
*/

// Note how we're not able to use `int` in the structs, but instead `plainInt`: this is so we can take address of those fields directly and return them as nodes.
//  We don't currently have concrete exported types that allow us to do this.  Maybe we should?

type K2 struct{ u, i plainString }
type T2 struct{ a, b, c, d plainInt }

func (K2) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (n *K2) LookupString(key string) (ipld.Node, error) {
	switch key {
	case "u":
		return &n.u, nil
	case "i":
		return &n.i, nil
	default:
		return nil, fmt.Errorf("no such field")
	}
}
func (n *K2) Lookup(key ipld.Node) (ipld.Node, error) {
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return n.LookupString(ks)
}
func (K2) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: "K2", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (n *K2) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupString(seg.String())
}
func (n *K2) MapIterator() ipld.MapIterator {
	return &_K2_MapIterator{n, 0}
}
func (K2) ListIterator() ipld.ListIterator {
	panic("no")
}
func (K2) Length() int {
	return -1
}
func (K2) IsUndefined() bool {
	return false
}
func (K2) IsNull() bool {
	return false
}
func (K2) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsString() (string, error) {
	return "", ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}
func (K2) Style() ipld.NodeStyle {
	panic("todo")
}

type _K2_MapIterator struct {
	n   *K2
	idx int
}

func (itr *_K2_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.idx >= 2 {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	switch itr.idx {
	case 0:
		k = plainString("u") // TODO: I guess we should generate const pools for struct field names?
		v = &itr.n.u
	case 1:
		k = plainString("i")
		v = &itr.n.i
	default:
		panic("unreachable")
	}
	itr.idx++
	return
}
func (itr *_K2_MapIterator) Done() bool {
	return itr.idx >= 2
}

func (T2) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (n *T2) LookupString(key string) (ipld.Node, error) {
	switch key {
	case "a":
		return &n.a, nil
	case "b":
		return &n.b, nil
	case "c":
		return &n.c, nil
	case "d":
		return &n.d, nil
	default:
		return nil, fmt.Errorf("no such field")
	}
}
func (n *T2) Lookup(key ipld.Node) (ipld.Node, error) {
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return n.LookupString(ks)
}
func (T2) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: "T2", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (n *T2) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupString(seg.String())
}
func (n *T2) MapIterator() ipld.MapIterator {
	return &_T2_MapIterator{n, 0}
}
func (T2) ListIterator() ipld.ListIterator {
	panic("no")
}
func (T2) Length() int {
	return -1
}
func (T2) IsUndefined() bool {
	return false
}
func (T2) IsNull() bool {
	return false
}
func (T2) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsString() (string, error) {
	return "", ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}
func (T2) Style() ipld.NodeStyle {
	panic("todo")
}

type _T2_MapIterator struct {
	n   *T2
	idx int
}

func (itr *_T2_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.idx >= 4 {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	switch itr.idx {
	case 0:
		k = plainString("a") // TODO: I guess we should generate const pools for struct field names?
		v = &itr.n.a
	case 1:
		k = plainString("b")
		v = &itr.n.b
	case 2:
		k = plainString("c")
		v = &itr.n.c
	case 3:
		k = plainString("d")
		v = &itr.n.d
	default:
		panic("unreachable")
	}
	itr.idx++
	return
}
func (itr *_T2_MapIterator) Done() bool {
	return itr.idx >= 4
}

type _K2__Assembler struct {
	w *K2
}
type _K2__ReprAssembler struct {
	w *K2
}

type _T2__Assembler struct {
	w *T2
}
type _T2__ReprAssembler struct {
	w *T2
}

func (ta *_T2__Assembler) BeginMap(_ int) (ipld.MapNodeAssembler, error) {
	return ta, nil
}
func (_T2__Assembler) BeginList(_ int) (ipld.ListNodeAssembler, error) { panic("no") }
func (_T2__Assembler) AssignNull() error                               { panic("no") }
func (_T2__Assembler) AssignBool(bool) error                           { panic("no") }
func (_T2__Assembler) AssignInt(int) error                             { panic("no") }
func (_T2__Assembler) AssignFloat(float64) error                       { panic("no") }
func (_T2__Assembler) AssignString(v string) error                     { panic("no") }
func (_T2__Assembler) AssignBytes([]byte) error                        { panic("no") }
func (ta *_T2__Assembler) Assign(v ipld.Node) error {
	if v2, ok := v.(*T2); ok {
		*ta.w = *v2
		return nil
	}
	// todo: apply a generic 'copy' function.
	panic("later")
}
func (_T2__Assembler) Style() ipld.NodeStyle { panic("later") }

func (ta *_T2__Assembler) AssembleKey() ipld.NodeAssembler {
	// this'll be fun
	panic("soon")
}
func (ta *_T2__Assembler) AssembleValue() ipld.NodeAssembler {
	// also fun
	panic("soon")
}
func (ta *_T2__Assembler) Done() error {
	panic("soon")
}
func (_T2__Assembler) KeyStyle() ipld.NodeStyle   { panic("later") }
func (_T2__Assembler) ValueStyle() ipld.NodeStyle { panic("later") }

// --- okay, now the type of interest: the map. --->
/*	ipldsch:
	type Root struct { mp {K2:T2} } # nevermind the root part, the anonymous map is the point.
*/

type Map_K2_T2 struct {
	m map[K2]*T2          // used for quick lookup.
	t []_Map_K2_T2__entry // used both for order maintainence, and for allocation amortization for both keys and values.
}

type _Map_K2_T2__entry struct {
	k K2 // address of this used when we return keys as nodes, such as in iterators.  Need in one place to amortize shifts to heap when ptr'ing for iface.
	v T2 // address of this is used in map values and to return.
}

func (n *Map_K2_T2) LookupString(key string) (ipld.Node, error) {
	panic("decision") // FIXME: What's this supposed to do?  does this error for maps with complex keys?
}

type _Map_K2_T2__Assembler struct {
	w  *Map_K2_T2
	ka _K2__Assembler
	va _T2__Assembler
}
type _Map_K2_T2__ReprAssembler struct {
	w  *Map_K2_T2
	ka _K2__ReprAssembler
	va _T2__ReprAssembler
}
