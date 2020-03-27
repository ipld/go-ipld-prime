package gendemo

// Map_K_T and this file is how a codegen'd map type would work.  it's allowed to use concrete key and value types.

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
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
	return nil
}
func (K) ListIterator() ipld.ListIterator {
	return nil
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
	return nil
}
func (T) ListIterator() ipld.ListIterator {
	return nil
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

func (_K__Assembler) BeginMap(_ int) (ipld.MapAssembler, error)   { panic("no") }
func (_K__Assembler) BeginList(_ int) (ipld.ListAssembler, error) { panic("no") }
func (_K__Assembler) AssignNull() error                               { panic("no") }
func (_K__Assembler) AssignBool(bool) error                           { panic("no") }
func (_K__Assembler) AssignInt(v int) error                           { panic("no") }
func (_K__Assembler) AssignFloat(float64) error                       { panic("no") }
func (ta *_K__Assembler) AssignString(v string) error {
	ta.w.x = v
	return nil
}
func (_K__Assembler) AssignBytes([]byte) error { panic("no") }
func (ta *_K__Assembler) AssignNode(v ipld.Node) error {
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

func (_T__Assembler) BeginMap(_ int) (ipld.MapAssembler, error)   { panic("no") }
func (_T__Assembler) BeginList(_ int) (ipld.ListAssembler, error) { panic("no") }
func (_T__Assembler) AssignNull() error                           { panic("no") }
func (_T__Assembler) AssignBool(bool) error                       { panic("no") }
func (ta *_T__Assembler) AssignInt(v int) error {
	ta.w.x = v
	return nil
}
func (_T__Assembler) AssignFloat(float64) error   { panic("no") }
func (_T__Assembler) AssignString(v string) error { panic("no") }
func (_T__Assembler) AssignBytes([]byte) error    { panic("no") }
func (ta *_T__Assembler) AssignNode(v ipld.Node) error {
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
	return nil
}
func (n *Map_K_T) Length() int {
	return len(n.t)
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
	return Type__Map_K_T{}
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

// Type__Map_K_T  implements both schema.Type and ipld.NodeStyle.
//
// REVIEW: Should this just be exported?  I think probably yes.
// Alternatives: `Types().Map_K_T().NewBuilder()`; or, `Types` as a large const?
type Type__Map_K_T struct{}

func (Type__Map_K_T) NewBuilder() ipld.NodeBuilder {
	return &_Map_K_T__Builder{_Map_K_T__Assembler{
		w: &Map_K_T{},
	}}
}

// Overall assembly flow:
// - ('w' must already be set before beginning.)
// - BeginMap -- initializes the contents of the node in 'w'.
// - AssembleKey -- extends 'w.t', and sets up 'ka.ca.w' to point to the 'k' in the very tail of 'w.t'.
// - !branch:
//   - AssignString -- delegates to _K__Assembler (which may run validations); then checks for repeated key, errors if so; in case of either of those errors, un-extends 'w.t'.
//   - Assign -- more or less does one of the other two, above or below.
//   - BeginMap -- doesn't apply in this case (key is not complex), but if it was/did...
//     - Finish -- is implemented on _Map_K_T__KeyAssembler and delegates to _K__Assembler, because must do the repeated key check.
// - (okay, the key is now confirmed.  but, keep in mind: we still might need to back out if the value assignment errors.)
// - AssembleValue -- sets up 'va.ca.w' to point to the 'v' in the very tail of 'w.t'.
// - ...
//
// (Yep, basically any path through key *or* value assembly may error, and if they do,
// the parent has to roll back the last entry in 'w.t' -- so everything has a wrapper.)
type _Map_K_T__Assembler struct {
	w  *Map_K_T
	ka _Map_K_T__KeyAssembler
	va _Map_K_T__ValueAssembler

	state maState
}
type _Map_K_T__Builder struct {
	_Map_K_T__Assembler
}
type _Map_K_T__KeyAssembler struct {
	ma *_Map_K_T__Assembler // annoyingly cyclic but needed to do dupkey checks.
	ca _K__Assembler
}
type _Map_K_T__ValueAssembler struct {
	ma *_Map_K_T__Assembler // annoyingly cyclic but needed to reset the midappend state.
	ca _T__Assembler
}

func (nb *_Map_K_T__Builder) Build() ipld.Node {
	result := nb.w
	nb.w = nil
	return result
}
func (nb *_Map_K_T__Builder) Reset() {
	*nb = _Map_K_T__Builder{}
	nb.w = &Map_K_T{}
}

func (na *_Map_K_T__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
	// Allocate storage space.
	na.w.t = make([]_Map_K_T__entry, 0, sizeHint)
	na.w.m = make(map[K]*T, sizeHint)
	// Initialize the key and value assemblers with pointers back to the whole.
	na.ka.ma = na
	na.va.ma = na
	// That's it; return self as the MapAssembler.  We already have all the right methods on this structure.
	return na, nil
}
func (_Map_K_T__Assembler) BeginList(_ int) (ipld.ListAssembler, error) { panic("no") }
func (_Map_K_T__Assembler) AssignNull() error                               { panic("no") }
func (_Map_K_T__Assembler) AssignBool(bool) error                           { panic("no") }
func (_Map_K_T__Assembler) AssignInt(v int) error                           { panic("no") }
func (_Map_K_T__Assembler) AssignFloat(float64) error                       { panic("no") }
func (_Map_K_T__Assembler) AssignString(v string) error                     { panic("no") }
func (_Map_K_T__Assembler) AssignBytes([]byte) error                        { panic("no") }
func (_Map_K_T__Assembler) AssignLink(ipld.Link) error                      { panic("no") }
func (ta *_Map_K_T__Assembler) AssignNode(v ipld.Node) error {
	if v2, ok := v.(*Map_K_T); ok {
		*ta.w = *v2
		ta.w = nil // block further mutation
		return nil
	}
	// todo: apply a generic 'copy' function.
	panic("later")
}
func (_Map_K_T__Assembler) Style() ipld.NodeStyle { panic("later") }

func (ma *_Map_K_T__Assembler) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_midValue
	// Check for dup keys; error if so.
	_, exists := ma.w.m[K{k}]
	if exists {
		return nil, ipld.ErrRepeatedMapKey{&K{k}}
	}
	// Extend entry table and update map to point into the new row.
	l := len(ma.w.t)
	ma.w.t = append(ma.w.t, _Map_K_T__entry{k: K{k}})
	ma.w.m[K{k}] = &ma.w.t[l].v
	// Init the value assembler with a pointer to its target and yield it.
	ma.va.ca.w = &ma.w.t[l].v
	return &ma.va, nil
}

func (ma *_Map_K_T__Assembler) AssembleKey() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_midKey
	// Extend entry table.
	l := len(ma.w.t)
	ma.w.t = append(ma.w.t, _Map_K_T__entry{})
	// Init the key assembler with a pointer to its target and to whole 'ma' and yield it.
	ma.ka.ma = ma
	ma.ka.ca.w = &ma.w.t[l].k
	return &ma.ka
}
func (ma *_Map_K_T__Assembler) AssembleValue() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_expectValue {
		panic("misuse")
	}
	ma.state = maState_midValue
	// Init the value assembler with a pointer to its target and yield it.
	ma.va.ca.w = &ma.w.t[len(ma.w.t)-1].v
	return &ma.va
}
func (ma *_Map_K_T__Assembler) Finish() error {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_finished
	// validators could run and report errors promptly, if this type had any.
	return nil
}
func (_Map_K_T__Assembler) KeyStyle() ipld.NodeStyle           { panic("later") }
func (_Map_K_T__Assembler) ValueStyle(_ string) ipld.NodeStyle { panic("later") }

