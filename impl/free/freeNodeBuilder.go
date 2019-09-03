package ipldfree

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
)

// NodeBuilder returns a new ipld.NodeBuilder implementation that will produce
// ipldfree.Node instances.
//
// There are no constraints on free nodes, so none of the create methods
// will ever return errors.
func NodeBuilder() ipld.NodeBuilder {
	return nodeBuilder{}
}

type nodeBuilder struct {
	predecessor *Node // optional; only relevant for "Amend*" methods.
}

func (nb nodeBuilder) CreateMap() (ipld.MapBuilder, error) {
	return &mapBuilder{n: Node{kind: ipld.ReprKind_Map, _map: make(map[string]ipld.Node)}}, nil
}
func (nb nodeBuilder) AmendMap() (ipld.MapBuilder, error) {
	if nb.predecessor == nil {
		return nb.CreateMap()
	}
	if nb.predecessor.kind != ipld.ReprKind_List {
		return nil, fmt.Errorf("AmendMap cannot be used when predecessor was not a map")
	}
	newMap := make(map[string]ipld.Node, len(nb.predecessor._map))
	for k, v := range nb.predecessor._map {
		newMap[k] = v
	}
	newArr := make([]string, len(nb.predecessor._mapOrd))
	copy(newArr, nb.predecessor._mapOrd)
	return &mapBuilder{n: Node{kind: ipld.ReprKind_Map, _map: newMap, _mapOrd: newArr}}, nil
}
func (nb nodeBuilder) CreateList() (ipld.ListBuilder, error) {
	return &listBuilder{n: Node{kind: ipld.ReprKind_List}}, nil
}
func (nb nodeBuilder) AmendList() (ipld.ListBuilder, error) {
	if nb.predecessor == nil {
		return nb.CreateList()
	}
	if nb.predecessor.kind != ipld.ReprKind_List {
		return nil, fmt.Errorf("AmendList cannot be used when predecessor was not a list")
	}
	newArr := make([]ipld.Node, len(nb.predecessor._arr))
	copy(newArr, nb.predecessor._arr)
	return &listBuilder{n: Node{kind: ipld.ReprKind_List, _arr: newArr}}, nil
}
func (nb nodeBuilder) CreateNull() (ipld.Node, error) {
	return &Node{kind: ipld.ReprKind_Null}, nil
}
func (nb nodeBuilder) CreateBool(v bool) (ipld.Node, error) {
	return &Node{kind: ipld.ReprKind_Bool, _bool: v}, nil
}
func (nb nodeBuilder) CreateInt(v int) (ipld.Node, error) {
	return &Node{kind: ipld.ReprKind_Int, _int: v}, nil
}
func (nb nodeBuilder) CreateFloat(v float64) (ipld.Node, error) {
	return &Node{kind: ipld.ReprKind_Float, _float: v}, nil
}
func (nb nodeBuilder) CreateString(v string) (ipld.Node, error) {
	return &Node{kind: ipld.ReprKind_String, _str: v}, nil
}
func (nb nodeBuilder) CreateBytes(v []byte) (ipld.Node, error) {
	return &Node{kind: ipld.ReprKind_Bytes, _bytes: v}, nil
}
func (nb nodeBuilder) CreateLink(v ipld.Link) (ipld.Node, error) {
	return &Node{kind: ipld.ReprKind_Link, _link: v}, nil
}

type mapBuilder struct {
	n Node // a wip node; initialized at construction.
	// whole builder object nil'd after terminal `Build()` call to prevent reuse.
}

// Insert adds a k:v pair to the map.
//
// As is usual for maps, the key must have kind==ReprKind_String.
//
// Keys not already present in the map will be appened to the end of the
// iteration order; keys already present retain their original order.
func (mb *mapBuilder) Insert(k, v ipld.Node) error {
	ks, err := k.AsString()
	if err != nil {
		return fmt.Errorf("invalid node for map key: %s", err)
	}
	_, exists := mb.n._map[ks]
	if exists {
		return fmt.Errorf("repeated map key: %s", ks)
	}
	mb.n._map[ks] = v
	mb.n._mapOrd = append(mb.n._mapOrd, ks)
	return nil
}
func (mb *mapBuilder) Delete(k ipld.Node) error {
	panic("NYI") // and see the "review: MapBuilder.Delete" comment in the interface defn file.
}
func (mb *mapBuilder) Build() (ipld.Node, error) {
	v := mb.n
	mb = nil
	return &v, nil
}
func (mapBuilder) BuilderForKeys() ipld.NodeBuilder {
	return justStringNodeBuilder{}
}
func (mapBuilder) BuilderForValue(_ string) ipld.NodeBuilder {
	return nodeBuilder{}
}

type listBuilder struct {
	n Node // a wip node; initialized at construction.
	// whole builder object nil'd after terminal `Build()` call to prevent reuse.
}

func (lb *listBuilder) AppendAll(vs []ipld.Node) error {
	off := len(lb.n._arr)
	new := off + len(vs)
	growList(&lb.n._arr, new-1)
	copy(lb.n._arr[off:new], vs)
	return nil
}
func (lb *listBuilder) Append(v ipld.Node) error {
	lb.n._arr = append(lb.n._arr, v)
	return nil
}
func (lb *listBuilder) Set(idx int, v ipld.Node) error {
	growList(&lb.n._arr, idx)
	lb.n._arr[idx] = v
	return nil
}
func (lb *listBuilder) Build() (ipld.Node, error) {
	v := lb.n
	lb = nil
	return &v, nil
}
func (listBuilder) BuilderForValue(_ int) ipld.NodeBuilder {
	return nodeBuilder{}
}

func growList(l *[]ipld.Node, k int) {
	oldLen := len(*l)
	minLen := k + 1
	if minLen > oldLen {
		// Grow.
		oldCap := cap(*l)
		if minLen > oldCap {
			// Out of cap; do whole new backing array allocation.
			//  Growth maths are per stdlib's reflect.grow.
			// First figure out how much growth to do.
			newCap := oldCap
			if newCap == 0 {
				newCap = minLen
			} else {
				for minLen > newCap {
					if minLen < 1024 {
						newCap += newCap
					} else {
						newCap += newCap / 4
					}
				}
			}
			// Now alloc and copy over old.
			newArr := make([]ipld.Node, minLen, newCap)
			copy(newArr, *l)
			*l = newArr
		} else {
			// Still have cap, just extend the slice.
			*l = (*l)[0:minLen]
		}
	}
}
