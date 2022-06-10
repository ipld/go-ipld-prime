package amend

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

var (
	_ datamodel.Node = &mapAmender{}
	_ Amender        = &mapAmender{}
)

type mapAmender struct {
	base    datamodel.Node
	parent  Amender
	created bool
	// This is the information needed to present an accurate "effective" view of the base node and all accumulated
	// modifications.
	mods *linkedhashmap.Map
	// This is the count of children *present in the base node* that are removed. Knowing this count allows accurate
	// traversal of the "effective" node view.
	rems int
	// This is the count of new children. If an added node is removed, this count should be decremented instead of
	// `rems`.
	adds int
}

func newMapAmender(base datamodel.Node, parent Amender, create bool) Amender {
	// If the base node is already a map-amender *for the same base node*, reuse the modification metadata but reset
	// other information (viz. parent, created).
	if amd, castOk := base.(*mapAmender); castOk && (base == amd.base) {
		return &mapAmender{base, parent, create, amd.mods, amd.rems, amd.adds}
	} else {
		// Start with fresh state because existing metadata could not be reused.
		return &mapAmender{base, parent, create, linkedhashmap.New(), 0, 0}
	}
}

func (a *mapAmender) Build() datamodel.Node {
	// `mapAmender` is also a `Node`.
	return (datamodel.Node)(a)
}

func (a *mapAmender) Kind() datamodel.Kind {
	return datamodel.Kind_Map
}

func (a *mapAmender) LookupByString(key string) (datamodel.Node, error) {
	seg := datamodel.PathSegmentOfString(key)
	// Added/removed nodes override the contents of the base node
	if mod, exists := a.mods.Get(seg); exists {
		v := mod.(datamodel.Node)
		if v.IsNull() {
			return nil, datamodel.ErrNotExists{Segment: seg}
		}
		return v, nil
	}
	// Fallback to base node
	if a.base != nil {
		return a.base.LookupByString(key)
	}
	return nil, datamodel.ErrNotExists{Segment: seg}
}

func (a *mapAmender) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return a.LookupByString(ks)
}

func (a *mapAmender) LookupByIndex(idx int64) (datamodel.Node, error) {
	return mixins.Map{TypeName: "mapAmender"}.LookupByIndex(idx)
}

func (a *mapAmender) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	return a.LookupByString(seg.String())
}

func (a *mapAmender) MapIterator() datamodel.MapIterator {
	var baseItr datamodel.MapIterator = nil
	// If all children were removed from the base node, or no base node was specified, there is nothing to iterate
	// over w.r.t. that node.
	if (a.base != nil) && (int64(a.rems) < a.base.Length()) {
		baseItr = a.base.MapIterator()
	}
	var modsItr *linkedhashmap.Iterator
	if (a.rems != 0) || (a.adds != 0) {
		itr := a.mods.Iterator()
		modsItr = &itr
	}
	return &mapAmender_Iterator{a, modsItr, baseItr, 0}
}

func (a *mapAmender) ListIterator() datamodel.ListIterator {
	return nil
}

func (a *mapAmender) Length() int64 {
	length := int64(a.adds - a.rems)
	if a.base != nil {
		length = length + a.base.Length()
	}
	return length
}

func (a *mapAmender) IsAbsent() bool {
	return false
}

func (a *mapAmender) IsNull() bool {
	return false
}

func (a *mapAmender) AsBool() (bool, error) {
	return mixins.Map{TypeName: "mapAmender"}.AsBool()
}

func (a *mapAmender) AsInt() (int64, error) {
	return mixins.Map{TypeName: "mapAmender"}.AsInt()
}

func (a *mapAmender) AsFloat() (float64, error) {
	return mixins.Map{TypeName: "mapAmender"}.AsFloat()
}

func (a *mapAmender) AsString() (string, error) {
	return mixins.Map{TypeName: "mapAmender"}.AsString()
}

func (a *mapAmender) AsBytes() ([]byte, error) {
	return mixins.Map{TypeName: "mapAmender"}.AsBytes()
}

func (a *mapAmender) AsLink() (datamodel.Link, error) {
	return mixins.Map{TypeName: "mapAmender"}.AsLink()
}

func (a *mapAmender) Prototype() datamodel.NodePrototype {
	return basicnode.Prototype.Map
}

type mapAmender_Iterator struct {
	amd     *mapAmender
	modsItr *linkedhashmap.Iterator
	baseItr datamodel.MapIterator
	idx     int
}

