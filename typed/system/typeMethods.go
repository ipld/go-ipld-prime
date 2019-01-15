package typesystem

import (
	"github.com/ipld/go-ipld-prime"
)

func (TypeBool) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Bool
}
func (TypeString) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (TypeBytes) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Bytes
}
func (TypeInt) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Int
}
func (TypeFloat) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Float
}
func (TypeMap) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (TypeList) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_List
}
func (TypeLink) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Link
}
func (tv TypeUnion) ReprKind() ipld.ReprKind {
	// REVIEW: this may fib; has the bizarre property of being dependent on the *concrete value* for kinded unions!
	if tv.Style == UnionStyle_Kinded {
		return ipld.ReprKind_Invalid
	} else {
		return ipld.ReprKind_Map
	}
}
func (tv TypeObject) ReprKind() ipld.ReprKind {
	if tv.TupleStyle {
		return ipld.ReprKind_List
	} else {
		return ipld.ReprKind_Map
	}
}
func (TypeEnum) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
