package ipld

// ReprKind represents the primitive kind in the IPLD data model.
// All of these kinds map directly onto serializable data.
//
// Note that ReprKind contains the concept of "map", but not "struct"
// or "object" -- those are a concepts that could be introduced in a
// type system layers, but are *not* present in the data model layer,
// and therefore they aren't included in the ReprKind enum.
type ReprKind uint8

const (
	ReprKind_Invalid ReprKind = 0
	ReprKind_Map     ReprKind = '{'
	ReprKind_List    ReprKind = '['
	ReprKind_Null    ReprKind = '0'
	ReprKind_Bool    ReprKind = 'b'
	ReprKind_Int     ReprKind = 'i'
	ReprKind_Float   ReprKind = 'f'
	ReprKind_String  ReprKind = 's'
	ReprKind_Bytes   ReprKind = 'x'
	ReprKind_Link    ReprKind = '/'
)

func (k ReprKind) String() string {
	switch k {
	case ReprKind_Invalid:
		return "Invalid"
	case ReprKind_Map:
		return "Map"
	case ReprKind_List:
		return "List"
	case ReprKind_Null:
		return "Null"
	case ReprKind_Bool:
		return "Bool"
	case ReprKind_Int:
		return "Int"
	case ReprKind_Float:
		return "Float"
	case ReprKind_String:
		return "String"
	case ReprKind_Bytes:
		return "Bytes"
	case ReprKind_Link:
		return "Link"
	default:
		panic("invalid enumeration value!")
	}
}
