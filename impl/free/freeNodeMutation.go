package ipldfree

import (
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

func (n *Node) SetField(k string, v ipld.Node) {
	n.coerceType(tMap)
	//TODO
}
func (n *Node) SetIndex(k int, v ipld.Node) {
	n.coerceType(tArr)
	//TODO
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
		n._map = make(map[string]*Node)
	}
	// You'll still want to set the value itself after this.
}
