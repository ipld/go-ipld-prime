package ipldfree

import (
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

func (n *Node) SetField(k string, v ipld.Node) {
	n.coerceType(tMap)
	n._map[k] = v
}
func (n *Node) SetIndex(k int, v ipld.Node) {
	n.coerceType(tArr)
	// REVIEW: there are implications to serial arrays as we spec'd them.
	//  Namely, they can't be sparse.  It's just not defined.
	//  And that means we simply have to have a way to define the length.
	//  We can do implicit grows via this setter; I'm fine with that.
	//  But we'll also need, evidently, a Truncate method.
	//  (Or, a magical sentinel value for node that says EOL.)
	oldLen := len(n._arr)
	minLen := k + 1
	if minLen > oldLen {
		// Grow.
		oldCap := cap(n._arr)
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
			copy(newArr, n._arr)
			n._arr = newArr
		} else {
			// Still have cap, just extend the slice.
			n._arr = n._arr[0:minLen]
		}
	}
	n._arr[k] = v
	//fmt.Printf("len,cap is now %d,%d\n", len(n._arr), cap(n._arr))
}
func (n *Node) SetBool(v bool) {
	n.coerceType(tBool)
	n._bool = v
}
func (n *Node) SetString(v string) {
	n.coerceType(tStr)
	n._str = v
}
func (n *Node) SetInt(v int) {
	n.coerceType(tInt)
	n._int = v
}
func (n *Node) SetFloat(v float64) {
	n.coerceType(tFloat)
	n._float = v
}
func (n *Node) SetBytes(v []byte) {
	n.coerceType(tBytes)
	n._bytes = v
}
func (n *Node) SetLink(v cid.Cid) {
	n.coerceType(tLink)
	n._link = v
}

func (n *Node) coerceType(newKind typ) {
	// If this node pointer has actually just been nil, initialize.
	//  (Our arrays sometimes initialize full of nils, so this comes up.)
	// TODO
	// REVIEW actually it's pretty dubious that we should return those.
	//  Nobody ever said our concept of array should be non-sparse and get nils in it.

	// Clear previous data, if relevant.
	//  Don't bother with zeroing finite-size scalars.
	switch n.typ {
	case tMap:
		switch newKind {
		case tMap:
			return
		default:
			n._map = nil
		}
	case tArr:
		switch newKind {
		case tArr:
			return
		default:
			n._arr = nil
		}
	case tStr:
		switch newKind {
		case tStr:
			return
		default:
			n._str = ""
		}
	case tBytes:
		switch newKind {
		case tBytes:
			return
		default:
			n._bytes = nil
		}
	case tLink:
		switch newKind {
		case tLink:
			return
		default:
			n._link = cid.Undef
		}
	}
	// Set new type union marker.
	//  Initialize empty value if necessary (maps).
	n.typ = newKind
	switch newKind {
	case tMap:
		n._map = make(map[string]ipld.Node)
	}
	// You'll still want to set the value itself after this.
}
