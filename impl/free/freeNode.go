package ipldfree

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
)

var (
	_ ipld.Node = &Node{}
)

/*
	Node is an implementatin of `ipld.Node` that can contain any content.

	This implementation is extremely simple; it is general-purpose,
	but not optimized for any particular purpose.

	The "zero" value of this struct has a kind of ReprKind_Invalid.
	NodeBuilder must be used to produce valid instances of Node.
*/
type Node struct {
	kind ipld.ReprKind

	_map    map[string]ipld.Node // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_mapOrd []string             // Conjugate to _map, only has meaning depending on the value of 'Type'.
	_arr    []ipld.Node          // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_bool   bool                 // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_int    int                  // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_float  float64              // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_str    string               // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_bytes  []byte               // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_link   ipld.Link            // Value union.  Only one of these has meaning, depending on the value of 'Type'.
}

func (n *Node) ReprKind() ipld.ReprKind {
	return n.kind
}

func (Node) IsUndefined() bool {
	return false
}
func (n *Node) IsNull() bool {
	return n.kind == ipld.ReprKind_Null
}
func (n *Node) AsBool() (v bool, _ error) {
	if n.kind != ipld.ReprKind_Bool {
		return false, ipld.ErrWrongKind{MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: n.kind}
	}
	return n._bool, nil
}
func (n *Node) AsInt() (v int, _ error) {
	if n.kind != ipld.ReprKind_Int {
		return 0, ipld.ErrWrongKind{MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: n.kind}
	}
	return n._int, nil
}
func (n *Node) AsFloat() (v float64, _ error) {
	if n.kind != ipld.ReprKind_Float {
		return 0, ipld.ErrWrongKind{MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: n.kind}
	}
	return n._float, nil
}
func (n *Node) AsString() (v string, _ error) {
	if n.kind != ipld.ReprKind_String {
		return "", ipld.ErrWrongKind{MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: n.kind}
	}
	return n._str, nil
}
func (n *Node) AsBytes() (v []byte, _ error) {
	if n.kind != ipld.ReprKind_Bytes {
		return nil, ipld.ErrWrongKind{MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: n.kind}
	}
	return n._bytes, nil
}
func (n *Node) AsLink() (v ipld.Link, _ error) {
	if n.kind != ipld.ReprKind_Link {
		return nil, ipld.ErrWrongKind{MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: n.kind}
	}
	return n._link, nil
}

func (n *Node) NodeBuilder() ipld.NodeBuilder {
	return nodeBuilder{n}
}

func (n *Node) MapIterator() ipld.MapIterator {
	if n.kind != ipld.ReprKind_Map {
		return &mapIterator{n, 0, ipld.ErrWrongKind{MethodName: "MapIterator", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: n.kind}}
	}
	return &mapIterator{n, 0, nil}
}

type mapIterator struct {
	node *Node
	idx  int
	err  error
}

func (itr *mapIterator) Next() (ipld.Node, ipld.Node, error) {
	if itr.err != nil {
		return nil, nil, itr.err
	}
	k := itr.node._mapOrd[itr.idx]
	v := itr.node._map[k]
	itr.idx++
	return &Node{kind: ipld.ReprKind_String, _str: k}, v, nil
}
func (itr *mapIterator) Done() bool {
	if itr.err != nil {
		return false
	}
	return itr.idx >= len(itr.node._mapOrd)
}

func (n *Node) ListIterator() ipld.ListIterator {
	if n.kind != ipld.ReprKind_List {
		return &listIterator{n, 0, ipld.ErrWrongKind{MethodName: "ListIterator", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: n.kind}}
	}
	return &listIterator{n, 0, nil}
}

type listIterator struct {
	node *Node
	idx  int
	err  error
}

func (itr *listIterator) Next() (int, ipld.Node, error) {
	if itr.err != nil {
		return -1, nil, itr.err
	}
	v := itr.node._arr[itr.idx]
	idx := itr.idx
	itr.idx++
	return idx, v, nil
}
func (itr *listIterator) Done() bool {
	if itr.err != nil {
		return false
	}
	return itr.idx >= len(itr.node._arr)
}

func (n *Node) Length() int {
	switch n.ReprKind() {
	case ipld.ReprKind_Map:
		return len(n._mapOrd)
	case ipld.ReprKind_List:
		return len(n._arr)
	default:
		return -1
	}
}

func (n *Node) LookupString(pth string) (ipld.Node, error) {
	switch n.kind {
	case ipld.ReprKind_Map:
		v, exists := n._map[pth]
		if !exists {
			return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(pth)}
		}
		return v, nil
	case ipld.ReprKind_Invalid,
		ipld.ReprKind_Null,
		ipld.ReprKind_List,
		ipld.ReprKind_String,
		ipld.ReprKind_Bytes,
		ipld.ReprKind_Int,
		ipld.ReprKind_Float,
		ipld.ReprKind_Link:
		return nil, ipld.ErrWrongKind{MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: n.kind}
	default:
		panic("unreachable")
	}
}

func (n *Node) Lookup(key ipld.Node) (ipld.Node, error) {
	switch n.kind {
	case ipld.ReprKind_Map:
		ks, err := key.AsString()
		if err != nil {
			return nil, ipld.ErrInvalidKey{fmt.Sprintf("got %s, need string", key.ReprKind())}
		}
		v, exists := n._map[ks]
		if !exists {
			return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(ks)}
		}
		return v, nil
	default:
		return nil, ipld.ErrWrongKind{MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: n.kind}
	}
}

func (n *Node) LookupIndex(idx int) (ipld.Node, error) {
	switch n.kind {
	case ipld.ReprKind_List:
		if idx >= len(n._arr) {
			return nil, ipld.ErrNotExists{ipld.PathSegmentOfInt(idx)}
		}
		if n._arr[idx] == nil {
			return nil, ipld.ErrNotExists{ipld.PathSegmentOfInt(idx)}
		}
		return n._arr[idx], nil
	case ipld.ReprKind_Invalid,
		ipld.ReprKind_Null,
		ipld.ReprKind_Map,
		ipld.ReprKind_Bool,
		ipld.ReprKind_String,
		ipld.ReprKind_Bytes,
		ipld.ReprKind_Int,
		ipld.ReprKind_Float,
		ipld.ReprKind_Link:
		return nil, ipld.ErrWrongKind{MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: n.kind}
	default:
		panic("unreachable")
	}
}

func (n *Node) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	switch n.kind {
	case ipld.ReprKind_Map:
		return n.LookupString(seg.String())
	case ipld.ReprKind_List:
		idx, err := seg.Index()
		if err != nil {
			return nil, err
		}
		return n.LookupIndex(idx)
	case ipld.ReprKind_Invalid,
		ipld.ReprKind_Null,
		ipld.ReprKind_Bool,
		ipld.ReprKind_String,
		ipld.ReprKind_Bytes,
		ipld.ReprKind_Int,
		ipld.ReprKind_Float,
		ipld.ReprKind_Link:
		return nil, ipld.ErrWrongKind{MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: n.kind}
	default:
		panic("unreachable")
	}
}
