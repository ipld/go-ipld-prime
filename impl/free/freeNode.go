package ipldfree

import (
	"fmt"
	"strconv"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

var (
	_ ipld.Node = &Node{}
)

/*
	Node has some internal

	This implementation of `ipld.Node` is pretty comparable to `ipldbind.Node`,
	but is somewhat simpler in implementation because values of this type can
	only be produced by its own builder patterns (and thus it requires much
	less reflection and in particular does not depend on refmt for object mapping).

	The "zero" value of this struct is interpreted as an empty map.

	This binding does not provide a serialization valid for hashing; to
	compute a CID, you'll have to convert to another kind of node.
	If you're not sure which kind serializable node to use, try `ipldcbor.Node`.
*/
type Node struct {
	kind ipld.ReprKind

	_map   map[string]ipld.Node // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_map2  map[int]ipld.Node    // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_arr   []ipld.Node          // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_bool  bool                 // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_str   string               // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_int   int                  // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_float float64              // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_bytes []byte               // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_link  cid.Cid              // Value union.  Only one of these has meaning, depending on the value of 'Type'.
}

func (n *Node) Kind() ipld.ReprKind {
	return n.kind
}

func (n *Node) IsNull() bool {
	return n.kind == ipld.ReprKind_Null
}
func (n *Node) AsBool() (v bool, _ error) {
	return n._bool, expectTyp(ipld.ReprKind_Bool, n.kind)
}
func (n *Node) AsString() (v string, _ error) {
	return n._str, expectTyp(ipld.ReprKind_String, n.kind)
}
func (n *Node) AsInt() (v int, _ error) {
	return n._int, expectTyp(ipld.ReprKind_Int, n.kind)
}
func (n *Node) AsLink() (v cid.Cid, _ error) {
	return n._link, expectTyp(ipld.ReprKind_Link, n.kind)
}

func (n *Node) Keys() ([]string, int) {
	return nil, 0 // FIXME
	// TODO need to maintain map order now, apparently, sigh
}

func (n *Node) TraverseField(pth string) (ipld.Node, error) {
	switch n.kind {
	case ipld.ReprKind_Null:
		return nil, fmt.Errorf("cannot traverse terminals")
	case ipld.ReprKind_Map:
		switch {
		case n._map != nil:
			v, _ := n._map[pth]
			return v, nil
		case n._map2 != nil:
			i, err := strconv.Atoi(pth)
			if err != nil {
				return nil, fmt.Errorf("404")
			}
			v, _ := n._map2[i]
			return v, nil
		default:
			panic("unreachable")
		}
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
