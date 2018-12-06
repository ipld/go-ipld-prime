package typed

func (TypeBool) ReprKind() ReprKind {
	return ReprKind_Bool
}
func (TypeString) ReprKind() ReprKind {
	return ReprKind_String
}
func (TypeBytes) ReprKind() ReprKind {
	return ReprKind_Bytes
}
func (TypeInt) ReprKind() ReprKind {
	return ReprKind_Int
}
func (TypeFloat) ReprKind() ReprKind {
	return ReprKind_Float
}
func (TypeMap) ReprKind() ReprKind {
	return ReprKind_Map
}
func (TypeList) ReprKind() ReprKind {
	return ReprKind_List
}
func (TypeLink) ReprKind() ReprKind {
	return ReprKind_Link
}
func (tv TypeUnion) ReprKind() ReprKind {
	// REVIEW: this may fib; has the bizarre property of being dependent on the *concrete value* for kinded unions!
	if tv.Style == UnionStyle_Kinded {
		return ReprKind_Invalid
	} else {
		return ReprKind_Map
	}
}
func (tv TypeObject) ReprKind() ReprKind {
	if tv.TupleStyle {
		return ReprKind_List
	} else {
		return ReprKind_Map
	}
}
func (TypeEnum) ReprKind() ReprKind {
	return ReprKind_String
}
