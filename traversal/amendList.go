package traversal

import (
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

var (
	_ datamodel.Node = &listAmender{}
	_ Amender        = &listAmender{}
)

type listElement struct {
	baseIdx int
	elem    datamodel.Node
}

type listAmender struct {
	cfg     *AmendOptions
	base    datamodel.Node
	parent  Amender
	created bool
	mods    arraylist.List
}

func (cfg *AmendOptions) newListAmender(base datamodel.Node, parent Amender, create bool) Amender {
	// If the base node is already a list-amender, reuse the mutation state but reset `parent` and `created`.
	if amd, castOk := base.(*listAmender); castOk {
		return &listAmender{cfg, amd.base, parent, create, amd.mods}
	} else {
		// Start with fresh state because existing metadata could not be reused.
		var elems []interface{}
		if base != nil {
			elems = make([]interface{}, base.Length())
			for i := range elems {
				elems[i] = listElement{i, nil}
			}
		} else {
			elems = make([]interface{}, 0)
		}
		return &listAmender{cfg, base, parent, create, *arraylist.New(elems...)}
	}
}

func (a *listAmender) Build() datamodel.Node {
	// `listAmender` is also a `Node`.
	return (datamodel.Node)(a)
}

func (a *listAmender) Kind() datamodel.Kind {
	return datamodel.Kind_List
}

func (a *listAmender) LookupByString(key string) (datamodel.Node, error) {
	return mixins.List{TypeName: "listAmender"}.LookupByString(key)
}

func (a *listAmender) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	return mixins.List{TypeName: "listAmender"}.LookupByNode(key)
}

func (a *listAmender) LookupByIndex(idx int64) (datamodel.Node, error) {
	seg := datamodel.PathSegmentOfInt(idx)
	if mod, exists := a.mods.Get(int(idx)); exists {
		child := mod.(listElement)
		if child.elem == nil {
			bn, err := a.base.LookupByIndex(int64(child.baseIdx))
			if err != nil {
				return nil, err
			}
			child.elem = bn
			return bn, nil
		}
		return child.elem, nil
	}
	return nil, datamodel.ErrNotExists{Segment: seg}
}

func (a *listAmender) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	idx, err := seg.Index()
	if err != nil {
		return nil, datamodel.ErrInvalidSegmentForList{TroubleSegment: seg, Reason: err}
	}
	return a.LookupByIndex(idx)
}

func (a *listAmender) MapIterator() datamodel.MapIterator {
	return nil
}

func (a *listAmender) ListIterator() datamodel.ListIterator {
	modsItr := a.mods.Iterator()
	return &listAmender_Iterator{a, &modsItr, 0}
}

func (a *listAmender) Length() int64 {
	return int64(a.mods.Size())
}

func (a *listAmender) IsAbsent() bool {
	return false
}

func (a *listAmender) IsNull() bool {
	return false
}

func (a *listAmender) AsBool() (bool, error) {
	return mixins.Map{TypeName: "listAmender"}.AsBool()
}

func (a *listAmender) AsInt() (int64, error) {
	return mixins.Map{TypeName: "listAmender"}.AsInt()
}

func (a *listAmender) AsFloat() (float64, error) {
	return mixins.Map{TypeName: "listAmender"}.AsFloat()
}

func (a *listAmender) AsString() (string, error) {
	return mixins.Map{TypeName: "listAmender"}.AsString()
}

func (a *listAmender) AsBytes() ([]byte, error) {
	return mixins.Map{TypeName: "listAmender"}.AsBytes()
}

func (a *listAmender) AsLink() (datamodel.Link, error) {
	return mixins.Map{TypeName: "listAmender"}.AsLink()
}

func (a *listAmender) Prototype() datamodel.NodePrototype {
	return basicnode.Prototype.List
}

type listAmender_Iterator struct {
	amd     *listAmender
	modsItr *arraylist.Iterator
	idx     int
}

func (itr *listAmender_Iterator) Next() (idx int64, v datamodel.Node, err error) {
	if itr.Done() {
		return -1, nil, datamodel.ErrIteratorOverread{}
	}
	if itr.modsItr.Next() {
		idx = int64(itr.modsItr.Index())
		v, err = itr.amd.LookupByIndex(idx)
		if err != nil {
			return -1, nil, err
		}
		itr.idx++
		return
	}
	return -1, nil, datamodel.ErrIteratorOverread{}
}

