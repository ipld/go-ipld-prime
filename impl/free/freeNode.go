package ipldfree

import (
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
	typ typ

	_map   map[string]interface{} // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_map2  map[int]interface{}    // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_arr   []interface{}          // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_str   string                 // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_bytes []byte                 // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_bool  bool                   // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_int   int64                  // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_uint  uint64                 // Value union.  Only one of these has meaning, depending on the value of 'Type'.
	_float float64                // Value union.  Only one of these has meaning, depending on the value of 'Type'.

}

type typ struct{ t byte }

var (
	tZero  = typ{} // treat as tMap
	tMap   = typ{'{'}
	tArr   = typ{'['}
	tStr   = typ{'s'}
	tBytes = typ{'x'}
	tBool  = typ{'b'}
	tInt   = typ{'i'}
	tUint  = typ{'u'}
	tFloat = typ{'f'}
)

func (n *Node) GetField(pth []string) (v interface{}, _ error) {
	return v, traverse(n.bound, pth, n.atlas, reflect.ValueOf(v))
}
