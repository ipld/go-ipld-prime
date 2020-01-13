package impls

// Map_K_T and this file is how a codegen'd map type would work.  it's allowed to use concrete key and value types.

import (
	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

// --- we need some types to use for keys and values: --->
/*	ipldsch:
	type K string
	type T int
*/

type K struct{ x string }
type T struct{ x int }

func (K) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (K) LookupString(string) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (K) Lookup(key ipld.Node) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_String}
}
func (K) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_String}
}
func (K) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: ipld.ReprKind_String}
}
func (K) MapIterator() ipld.MapIterator {
	panic("no")
}
func (K) ListIterator() ipld.ListIterator {
	panic("no")
}
func (K) Length() int {
	return -1
}
func (K) IsUndefined() bool {
	return false
}
func (K) IsNull() bool {
	return false
}
func (K) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_String}
}
func (K) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_String}
}
func (K) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_String}
}
func (n *K) AsString() (string, error) {
	return n.x, nil
}
func (K) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_String}
}
func (K) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_String}
}
func (K) Style() ipld.NodeStyle {
	panic("todo")
}

func (T) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (T) LookupString(string) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_Int}
}
func (T) Lookup(key ipld.Node) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_Int}
}
func (T) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Int}
}
func (T) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: ipld.ReprKind_Int}
}
func (T) MapIterator() ipld.MapIterator {
	panic("no")
}
func (T) ListIterator() ipld.ListIterator {
	panic("no")
}
func (T) Length() int {
	return -1
}
func (T) IsUndefined() bool {
	return false
}
func (T) IsNull() bool {
	return false
}
func (T) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Int}
}
func (n *T) AsInt() (int, error) {
	return n.x, nil
}
func (T) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Int}
}
func (T) AsString() (string, error) {
	return "", ipld.ErrWrongKind{MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Int}
}
func (T) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Int}
}
func (T) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Int}
}
func (T) Style() ipld.NodeStyle {
	panic("todo")
}

type _K__Assembler struct {
	w *K
}
type _K__ReprAssembler struct {
	w *K
}

func (_K__Assembler) BeginMap(_ int) (ipld.MapNodeAssembler, error)   { panic("no") }
func (_K__Assembler) BeginList(_ int) (ipld.ListNodeAssembler, error) { panic("no") }
func (_K__Assembler) AssignNull() error                               { panic("no") }
func (_K__Assembler) AssignBool(bool) error                           { panic("no") }
func (_K__Assembler) AssignInt(v int) error                           { panic("no") }
func (_K__Assembler) AssignFloat(float64) error                       { panic("no") }
func (ta *_K__Assembler) AssignString(v string) error {
	ta.w.x = v
	return nil
}
func (_K__Assembler) AssignBytes([]byte) error { panic("no") }
func (ta *_K__Assembler) Assign(v ipld.Node) error {
	if v2, ok := v.(*K); ok {
		*ta.w = *v2
		return nil
	}
	v2, err := v.AsString()
	if err != nil {
		return err // TODO:errors: probably deserves a layer of decoration being more explicit about invalid assignment.
	}
	return ta.AssignString(v2)
}
func (_K__Assembler) Style() ipld.NodeStyle { panic("later") }

type _T__Assembler struct {
	w *T
}
type _T__ReprAssembler struct {
	w *T
}

func (_T__Assembler) BeginMap(_ int) (ipld.MapNodeAssembler, error)   { panic("no") }
func (_T__Assembler) BeginList(_ int) (ipld.ListNodeAssembler, error) { panic("no") }
func (_T__Assembler) AssignNull() error                               { panic("no") }
func (_T__Assembler) AssignBool(bool) error                           { panic("no") }
func (ta *_T__Assembler) AssignInt(v int) error {
	ta.w.x = v
	return nil
}
func (_T__Assembler) AssignFloat(float64) error   { panic("no") }
func (_T__Assembler) AssignString(v string) error { panic("no") }
func (_T__Assembler) AssignBytes([]byte) error    { panic("no") }
func (ta *_T__Assembler) Assign(v ipld.Node) error {
	if v2, ok := v.(*T); ok {
		*ta.w = *v2
		return nil
	}
	v2, err := v.AsInt()
	if err != nil {
		return err // TODO:errors: probably deserves a layer of decoration being more explicit about invalid assignment.
	}
	return ta.AssignInt(v2)
}
func (_T__Assembler) Style() ipld.NodeStyle { panic("later") }

// --- okay, now the type of interest: the map. --->
/*	ipldsch:
	type Root struct { mp {K:T} } # nevermind the root part, the anonymous map is the point.
*/

type Map_K_T struct {
	m map[K]*T          // used for quick lookup.
	t []_Map_K_T__entry // used both for order maintainence, and for allocation amortization for both keys and values.
}

type _Map_K_T__entry struct {
	k K // address of this used when we return keys as nodes, such as in iterators.  Need in one place to amortize shifts to heap when ptr'ing for iface.
	v T // address of this is used in map values and to return.
}