func (_Map_K_T__KeyAssembler) BeginMap(sizeHint int) (ipld.MapAssembler, error)   { panic("no") }
func (_Map_K_T__KeyAssembler) BeginList(sizeHint int) (ipld.ListAssembler, error) { panic("no") }
func (_Map_K_T__KeyAssembler) AssignNull() error                                      { panic("no") }
func (_Map_K_T__KeyAssembler) AssignBool(bool) error                                  { panic("no") }
func (_Map_K_T__KeyAssembler) AssignInt(int) error                                    { panic("no") }
func (_Map_K_T__KeyAssembler) AssignFloat(float64) error                              { panic("no") }
func (mka *_Map_K_T__KeyAssembler) AssignString(v string) error {
	// Check for dup keys; error if so.
	_, exists := mka.ma.w.m[K{v}]
	if exists {
		k := K{v}
		return ipld.ErrRepeatedMapKey{&k}
	}
	// Delegate to the key type's assembler.  It may run validations and may error.
	//  This results in the entry table memory being updated.
	//  When it returns, the delegated assembler should've already nil'd its 'w' to prevent further mutation.
	if err := mka.ca.AssignString(v); err != nil {
		return err // REVIEW:errors: probably deserves a wrapper indicating the error came during key coersion.
	}
	// Update the map to point into the entry value!
	//  (Hopefully the go compiler recognizes our assignment after existence check and optimizes appropriately.)
	mka.ma.w.m[K{v}] = &mka.ma.w.t[len(mka.ma.w.t)-1].v
	// Update parent assembler state: clear to proceed.
	mka.ma.state = maState_expectValue
	mka.ma = nil // invalidate self to prevent further incorrect use.
	return nil
}
func (_Map_K_T__KeyAssembler) AssignBytes([]byte) error   { panic("no") }
func (_Map_K_T__KeyAssembler) AssignLink(ipld.Link) error { panic("no") }
func (mka *_Map_K_T__KeyAssembler) AssignNode(v ipld.Node) error {
	vs, err := v.AsString()
	if err != nil {
		return fmt.Errorf("cannot assign non-string node into map key assembler") // FIXME:errors: this doesn't quite fit in ErrWrongKind cleanly; new error type?
	}
	return mka.AssignString(vs)
}
func (_Map_K_T__KeyAssembler) Style() ipld.NodeStyle { panic("later") } // probably should give the style of plainString, which could say "only stores string kind" (though we haven't made such a feature part of the interface yet).

