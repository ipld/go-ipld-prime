package basicnode

import (
	"fmt"
	"reflect"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

var (
	_ datamodel.Node                             = &plainList{}
	_ datamodel.NodePrototype                    = Prototype__List{}
	_ datamodel.NodePrototypeSupportingListAmend = Prototype__List{}
	_ datamodel.NodeBuilder                      = &plainList__Builder{}
	_ datamodel.NodeAssembler                    = &plainList__Assembler{}
)

// plainList is a concrete type that provides a list-kind datamodel.Node.
// It can contain any kind of value.
// plainList is also embedded in the 'any' struct and usable from there.
type plainList struct {
	x []datamodel.NodeAmender
}

// -- Node interface methods -->

func (plainList) Kind() datamodel.Kind {
	return datamodel.Kind_List
}
func (plainList) LookupByString(string) (datamodel.Node, error) {
	return mixins.List{TypeName: "list"}.LookupByString("")
}
func (plainList) LookupByNode(datamodel.Node) (datamodel.Node, error) {
	return mixins.List{TypeName: "list"}.LookupByNode(nil)
}
func (n *plainList) LookupByIndex(idx int64) (datamodel.Node, error) {
	if v, err := n.lookupAmenderByIndex(idx); err != nil {
		return nil, err
	} else {
		return v.Build(), nil
	}
}
func (n *plainList) lookupAmenderByIndex(idx int64) (datamodel.NodeAmender, error) {
	if n.Length() <= idx {
		return nil, datamodel.ErrNotExists{Segment: datamodel.PathSegmentOfInt(idx)}
	}
	return n.x[idx], nil
}
func (n *plainList) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	if v, err := n.lookupAmenderBySegment(seg); err != nil {
		return nil, err
	} else {
		return v.Build(), nil
	}
}
func (n *plainList) lookupAmenderBySegment(seg datamodel.PathSegment) (datamodel.NodeAmender, error) {
	idx, err := seg.Index()
	if err != nil {
		return nil, datamodel.ErrInvalidSegmentForList{TroubleSegment: seg, Reason: err}
	}
	return n.lookupAmenderByIndex(idx)
}
func (plainList) MapIterator() datamodel.MapIterator {
	return nil
}
func (n *plainList) ListIterator() datamodel.ListIterator {
	return &plainList_ListIterator{n, 0}
}
func (n *plainList) Length() int64 {
	return int64(len(n.x))
}
func (plainList) IsAbsent() bool {
	return false
}
func (plainList) IsNull() bool {
	return false
}
func (plainList) AsBool() (bool, error) {
	return mixins.List{TypeName: "list"}.AsBool()
}
func (plainList) AsInt() (int64, error) {
	return mixins.List{TypeName: "list"}.AsInt()
}
func (plainList) AsFloat() (float64, error) {
	return mixins.List{TypeName: "list"}.AsFloat()
}
func (plainList) AsString() (string, error) {
	return mixins.List{TypeName: "list"}.AsString()
}
func (plainList) AsBytes() ([]byte, error) {
	return mixins.List{TypeName: "list"}.AsBytes()
}
func (plainList) AsLink() (datamodel.Link, error) {
	return mixins.List{TypeName: "list"}.AsLink()
}
func (plainList) Prototype() datamodel.NodePrototype {
	return Prototype.List
}

type plainList_ListIterator struct {
	n   *plainList
	idx int
}

func (itr *plainList_ListIterator) Next() (idx int64, v datamodel.Node, _ error) {
	if itr.Done() {
		return -1, nil, datamodel.ErrIteratorOverread{}
	}
	v = itr.n.x[itr.idx].Build()
	idx = int64(itr.idx)
	itr.idx++
	return
}
func (itr *plainList_ListIterator) Done() bool {
	return itr.idx >= len(itr.n.x)
}

// -- NodePrototype -->

type Prototype__List struct{}

func (p Prototype__List) NewBuilder() datamodel.NodeBuilder {
	return p.AmendingBuilder(nil)
}

// -- NodePrototypeSupportingListAmend -->