func (Map_K_T) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (n *Map_K_T) Get(key *K) (*T, error) {
	v, exists := n.m[*key]
	if !exists {
		return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(key.x)}
	}
	return v, nil
}
func (n *Map_K_T) LookupString(key string) (ipld.Node, error) {
	v, exists := n.m[K{key}]
	if !exists {
		return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
	}
	return v, nil
}
func (n *Map_K_T) Lookup(key ipld.Node) (ipld.Node, error) {
	if k2, ok := key.(*K); ok {
		return n.Get(k2)
	}
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return n.LookupString(ks)
}
func (Map_K_T) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: "Map_K_T", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (n *Map_K_T) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupString(seg.String())
}
func (n *Map_K_T) MapIterator() ipld.MapIterator {
	return &_Map_K_T_MapIterator{n, 0}
}
func (Map_K_T) ListIterator() ipld.ListIterator {
	panic("no")
}
func (Map_K_T) Length() int {
	return -1
}
func (Map_K_T) IsUndefined() bool {
	return false
}
func (Map_K_T) IsNull() bool {
	return false
}
func (Map_K_T) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{TypeName: "Map_K_T", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (Map_K_T) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{TypeName: "Map_K_T", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (Map_K_T) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{TypeName: "Map_K_T", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (Map_K_T) AsString() (string, error) {
	return "", ipld.ErrWrongKind{TypeName: "Map_K_T", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (Map_K_T) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{TypeName: "Map_K_T", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (Map_K_T) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{TypeName: "Map_K_T", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}
func (Map_K_T) Style() ipld.NodeStyle {
	panic("todo")
}

type _Map_K_T_MapIterator struct {
	n   *Map_K_T
	idx int
}

func (itr *_Map_K_T_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.Done() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	k = &itr.n.t[itr.idx].k
	v = &itr.n.t[itr.idx].v
	itr.idx++
	return
}
func (itr *_Map_K_T_MapIterator) Done() bool {
	return itr.idx >= len(itr.n.t)
}

type _Map_K_T__Assembler struct {
	w  *Map_K_T
	ka _K__Assembler
	va _T__Assembler

	midappend bool
}
type _Map_K_T__Builder struct {
	_Map_K_T__Assembler
}
type _Map_K_T__ReprAssembler struct {
	w  *Map_K_T
	ka _K__ReprAssembler
	va _T__ReprAssembler
}

func NewBuilder_Map_K_T() ipld.NodeBuilder {
	return &_Map_K_T__Builder{_Map_K_T__Assembler{w: &Map_K_T{}}}
}

func (nb *_Map_K_T__Builder) Build() (ipld.Node, error) {
	// want to run validators here if present.
	// would also like to check for half-done work (e.g. midappend==true anywhere deeply, or not-Done maps, etc etc!).
	return nb.w, nil
}
func (nb *_Map_K_T__Builder) Reset() {
	nb.w = &Map_K_T{}
}

func (ta *_Map_K_T__Assembler) BeginMap(_ int) (ipld.MapNodeAssembler, error) {
	return ta, nil
}
func (_Map_K_T__Assembler) BeginList(_ int) (ipld.ListNodeAssembler, error) { panic("no") }
func (_Map_K_T__Assembler) AssignNull() error                               { panic("no") }
func (_Map_K_T__Assembler) AssignBool(bool) error                           { panic("no") }
func (_Map_K_T__Assembler) AssignInt(v int) error                           { panic("no") }
func (_Map_K_T__Assembler) AssignFloat(float64) error                       { panic("no") }
func (_Map_K_T__Assembler) AssignString(v string) error                     { panic("no") }
func (_Map_K_T__Assembler) AssignBytes([]byte) error                        { panic("no") }
func (ta *_Map_K_T__Assembler) Assign(v ipld.Node) error {
	if v2, ok := v.(*Map_K_T); ok {
		*ta.w = *v2
		return nil
	}
	// todo: apply a generic 'copy' function.
	panic("later")
}
func (_Map_K_T__Assembler) Style() ipld.NodeStyle { panic("later") }

func (ma *_Map_K_T__Assembler) AssembleDirectly(k string) (ipld.NodeAssembler, error) {
	if ma.midappend == true {
		panic("misuse")
	}
	_, exists := ma.w.m[K{k}]
	if exists {
		return nil, ipld.ErrRepeatedMapKey{&K{k}}
	}
	l := len(ma.w.t)
	ma.w.t = append(ma.w.t, _Map_K_T__entry{k: K{k}})
	return &_T__Assembler{&ma.w.t[l].v}, nil
}

func (ma *_Map_K_T__Assembler) AssembleKey() ipld.NodeAssembler {
	if ma.midappend == true {
		panic("misuse")
	}
	ma.midappend = true
	l := len(ma.w.t)
	ma.w.t = append(ma.w.t, _Map_K_T__entry{})
	return &_K__Assembler{&ma.w.t[l].k}
}
func (ma *_Map_K_T__Assembler) AssembleValue() ipld.NodeAssembler {
	if ma.midappend == false {
		panic("misuse")
	}
	ma.midappend = false // REVIEW: kinda sketchy to set this so early... but there's only so much hand-holding we can do!
	return &_T__Assembler{&ma.w.t[len(ma.w.t)-1].v}
}
func (ta *_Map_K_T__Assembler) Done() error {
	if ta.midappend == true {
		panic("misuse")
	}
	// ... i... thought i was gonna need to do more work here.  i guess not.
	// validators could run and report errors promptly, if this type had any.
	return nil
}
func (_Map_K_T__Assembler) KeyStyle() ipld.NodeStyle   { panic("later") }
func (_Map_K_T__Assembler) ValueStyle() ipld.NodeStyle { panic("later") }