func (itr *mapAmender_Iterator) Next() (k datamodel.Node, v datamodel.Node, _ error) {
	if itr.Done() {
		return nil, nil, datamodel.ErrIteratorOverread{}
	}
	if itr.baseItr != nil {
		// Iterate over base node first to maintain ordering.
		var err error
		for !itr.baseItr.Done() {
			k, v, err = itr.baseItr.Next()
			if err != nil {
				return nil, nil, err
			}
			ks, _ := k.AsString()
			if err != nil {
				return nil, nil, err
			}
			if mod, exists := itr.amd.mods.Get(datamodel.PathSegmentOfString(ks)); exists {
				v = mod.(datamodel.Node)
				// Skip removed nodes
				if v.IsNull() {
					continue
				}
				// Fall-through and return wrapped nodes
			}
			// We found a "real" node to return, increment the counter.
			itr.idx++
			return
		}
	}
	if itr.modsItr != nil {
		// Iterate over mods, skipping removed nodes.
		for itr.modsItr.Next() {
			key := itr.modsItr.Key()
			k = basicnode.NewString(key.(datamodel.PathSegment).String())
			v = itr.modsItr.Value().(datamodel.Node)
			// Skip removed nodes.
			if v.IsNull() {
				continue
			}
			// Skip "wrapper" nodes that represent existing sub-nodes in the hierarchy corresponding to an added leaf
			// node.
			if amd, castOk := v.(Amender); castOk && !isCreated(amd) {
				continue
			}
			// We found a "real" node to return, increment the counter.
			itr.idx++
			return
		}
	}
	return nil, nil, datamodel.ErrIteratorOverread{}
}

func (itr *mapAmender_Iterator) Done() bool {
	// Iteration is complete when all source nodes have been processed (skipping removed nodes) and all mods have been
	// processed.
	return int64(itr.idx) >= itr.amd.Length()
}

func (a *mapAmender) Get(path datamodel.Path) (datamodel.Node, error) {
	childSeg, remainingPath := path.Shift()
	childVal, err := a.LookupBySegment(childSeg)
	atLeaf := remainingPath.Len() == 0
	// Since we're explicitly looking for a node, look for the child node in the current amender state and throw an
	// error if it does not exist.
	if err != nil {
		return nil, err
	}
	childAmender := newAmender(childVal, a, childVal.Kind(), false)
	a.mods.Put(childSeg, childAmender)
	if atLeaf {
		return childVal, nil
	} else {
		return childAmender.Get(remainingPath)
	}
}

func (a *mapAmender) Add(path datamodel.Path, value datamodel.Node, createParents bool) error {
	childSeg, remainingPath := path.Shift()
	atLeaf := remainingPath.Len() == 0
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
	create := false
	if childVal == nil {
		a.adds++
		create = true
	}
	var childKind datamodel.Kind
	if atLeaf {
		if childVal != nil {
			// The leaf must not already exist.
			return datamodel.ErrRepeatedMapKey{Key: basicnode.NewString(childSeg.String())}
		}
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
	childAmender := newAmender(childVal, a, childKind, create)
	a.mods.Put(childSeg, childAmender)
	if atLeaf {
		return nil
	} else {
		return childAmender.Add(remainingPath, value, createParents)
	}
}

func (a *mapAmender) Remove(path datamodel.Path) (datamodel.Node, error) {
	childSeg, remainingPath := path.Shift()
	childVal, err := a.LookupBySegment(childSeg)
	atLeaf := remainingPath.Len() == 0
	// Since we're explicitly looking for a node, look for the child node in the current amender state and throw an
	// error if it does not exist.
	if err != nil {
		return nil, err
	}
	if atLeaf {
		// Use the "Null" node to indicate a removed child.
		a.mods.Put(childSeg, datamodel.Null)
		// If this parent node is an amender and present in the base hierarchy, increment `rems`, otherwise decrement
		// `adds`. This allows us to retain knowledge about the "history" of the base hierarchy.
		if ma, mapCastOk := childVal.(*mapAmender); mapCastOk {
			if ma.base != nil {
				a.rems++
			} else {
				a.adds--
			}
		} else if la, listCastOk := childVal.(*listAmender); listCastOk {
			if la.base != nil {
				a.rems++
			} else {
				a.adds--
			}
		} else {
			a.rems++
		}
		return childVal, nil
	} else {
		childAmender := newAmender(childVal, a, childVal.Kind(), false)
		// No need to update `rems` since we haven't reached the parent whose child is being removed.
		a.mods.Put(childSeg, childAmender)
		return childAmender.Remove(remainingPath)
	}
}

func (a *mapAmender) Replace(path datamodel.Path, value datamodel.Node) (datamodel.Node, error) {
	childSeg, remainingPath := path.Shift()
	childVal, err := a.LookupBySegment(childSeg)
	atLeaf := remainingPath.Len() == 0
	// Since we're explicitly looking for a node, look for the child node in the current amender state and throw an
	// error if it does not exist.
	if err != nil {
		return nil, err
	}
	var childKind datamodel.Kind
	if atLeaf {
		childVal = value
		childKind = value.Kind()
	} else if _, err := childSeg.Index(); err == nil {
		// As per the discussion [here](https://github.com/smrz2001/go-ipld-prime/pull/1#issuecomment-1143035685), this
		// code assumes that if we're dealing with an integral path segment, it corresponds to a list index.
		childKind = datamodel.Kind_List
	} else {
		// From the same discussion as above, any non-integral, intermediate path can be assumed to be a map key.
		childKind = datamodel.Kind_Map
	}
	childAmender := newAmender(childVal, a, childKind, false)
	a.mods.Put(childSeg, childAmender)
	if atLeaf {
		return childVal, nil
	} else {
		return childAmender.Replace(remainingPath, value)
	}
}