func (p Prototype__List) AmendingBuilder(base datamodel.Node) datamodel.ListAmender {
	var w *plainList
	if base != nil {
		if baseList, castOk := base.(*plainList); !castOk {
			panic("misuse")
		} else {
			w = baseList
		}
	} else {
		w = &plainList{}
	}
	return &plainList__Builder{plainList__Assembler{w: w}}
}

// -- NodeBuilder -->

type plainList__Builder struct {
	plainList__Assembler
}

func (nb *plainList__Builder) Build() datamodel.Node {
	if (nb.state != laState_initial) && (nb.state != laState_finished) {
		panic("invalid state: assembly in progress must be 'finished' before Build can be called!")
	}
	return nb.w
}
func (nb *plainList__Builder) Reset() {
	*nb = plainList__Builder{}
	nb.w = &plainList{}
}

// -- NodeAmender -->

func (nb *plainList__Builder) Transform(path datamodel.Path, transform datamodel.AmendFn) (datamodel.Node, error) {
	// Can only transform the root of the node or an immediate child.
	if path.Len() > 1 {
		panic("misuse")
	}
	// Allow the root of the node to be replaced.
	if path.Len() == 0 {
		prevNode := nb.Build()
		if newNode, err := transform(prevNode); err != nil {
			return nil, err
		} else if newLb, castOk := newNode.(*plainList__Builder); !castOk {
			return nil, fmt.Errorf("transform: cannot transform root into incompatible type: %v", reflect.TypeOf(newLb))
		} else {
			*nb.w = *newLb.w
			return prevNode, nil
		}
	}
	childSeg, _ := path.Shift()
	childIdx, err := childSeg.Index()
	var childAmender datamodel.NodeAmender
	if err != nil {
		if childSeg.String() == "-" {
			// "-" indicates appending a new element to the end of the list.
			childIdx = nb.w.Length()
		} else {
			return nil, datamodel.ErrInvalidSegmentForList{TroubleSegment: childSeg, Reason: err}
		}
	} else {
		// Don't allow the index to be equal to the length if the segment was not "-".
		if childIdx >= nb.w.Length() {
			return nil, fmt.Errorf("transform: cannot navigate path segment %q at %q because it is beyond the list bounds", childSeg, path)
		}
		// Only lookup the segment if it was within range of the list elements. If `childIdx` is equal to the length of
		// the list, then we fall-through and append an element to the end of the list.
		childAmender, err = nb.w.lookupAmenderByIndex(childIdx)
		if err != nil {
			// Return any error other than "not exists"
			if _, notFoundErr := err.(datamodel.ErrNotExists); !notFoundErr {
				return nil, fmt.Errorf("transform: child at %q did not exist)", path)
			}
		}
	}
	// The default behaviour will be to update the element at the specified index (if it exists). New list elements can
	// be added in two cases:
	//  - If an element is being appended to the end of the list.
	//  - If the transformation of the target node results in a list of nodes, use the first node in the list to replace
	//    the target node and then "add" the rest after. This is a bit of an ugly hack but is required for compatibility
	//    with two conflicting sets of semantics - the current `focus` and `walk`, which (quite reasonably) do an
	//    in-place replacement of list elements, and JSON Patch (https://datatracker.ietf.org/doc/html/rfc6902), which
	//    does not specify list element replacement. The only "compliant" way to do this today is to first "remove" the
	//    target node and then "add" its replacement at the same index, which seems inefficient.
	var prevChildVal datamodel.Node = nil
	if childAmender != nil {
		prevChildVal = childAmender.Build()
	}
	if newChildVal, err := transform(prevChildVal); err != nil {
		return nil, err
	} else if newChildVal == nil {
		newX := make([]datamodel.NodeAmender, nb.w.Length()-1)
		copy(newX, nb.w.x[:childIdx])
		copy(newX[:childIdx], nb.w.x[childIdx+1:])
		nb.w.x = newX
	} else if err = nb.storeChildAmender(childIdx, newChildVal); err != nil {
		return nil, err
	}
	return prevChildVal, nil
}