func (itr *listAmender_Iterator) Done() bool {
	return int64(itr.idx) >= itr.amd.Length()
}

func (a *listAmender) Get(prog *Progress, path datamodel.Path, trackProgress bool) (datamodel.Node, error) {
	// If the root is requested, return the `Node` view of the amender.
	if path.Len() == 0 {
		return a.Build(), nil
	}
	// Check the budget
	if prog.Budget != nil {
		if prog.Budget.NodeBudget <= 0 {
			return nil, &ErrBudgetExceeded{BudgetKind: "node", Path: prog.Path}
		}
		prog.Budget.NodeBudget--
	}
	childSeg, remainingPath := path.Shift()
	prog.Path = prog.Path.AppendSegment(childSeg)
	childVal, err := a.LookupBySegment(childSeg)
	// Since we're explicitly looking for a node, look for the child node in the current amender state and throw an
	// error if it does not exist.
	if err != nil {
		return nil, err
	}
	childIdx, err := childSeg.Index()
	if err != nil {
		return nil, err
	}
	childAmender, err := a.storeChildAmender(childIdx, childVal, childVal.Kind(), false, trackProgress)
	if err != nil {
		return nil, err
	}
	return childAmender.Get(prog, remainingPath, trackProgress)
}

func (a *listAmender) Transform(prog *Progress, path datamodel.Path, fn TransformFn, createParents bool) (datamodel.Node, error) {
	// Allow the base node to be replaced.
	if path.Len() == 0 {
		prevNode := a.Build()
		if newNode, err := fn(*prog, prevNode); err != nil {
			return nil, err
		} else if newNode.Kind() != datamodel.Kind_List {
			return nil, fmt.Errorf("transform: cannot transform root into incompatible type: %q", newNode.Kind())
		} else {
			// Go through `newListAmender` in case `newNode` is already a list-amender.
			*a = *a.cfg.newListAmender(newNode, a.parent, a.created).(*listAmender)
			return prevNode, nil
		}
	}
	// Check the budget
	if prog.Budget != nil {
		if prog.Budget.NodeBudget <= 0 {
			return nil, &ErrBudgetExceeded{BudgetKind: "node", Path: prog.Path}
		}
		prog.Budget.NodeBudget--
	}
	childSeg, remainingPath := path.Shift()
	atLeaf := remainingPath.Len() == 0
	childIdx, err := childSeg.Index()
	var childVal datamodel.Node
	if err != nil {
		if childSeg.String() == "-" {
			// "-" indicates appending a new element to the end of the list.
			childIdx = a.Length()
			childSeg = datamodel.PathSegmentOfInt(childIdx)
		} else {
			return nil, datamodel.ErrInvalidSegmentForList{TroubleSegment: childSeg, Reason: err}
		}
	} else {
		// Don't allow the index to be equal to the length if the segment was not "-".
		if childIdx >= a.Length() {
			return nil, fmt.Errorf("transform: cannot navigate path segment %q at %q because it is beyond the list bounds", childSeg, prog.Path)
		}
		// Only lookup the segment if it was within range of the list elements. If `childIdx` is equal to the length of
		// the list, then we fall-through and append an element to the end of the list.
		childVal, err = a.LookupBySegment(childSeg)
		if err != nil {
			// - Return any error other than "not exists".
			// - If the child node does not exist and `createParents = true`, create the new hierarchy, otherwise throw
			//   an error.
			// - Even if `createParent = false`, if we're at the leaf, don't throw an error because we don't need to
			//   create any more intermediate parent nodes.
			if _, notFoundErr := err.(datamodel.ErrNotExists); !notFoundErr || !(atLeaf || createParents) {
				return nil, fmt.Errorf("transform: parent position at %q did not exist (and createParents was false)", prog.Path)
			}
		}
	}
	prog.Path = prog.Path.AppendSegment(childSeg)
	// The default behaviour will be to update the element at the specified index (if it exists). New list elements can
	// be added in two cases:
	//  - If an element is being appended to the end of the list.
	//  - If the transformation of the target node results in a list of nodes, use the first node in the list to replace
	//    the target node and then "add" the rest after. This is a bit of an ugly hack but is required for compatibility
	//    with two conflicting sets of semantics - the current `FocusedTransform`, which (quite reasonably) does an
	//    in-place replacement of list elements, and JSON Patch (https://datatracker.ietf.org/doc/html/rfc6902), which
	//    does not specify list element replacement. The only "compliant" way to do this today is to first "remove" the
	//    target node and then "add" its replacement at the same index, which seems incredibly inefficient.
	create := (childVal == nil) || atLeaf
	if atLeaf {
		if newChildVal, err := fn(*prog, childVal); err != nil {
			return nil, err
		} else if newChildVal == nil {
			a.mods.Remove(int(childIdx))
		} else if _, err = a.storeChildAmender(childIdx, newChildVal, newChildVal.Kind(), create, true); err != nil {
			return nil, err
		}
		return childVal, nil
	}
	// If we're not at the leaf yet, look ahead on the remaining path to determine what kind of intermediate parent
	// node we need to create.
	var childKind datamodel.Kind
	if childVal == nil {
		// If we're not at the leaf yet, look ahead on the remaining path to determine what kind of intermediate parent
		// node we need to create.
		nextChildSeg, _ := remainingPath.Shift()
		if _, err = nextChildSeg.Index(); err == nil {
			// As per the discussion [here](https://github.com/smrz2001/go-ipld-prime/pull/1#issuecomment-1143035685),
			// this code assumes that if we're dealing with an integral path segment, it corresponds to a list index.
			childKind = datamodel.Kind_List
		} else {
			// From the same discussion as above, any non-integral, intermediate path can be assumed to be a map key.
			childKind = datamodel.Kind_Map
		}
	} else {
		childKind = childVal.Kind()
	}
	childAmender, err := a.storeChildAmender(childIdx, childVal, childKind, create, true)
	if err != nil {
		return nil, err
	}
	return childAmender.Transform(prog, remainingPath, fn, createParents)
}

