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
	ReprKind_Invalid = 0
	ReprKind_Map     = '{'
	ReprKind_List    = '['
	ReprKind_Null    = '0'
	ReprKind_Bool    = 'b'
	ReprKind_Int     = 'i'
	ReprKind_Float   = 'f'
	ReprKind_String  = 's'
	ReprKind_Bytes   = 'x'
	ReprKind_Link    = '/'
)