func (mva *_Map_K_T__ValueAssembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
	panic("todo") // We would add the additional required methods to 'mva' to save another type... but in this case it's also clear to us at codegen time this method can just error.
}
func (mva *_Map_K_T__ValueAssembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
	panic("todo") // We would add the additional required methods to 'mva' to save another type... but in this case it's also clear to us at codegen time this method can just error.
}
func (mva *_Map_K_T__ValueAssembler) AssignNull() error     { panic("todo") } // All these scalar rejections also clear to us at codegen time.  We can report them without delegation.  Should?  Debatable; but will save SLOC.
func (mva *_Map_K_T__ValueAssembler) AssignBool(bool) error { panic("todo") }
func (mva *_Map_K_T__ValueAssembler) AssignInt(v int) error {
	if err := mva.ca.AssignInt(v); err != nil {
		return err
	}
	mva.flush()
	return nil
}
func (mva *_Map_K_T__ValueAssembler) AssignFloat(float64) error   { panic("todo") }
func (mva *_Map_K_T__ValueAssembler) AssignString(v string) error { panic("todo") }
func (mva *_Map_K_T__ValueAssembler) AssignBytes([]byte) error    { panic("todo") }
func (mva *_Map_K_T__ValueAssembler) AssignLink(ipld.Link) error  { panic("todo") }
func (mva *_Map_K_T__ValueAssembler) AssignNode(v ipld.Node) error {
	if err := mva.ca.AssignNode(v); err != nil {
		return err
	}
	mva.flush()
	return nil
}
func (mva *_Map_K_T__ValueAssembler) flush() {
	// The child assembler already assigned directly into the target memory,
	//  so there's not much to do here... except update the assembler state machine.
	// We also don't check the previous state because context makes us already sure:
	//  A) the appropriate time to do that would've been *before* assignments;
	//  A.2) accordingly, we did so before exposing this value assembler at all; and
	//  B) if we were in a wrong state because someone holds onto this too long,
	//   the invalidation we're about to do on `mva.ca.w` will make it impossible
	//    for them to make changes in appropriately.
	mva.ma.state = maState_initial
	mva.ca.w = nil
}
func (_Map_K_T__ValueAssembler) Style() ipld.NodeStyle { panic("later") }

// type _Map_K_T__ReprAssembler struct {
// 	w *Map_K_T
// 	// todo: ka _Map_K_T__KeyAssembler   ?  might need a different type for repr.
// 	// todo: va _Map_K_T__ValueAssembler ?  might need a different type for repr.
// }