func (a *listAmender) storeChildAmender(childIdx int64, n datamodel.Node, k datamodel.Kind, create bool, trackProgress bool) (Amender, error) {
	if trackProgress {
		var childAmender Amender
		idx := int(childIdx)
		if create && (n.Kind() == datamodel.Kind_List) && (n.Length() > 0) {
			first, err := n.LookupByIndex(0)
			if err != nil {
				return nil, err
			}
			// The following logic uses a transformed list (if there is one) to perform both insertions (needed by JSON
			// Patch) and replacements (needed by `FocusedTransform`), while also providing the flexibility to insert more
			// than one element at a particular index in the list.
			//
			// Rules:
			//  - If appending to the end of the main list, all elements from the transformed list should be considered
			//    "created" because they did not exist before.
			//   - If updating at a particular index in the main list, however, use the first element from the transformed
			//     list to replace the existing element at that index in the main list, then insert the rest of the
			//     transformed list elements after.
			//
			// A special case to consider is that of a list element genuinely being a list itself. If that is the case, the
			// transformation MUST wrap the element in another list so that, once unwrapped, the element can be replaced or
			// inserted without affecting its semantics. Otherwise, the sub-list's elements will get expanded onto that
			// index in the main list.
			childAmender = a.cfg.newAmender(first, a, first.Kind(), childIdx == a.Length())
			a.mods.Set(idx, listElement{-1, childAmender.Build()})
			if n.Length() > 1 {
				elems := make([]interface{}, n.Length()-1)
				for i := range elems {
					next, err := n.LookupByIndex(int64(i) + 1)
					if err != nil {
						return nil, err
					}
					elems[i] = listElement{-1, a.cfg.newAmender(next, a, next.Kind(), true).Build()}
				}
				a.mods.Insert(idx+1, elems...)
			}
		} else {
			childAmender = a.cfg.newAmender(n, a, k, create)
			a.mods.Set(idx, listElement{-1, childAmender.Build()})
		}
		return childAmender, nil
	}
	return a.cfg.newAmender(n, a, k, create), nil
}
