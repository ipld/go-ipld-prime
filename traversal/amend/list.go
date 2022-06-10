package amend

import (
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
	idx  int
	elem datamodel.Node
}

type listAmender struct {
	base    datamodel.Node
	parent  Amender
	created bool
	mods    *arraylist.List
}

func newListAmender(base datamodel.Node, parent Amender, create bool) Amender {
	var mods *arraylist.List
	// If the base node is already a list-amender *for the same base node*, reuse the modification metadata because that
	// encapsulates all accumulated modifications.
	if amd, castOk := base.(*listAmender); castOk && (base == amd.base) {
		mods = amd.mods
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
		mods = arraylist.New(elems...)
	}
	return &listAmender{base, parent, create, mods}
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
			baseNode, err := a.base.LookupByIndex(int64(child.idx))
			if err != nil {
				return nil, err
			}
			child.elem = baseNode
			return baseNode, nil
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

func (a *listAmender) Get(path datamodel.Path) (datamodel.Node, error) {
	childSeg, remainingPath := path.Shift()
	childVal, err := a.LookupBySegment(childSeg)
	atLeaf := remainingPath.Len() == 0
	// Since we're explicitly looking for a node, look for the child node in the current amender state and throw an
	// error if it does not exist.
	if err != nil {
		return nil, err
	}
	childIdx, err := childSeg.Index()
	if err != nil {
		return nil, err
	}
	childAmender := newAmender(childVal, a, childVal.Kind(), false)
	a.mods.Set(int(childIdx), listElement{int(childIdx), childAmender.(datamodel.Node)})
	if atLeaf {
		return childVal, nil
	} else {
		return childAmender.Get(remainingPath)
	}
}

func (a *listAmender) Add(path datamodel.Path, value datamodel.Node, createParents bool) error {
	childSeg, remainingPath := path.Shift()
	atLeaf := remainingPath.Len() == 0
	childIdx, err := childSeg.Index()
	if err != nil {
		return datamodel.ErrInvalidSegmentForList{TroubleSegment: childSeg, Reason: err}
	}
	// Allow the index to be equal to the length - this just means that a new element needs to be added to the end of
	// the list (i.e. appended).
	if childIdx > a.Length() {
		return datamodel.ErrNotExists{Segment: childSeg}
	}
	childVal, err := a.LookupBySegment(childSeg)
	if err != nil {
		// - Return any error other than "not exists".
		// - If the chile node does not exist and `createParents = true`, create the new hierarchy, otherwise throw an
		//   error.
		// - Even if `createParent = false`, if we're at the leaf, don't throw an error because we don't need to create
		//   any more intermediate parent nodes.
		if _, notFoundErr := err.(datamodel.ErrNotExists); !notFoundErr || !(atLeaf || createParents) {
			return err
		}
	}
	// While building the nested amender tree, only count nodes as "added" when they didn't exist and had to be created
	// to fill out the hierarchy.
	// In the case of a list, also consider a node "added" if we're at the leaf. Even if there already was a child at
	// that index, it just means we need to "insert" a new node at the index.
	create := false
	if (childVal == nil) || atLeaf {
		create = true
	}
	var childKind datamodel.Kind
	if atLeaf {
		childVal = value
		childKind = value.Kind()
	} else {
		// If we're not at the leaf yet, look ahead on the remaining path to determine what kind of intermediate parent
		// node we need to create.
		nextChildSeg, _ := remainingPath.Shift()
		if _, err := nextChildSeg.Index(); err == nil {
			// As per the discussion [here](https://github.com/smrz2001/go-ipld-prime/pull/1#issuecomment-1143035685),
			// this code assumes that if we're dealing with an integral path segment, it corresponds to a list index.
			childKind = datamodel.Kind_List
		} else {
			// From the same discussion as above, any non-integral, intermediate path can be assumed to be a map key.
			childKind = datamodel.Kind_Map
		}
	}
	// When adding to a list-amender we're *always* creating a new node, never "wrapping" an existing one. This is by
	// virtue of list semantics, where an addition means inserting a new element, even if one already existed at the
	// specified index.
	childAmender := newAmender(childVal, a, childKind, create)
	if create {
		a.mods.Insert(int(childIdx), listElement{int(childIdx), childAmender.(datamodel.Node)})
	} else {
		a.mods.Set(int(childIdx), listElement{int(childIdx), childAmender.(datamodel.Node)})
	}
	if atLeaf {
		return nil
	} else {
		return childAmender.Add(remainingPath, value, createParents)
	}
}

func (a *listAmender) Remove(path datamodel.Path) (datamodel.Node, error) {
	childSeg, remainingPath := path.Shift()
	childVal, err := a.LookupBySegment(childSeg)
	atLeaf := remainingPath.Len() == 0
	// Since we're explicitly looking for a node, look for the child node in the current amender state and throw an
	// error if it does not exist.
	if err != nil {
		return nil, err
	}
	childIdx, err := childSeg.Index()
	if err != nil {
		return nil, err
	}
	if atLeaf {
		a.mods.Remove(int(childIdx))
		return childVal, nil
	} else {
		childAmender := newAmender(childVal, a, childVal.Kind(), false)
		a.mods.Set(int(childIdx), listElement{int(childIdx), childAmender.(datamodel.Node)})
		return childAmender.Remove(remainingPath)
	}
}

func (a *listAmender) Replace(path datamodel.Path, value datamodel.Node) (datamodel.Node, error) {
	childSeg, remainingPath := path.Shift()
	childVal, err := a.LookupBySegment(childSeg)
	atLeaf := remainingPath.Len() == 0
	// Since we're explicitly looking for a node, look for the child node in the current amender state and throw an
	// error if it does not exist.
	if err != nil {
		return nil, err
	}
	childIdx, err := childSeg.Index()
	if err != nil {
		return nil, err
	}
	childAmender := newAmender(childVal, a, childVal.Kind(), false)
	a.mods.Set(int(childIdx), listElement{int(childIdx), childAmender.(datamodel.Node)})
	if atLeaf {
		return childVal, nil
	} else {
		return childAmender.Replace(remainingPath, value)
	}
}
