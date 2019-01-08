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
