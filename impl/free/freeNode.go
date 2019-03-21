package ipldfree

import (
	"fmt"
	"strconv"

	"github.com/ipld/go-ipld-prime"
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

func (n *Node) IsNull() bool {
	return n.kind == ipld.ReprKind_Null
}
func (n *Node) AsBool() (v bool, _ error) {
	return n._bool, expectTyp(ipld.ReprKind_Bool, n.kind)
}
func (n *Node) AsInt() (v int, _ error) {
	return n._int, expectTyp(ipld.ReprKind_Int, n.kind)
}
func (n *Node) AsFloat() (v float64, _ error) {
	return n._float, expectTyp(ipld.ReprKind_Float, n.kind)
}
func (n *Node) AsString() (v string, _ error) {
	return n._str, expectTyp(ipld.ReprKind_String, n.kind)
}
func (n *Node) AsBytes() (v []byte, _ error) {
	return n._bytes, expectTyp(ipld.ReprKind_Bytes, n.kind)
}
func (n *Node) AsLink() (v ipld.Link, _ error) {
	return n._link, expectTyp(ipld.ReprKind_Link, n.kind)
}

func (n *Node) NodeBuilder() ipld.NodeBuilder {
	return nodeBuilder{n}
}

func (n *Node) MapIterator() ipld.MapIterator {
	return &mapIterator{n, 0, expectTyp(ipld.ReprKind_Map, n.kind)}
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
	return &listIterator{n, 0, expectTyp(ipld.ReprKind_List, n.kind)}
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

func (n *Node) TraverseField(pth string) (ipld.Node, error) {
	switch n.kind {
	case ipld.ReprKind_Invalid:
		return nil, fmt.Errorf("cannot traverse a node that is undefined")
	case ipld.ReprKind_Null:
		return nil, fmt.Errorf("cannot traverse terminals")
	case ipld.ReprKind_Map:
		v, exists := n._map[pth]
		if !exists {
			return nil, fmt.Errorf("404")
		}
		return v, nil
	case ipld.ReprKind_List:
		i, err := strconv.Atoi(pth)
		if err != nil {
			return nil, fmt.Errorf("404")
		}
		if i >= len(n._arr) {
			return nil, fmt.Errorf("404")
		}
		return n._arr[i], nil
	case ipld.ReprKind_Bool,
		ipld.ReprKind_String,
		ipld.ReprKind_Bytes,
		ipld.ReprKind_Int,
		ipld.ReprKind_Float,
		ipld.ReprKind_Link:
		return nil, fmt.Errorf("cannot traverse terminals")
	default:
		panic("unreachable")
	}
}

func (n *Node) TraverseIndex(idx int) (ipld.Node, error) {
	switch n.kind {
	case ipld.ReprKind_Invalid:
		return nil, fmt.Errorf("cannot traverse a node that is undefined")
	case ipld.ReprKind_Null:
		return nil, fmt.Errorf("cannot traverse terminals")
	case ipld.ReprKind_Map:
		return nil, fmt.Errorf("cannot traverse map by numeric index")
		// REVIEW: there's an argument that maybe we should support this; would be '_map2' code.
	case ipld.ReprKind_List:
		if idx >= len(n._arr) {
			return nil, fmt.Errorf("404")
		}
		if n._arr[idx] == nil {
			return nil, fmt.Errorf("404")
		}
		return n._arr[idx], nil
	case ipld.ReprKind_Bool,
		ipld.ReprKind_String,
		ipld.ReprKind_Bytes,
		ipld.ReprKind_Int,
		ipld.ReprKind_Float,
		ipld.ReprKind_Link:
		return nil, fmt.Errorf("cannot traverse terminals")
	default:
		panic("unreachable")
	}
}

func expectTyp(expect, actual ipld.ReprKind) error {
	if expect == actual {
		return nil
	}
	return fmt.Errorf("type assertion rejected: node is %q, assertion was for %q", actual, expect)
}