func (nb *plainList__Builder) storeChildAmender(childIdx int64, a datamodel.NodeAmender) error {
	var elems []datamodel.NodeAmender
	n := a.Build()
	if (n.Kind() == datamodel.Kind_List) && (n.Length() > 0) {
		elems = make([]datamodel.NodeAmender, n.Length())
		// The following logic uses a transformed list (if there is one) to perform both insertions (needed by JSON
		// Patch) and replacements (needed by `focus` and `walk`), while also providing the flexibility to insert more
		// than one element at a particular index in the list.
		//
		// Rules:
		//   - If appending to the end of the main list, all elements from the transformed list will be individually
		//     appended to the end of the list.
		//   - If updating at a particular index in the main list, use the first element from the transformed list to
		//     replace the existing element at that index in the main list, then insert the rest of the transformed list
		//     elements after.
		//
		// A special case to consider is that of a list element genuinely being a list itself. If that is the case, the
		// transformation MUST wrap the element in another list so that, once unwrapped, the element can be replaced or
		// inserted without affecting its semantics. Otherwise, the sub-list's elements will get expanded onto that
		// index in the main list.
		for i := range elems {
			elem, err := n.LookupByIndex(int64(i))
			if err != nil {
				return err
			}
			elems[i] = Prototype.Any.AmendingBuilder(elem)
		}
	} else {
		elems = []datamodel.NodeAmender{Prototype.Any.AmendingBuilder(n)}
	}
	if childIdx == nb.w.Length() {
		nb.w.x = append(nb.w.x, elems...)
	} else {
		numElems := int64(len(elems))
		newX := make([]datamodel.NodeAmender, nb.w.Length()+numElems-1)
		copy(newX, nb.w.x[:childIdx])
		copy(newX[childIdx:], elems)
		copy(newX[childIdx+numElems:], nb.w.x[childIdx+1:])
		nb.w.x = newX
	}
	return nil
}

func (nb *plainList__Builder) Get(idx int64) (datamodel.Node, error) {
	return nb.w.LookupByIndex(idx)
}

func (nb *plainList__Builder) Remove(idx int64) error {
	_, err := nb.Transform(
		datamodel.NewPath([]datamodel.PathSegment{datamodel.PathSegmentOfInt(idx)}),
		func(_ datamodel.Node) (datamodel.NodeAmender, error) {
			return nil, nil
		},
	)
	return err
}

func (nb *plainList__Builder) Append(values datamodel.Node) error {
	// Passing an index equal to the length of the list will append the passed values to the end of the list
	return nb.Insert(nb.Length(), values)
}

func (nb *plainList__Builder) Insert(idx int64, values datamodel.Node) error {
	var ps datamodel.PathSegment
	if idx == nb.Length() {
		ps = datamodel.PathSegmentOfString("-") // indicates appending to the end of the list
	} else {
		ps = datamodel.PathSegmentOfInt(idx)
	}
	_, err := nb.Transform(
		datamodel.NewPath([]datamodel.PathSegment{ps}),
		func(_ datamodel.Node) (datamodel.NodeAmender, error) {
			return Prototype.Any.AmendingBuilder(values), nil
		},
	)
	return err
}

func (nb *plainList__Builder) Set(idx int64, value datamodel.Node) error {
	return nb.Insert(idx, value)
}

func (nb *plainList__Builder) Empty() bool {
	return nb.Length() == 0
}

func (nb *plainList__Builder) Length() int64 {
	return nb.w.Length()
}

func (nb *plainList__Builder) Clear() {
	nb.Reset()
}

func (nb *plainList__Builder) Values() (datamodel.Node, error) {
	return nb.Build(), nil
}

// -- NodeAssembler -->

type plainList__Assembler struct {
	w *plainList

	va plainList__ValueAssembler

	state laState
}
type plainList__ValueAssembler struct {
	la *plainList__Assembler
}

// laState is an enum of the state machine for a list assembler.
// (this might be something to export reusably, but it's also very much an impl detail that need not be seen, so, dubious.)
// it's similar to maState for maps, but has fewer states because we never have keys to assemble.
type laState uint8

const (
	laState_initial  laState = iota // also the 'expect value or finish' state
	laState_midValue                // waiting for a 'finished' state in the ValueAssembler.
	laState_finished                // 'w' will also be nil, but this is a politer statement
)

