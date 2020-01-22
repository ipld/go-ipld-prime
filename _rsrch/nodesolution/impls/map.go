package impls

// plainMap and this file is how ipldfree would do it: no concrete types, just interfaces.

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

type plainMap struct {
	m map[string]ipld.Node // string key -- even if a runtime schema wrapper is using us for storage, we must have a comparable type here, and string is all we know.
	t []plainMap__Entry
}

type plainMap__Entry struct {
	k plainString // address of this used when we return keys as nodes, such as in iterators.  Need in one place to amortize shifts to heap when ptr'ing for iface.
	v ipld.Node   // same as map values.  keeping them here simplifies iteration.  (in codegen'd maps, this position is also part of amortization, but in this implementation, that's less useful.)
	// todo: depends on what it... is.
	//  if we put the anyNode -- the whole dang steamrolled union -- here, then we get amortized values.
	//   it's unclear if that's worth it.  that's a really big struct, all told.  like... 9 words or something.
	//   we can make two different implementations which offer each choice to users, of course.  but I don't think many people will use such fine grained configurability; and which is appropriate might depend on some context and is a mid-tree decision; and all sorts of things that make it hard to make this ergonomically configurable.

	// on the bright side, we can at least get amortized keys for sure.  previous generation didn't have that.
	//  and that's probably going to translate to a significant improvement on selectors.  like... really significant.
}