func (plainList__Assembler) BeginMap(sizeHint int64) (datamodel.MapAssembler, error) {
	return mixins.ListAssembler{TypeName: "list"}.BeginMap(0)
}
func (na *plainList__Assembler) BeginList(sizeHint int64) (datamodel.ListAssembler, error) {
	if sizeHint < 0 {
		sizeHint = 0
	}
	// Allocate storage space.
	na.w.x = make([]datamodel.NodeAmender, 0, sizeHint)
	// That's it; return self as the ListAssembler.  We already have all the right methods on this structure.
	return na, nil
}
func (plainList__Assembler) AssignNull() error {
	return mixins.ListAssembler{TypeName: "list"}.AssignNull()
}
func (plainList__Assembler) AssignBool(bool) error {
	return mixins.ListAssembler{TypeName: "list"}.AssignBool(false)
}
func (plainList__Assembler) AssignInt(int64) error {
	return mixins.ListAssembler{TypeName: "list"}.AssignInt(0)
}
func (plainList__Assembler) AssignFloat(float64) error {
	return mixins.ListAssembler{TypeName: "list"}.AssignFloat(0)
}
func (plainList__Assembler) AssignString(string) error {
	return mixins.ListAssembler{TypeName: "list"}.AssignString("")
}
func (plainList__Assembler) AssignBytes([]byte) error {
	return mixins.ListAssembler{TypeName: "list"}.AssignBytes(nil)
}
func (plainList__Assembler) AssignLink(datamodel.Link) error {
	return mixins.ListAssembler{TypeName: "list"}.AssignLink(nil)
}
func (na *plainList__Assembler) AssignNode(v datamodel.Node) error {
	// Sanity check, then update, assembler state.
	//  Update of state to 'finished' comes later; where exactly depends on if shortcuts apply.
	if na.state != laState_initial {
		panic("misuse")
	}
	// Copy the content.
	if v2, ok := v.(*plainList); ok { // if our own type: shortcut.
		// Copy the structure by value.
		//  This means we'll have pointers into the same internal maps and slices;
		//   this is okay, because the Node type promises it's immutable, and we are going to instantly finish ourselves to also maintain that.
		// FIXME: the shortcut behaves differently than the long way: it discards any existing progress.  Doesn't violate immut, but is odd.
		*na.w = *v2
		na.state = laState_finished
		return nil
	}
	// If the above shortcut didn't work, resort to a generic copy.
	//  We call AssignNode for all the child values, giving them a chance to hit shortcuts even if we didn't.
	if v.Kind() != datamodel.Kind_List {
		return datamodel.ErrWrongKind{TypeName: "list", MethodName: "AssignNode", AppropriateKind: datamodel.KindSet_JustList, ActualKind: v.Kind()}
	}
	itr := v.ListIterator()
	for !itr.Done() {
		_, v, err := itr.Next()
		if err != nil {
			return err
		}
		if err := na.AssembleValue().AssignNode(v); err != nil {
			return err
		}
	}
	return na.Finish()
}
func (plainList__Assembler) Prototype() datamodel.NodePrototype {
	return Prototype.List
}

// -- ListAssembler -->

// AssembleValue is part of conforming to ListAssembler, which we do on
// plainList__Assembler so that BeginList can just return a retyped pointer rather than new object.
func (la *plainList__Assembler) AssembleValue() datamodel.NodeAssembler {
	// Sanity check, then update, assembler state.
	if la.state != laState_initial {
		panic("misuse")
	}
	la.state = laState_midValue
	// Make value assembler valid by giving it pointer back to whole 'la'; yield it.
	la.va.la = la
	return &la.va
}

// Finish is part of conforming to ListAssembler, which we do on
// plainList__Assembler so that BeginList can just return a retyped pointer rather than new object.
func (la *plainList__Assembler) Finish() error {
	// Sanity check, then update, assembler state.
	if la.state != laState_initial {
		panic("misuse")
	}
	la.state = laState_finished
	// validators could run and report errors promptly, if this type had any.
	return nil
}
func (plainList__Assembler) ValuePrototype(_ int64) datamodel.NodePrototype {
	return Prototype.Any
}

// -- ListAssembler.ValueAssembler -->