func (plainMap) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (n *plainMap) LookupString(key string) (ipld.Node, error) {
	v, exists := n.m[key]
	if !exists {
		return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
	}
	return v, nil
}
func (n *plainMap) Lookup(key ipld.Node) (ipld.Node, error) {
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return n.LookupString(ks)
}
func (plainMap) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: "map", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (n *plainMap) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupString(seg.String())
}
func (n *plainMap) MapIterator() ipld.MapIterator {
	return &plainMap_MapIterator{n, 0}
}
func (plainMap) ListIterator() ipld.ListIterator {
	panic("no")
}
func (n *plainMap) Length() int {
	return len(n.t)
}
func (plainMap) IsUndefined() bool {
	return false
}
func (plainMap) IsNull() bool {
	return false
}
func (plainMap) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsString() (string, error) {
	return "", ipld.ErrWrongKind{TypeName: "map", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) Style() ipld.NodeStyle {
	panic("todo")
}

type plainMap_MapIterator struct {
	n   *plainMap
	idx int
}

func (itr *plainMap_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.Done() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	k = &itr.n.t[itr.idx].k
	v = itr.n.t[itr.idx].v
	itr.idx++
	return
}
func (itr *plainMap_MapIterator) Done() bool {
	return itr.idx >= len(itr.n.t)
}

type Style__Map struct{}

func (Style__Map) NewBuilder() ipld.NodeBuilder {
	return &plainMap__Builder{plainMap__Assembler{w: &plainMap{}}}
}

type plainMap__Assembler struct {
	w  *plainMap
	ka plainMap__KeyAssembler
	va plainMap__ValueAssembler

	midappend bool // if true, next call must be 'AssembleValue'.
}
type plainMap__Builder struct {
	plainMap__Assembler
}
type plainMap__KeyAssembler struct {
	ma *plainMap__Assembler
}
type plainMap__ValueAssembler struct {
	ma *plainMap__Assembler
}

func (nb *plainMap__Builder) Build() (ipld.Node, error) {
	// want to run validators here if present.
	// would also like to check for half-done work (e.g. midappend==true anywhere deeply, or not-Done maps, etc etc!).
	return nb.w, nil
}
func (nb *plainMap__Builder) Reset() {
	*nb = plainMap__Builder{}
	nb.w = &plainMap{}
}

func (na *plainMap__Assembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error) {
	// Allocate storage space.
	na.w.t = make([]plainMap__Entry, 0, sizeHint)
	na.w.m = make(map[string]ipld.Node, sizeHint)
	// Initialize the key and value assemblers with pointers back to the whole.
	//  (REVIEW: Should we initialize these during `AssembleKey` and `AssembleValue` calls, and nil them when done with each value?  It might increase the safety of use.)
	na.ka.ma = na
	na.va.ma = na
	// That's it; return self as the MapNodeAssembler.  We already have all the right methods on this structure.
	return na, nil
}
func (plainMap__Assembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) { panic("no") }
func (plainMap__Assembler) AssignNull() error                                      { panic("no") }
func (plainMap__Assembler) AssignBool(bool) error                                  { panic("no") }
func (plainMap__Assembler) AssignInt(int) error                                    { panic("no") }
func (plainMap__Assembler) AssignFloat(float64) error                              { panic("no") }
func (plainMap__Assembler) AssignString(v string) error                            { panic("no") }
func (plainMap__Assembler) AssignBytes([]byte) error                               { panic("no") }
func (na *plainMap__Assembler) Assign(v ipld.Node) error {
	// todo: apply a generic 'copy' function.
	// todo: probably can also shortcut to copying na.t and na.m if it's our same concrete type?
	//  (can't quite just `na.w = v`, because we don't have 'freeze' features, and we don't wanna open door to mutation of 'v'.)
	//   (wait... actually, probably we can?  'Assign' is a "done" method.  we can&should invalidate the wip pointer here.)
	panic("later")
}
func (plainMap__Assembler) Style() ipld.NodeStyle { panic("later") }

// NOT yet an interface function... but wanting benchmarking on it.
// (how much it might show up for structs is another question.)
//  (there's... enough differences vs things with concrete type knowledge we might wanna do this all with one of those, too.)
func (ma *plainMap__Assembler) AssembleDirectly(k string) (ipld.NodeAssembler, error) {
	if ma.midappend == true {
		panic("misuse")
	}
	_, exists := ma.w.m[k]
	if exists {
		return nil, ipld.ErrRepeatedMapKey{String(k)}
	}
	//l := len(ma.w.t)
	ma.w.t = append(ma.w.t, plainMap__Entry{k: plainString(k)})
	// configure and return an anyAssembler, similar to below in prepareAssigner
	panic("todo")
}

// plainMap__Assembler also directly implements MapAssembler, so BeginMap can just return a retyped pointer rather than new object.
func (ma *plainMap__Assembler) AssembleKey() ipld.NodeAssembler {
	if ma.midappend == true {
		panic("misuse")
	}
	ma.midappend = true
	ma.w.t = append(ma.w.t, plainMap__Entry{})
	return &ma.ka
}
func (ma *plainMap__Assembler) AssembleValue() ipld.NodeAssembler {
	if ma.midappend == false {
		panic("misuse")
	}
	return &ma.va
}
func (ma *plainMap__Assembler) Done() error {
	if ma.midappend == true {
		panic("misuse")
	}
	return nil
}
func (plainMap__Assembler) KeyStyle() ipld.NodeStyle   { panic("later") }
func (plainMap__Assembler) ValueStyle() ipld.NodeStyle { panic("later") }

func (plainMap__KeyAssembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error)   { panic("no") }
func (plainMap__KeyAssembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) { panic("no") }
func (plainMap__KeyAssembler) AssignNull() error                                      { panic("no") }
func (plainMap__KeyAssembler) AssignBool(bool) error                                  { panic("no") }
func (plainMap__KeyAssembler) AssignInt(int) error                                    { panic("no") }
func (plainMap__KeyAssembler) AssignFloat(float64) error                              { panic("no") }
func (mka *plainMap__KeyAssembler) AssignString(v string) error {
	_, exists := mka.ma.w.m[v]
	if exists {
		return ipld.ErrRepeatedMapKey{String(v)}
	}
	mka.ma.w.t[len(mka.ma.w.t)-1].k = plainString(v)
	return nil
}
func (plainMap__KeyAssembler) AssignBytes([]byte) error { panic("no") }
func (mka *plainMap__KeyAssembler) Assign(v ipld.Node) error {
	vs, err := v.AsString()
	if err != nil {
		return fmt.Errorf("cannot assign non-string node into map key assembler") // FIXME:errors: this doesn't quite fit in ErrWrongKind cleanly; new error type?
	}
	return mka.AssignString(vs)
}
func (plainMap__KeyAssembler) Style() ipld.NodeStyle { panic("later") } // probably should give the style of plainString, which could say "only stores string kind" (though we haven't made such a feature part of the interface yet).

func (mva *plainMap__ValueAssembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error) {
	panic("todo") // now please
}
func (mva *plainMap__ValueAssembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) {
	panic("todo") // now please
}
func (mva *plainMap__ValueAssembler) AssignNull() error     { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignBool(bool) error { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignInt(v int) error {
	l := len(mva.ma.w.t) - 1
	vb := plainInt(v)
	mva.ma.w.t[l].v = &vb
	mva.ma.w.m[string(mva.ma.w.t[l].k)] = &vb
	mva.ma.midappend = false
	return nil
}
func (mva *plainMap__ValueAssembler) AssignFloat(float64) error   { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignString(v string) error { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignBytes([]byte) error    { panic("todo") }
func (mva *plainMap__ValueAssembler) Assign(v ipld.Node) error {
	l := len(mva.ma.w.t) - 1
	mva.ma.w.t[l].v = v
	mva.ma.w.m[string(mva.ma.w.t[l].k)] = v
	mva.ma.midappend = false
	return nil
}
func (plainMap__ValueAssembler) Style() ipld.NodeStyle { panic("later") }