func (lva *plainList__ValueAssembler) BeginMap(sizeHint int64) (datamodel.MapAssembler, error) {
	ma := plainList__ValueAssemblerMap{}
	ma.ca.w = &plainMap{}
	ma.p = lva.la
	_, err := ma.ca.BeginMap(sizeHint)
	return &ma, err
}
func (lva *plainList__ValueAssembler) BeginList(sizeHint int64) (datamodel.ListAssembler, error) {
	la := plainList__ValueAssemblerList{}
	la.ca.w = &plainList{}
	la.p = lva.la
	_, err := la.ca.BeginList(sizeHint)
	return &la, err
}
func (lva *plainList__ValueAssembler) AssignNull() error {
	return lva.AssignNode(datamodel.Null)
}
func (lva *plainList__ValueAssembler) AssignBool(v bool) error {
	vb := plainBool(v)
	return lva.AssignNode(&vb)
}
func (lva *plainList__ValueAssembler) AssignInt(v int64) error {
	vb := plainInt(v)
	return lva.AssignNode(&vb)
}
func (lva *plainList__ValueAssembler) AssignFloat(v float64) error {
	vb := plainFloat(v)
	return lva.AssignNode(&vb)
}
func (lva *plainList__ValueAssembler) AssignString(v string) error {
	vb := plainString(v)
	return lva.AssignNode(&vb)
}
func (lva *plainList__ValueAssembler) AssignBytes(v []byte) error {
	vb := plainBytes(v)
	return lva.AssignNode(&vb)
}
func (lva *plainList__ValueAssembler) AssignLink(v datamodel.Link) error {
	vb := plainLink{v}
	return lva.AssignNode(&vb)
}
func (lva *plainList__ValueAssembler) AssignNode(v datamodel.Node) error {
	lva.la.w.x = append(lva.la.w.x, Prototype.Any.AmendingBuilder(v))
	lva.la.state = laState_initial
	lva.la = nil // invalidate self to prevent further incorrect use.
	return nil
}
func (plainList__ValueAssembler) Prototype() datamodel.NodePrototype {
	return Prototype.Any
}

type plainList__ValueAssemblerMap struct {
	ca plainMap__Assembler
	p  *plainList__Assembler // pointer back to parent, for final insert and state bump
}

// we briefly state only the methods we need to delegate here.
// just embedding plainMap__Assembler also behaves correctly,
//  but causes a lot of unnecessary autogenerated functions in the final binary.

func (ma *plainList__ValueAssemblerMap) AssembleEntry(k string) (datamodel.NodeAssembler, error) {
	return ma.ca.AssembleEntry(k)
}
func (ma *plainList__ValueAssemblerMap) AssembleKey() datamodel.NodeAssembler {
	return ma.ca.AssembleKey()
}
func (ma *plainList__ValueAssemblerMap) AssembleValue() datamodel.NodeAssembler {
	return ma.ca.AssembleValue()
}
func (plainList__ValueAssemblerMap) KeyPrototype() datamodel.NodePrototype {
	return Prototype__String{}
}
func (plainList__ValueAssemblerMap) ValuePrototype(_ string) datamodel.NodePrototype {
	return Prototype.Any
}

func (ma *plainList__ValueAssemblerMap) Finish() error {
	if err := ma.ca.Finish(); err != nil {
		return err
	}
	w := ma.ca.w
	ma.ca.w = nil
	return ma.p.va.AssignNode(w)
}

type plainList__ValueAssemblerList struct {
	ca plainList__Assembler
	p  *plainList__Assembler // pointer back to parent, for final insert and state bump
}

// we briefly state only the methods we need to delegate here.
// just embedding plainList__Assembler also behaves correctly,
//  but causes a lot of unnecessary autogenerated functions in the final binary.

func (la *plainList__ValueAssemblerList) AssembleValue() datamodel.NodeAssembler {
	return la.ca.AssembleValue()
}
func (plainList__ValueAssemblerList) ValuePrototype(_ int64) datamodel.NodePrototype {
	return Prototype.Any
}

func (la *plainList__ValueAssemblerList) Finish() error {
	if err := la.ca.Finish(); err != nil {
		return err
	}
	w := la.ca.w
	la.ca.w = nil
	return la.p.va.AssignNode(w)
}
