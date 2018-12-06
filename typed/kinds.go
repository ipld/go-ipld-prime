package typed

// ReprKind represents the primitive kind in the IPLD data model.
// Note that it contains the concept of "map", but not "object" --
// "object" is a concept up in our type system layers, and *not*
// present in the data model layer.
type ReprKind uint8

const (
	ReprKind_Invalid = 0
	ReprKind_Bool    = 'b'
	ReprKind_String  = 's'
	ReprKind_Bytes   = 'x'
	ReprKind_Int     = 'i'
	ReprKind_Float   = 'f'
	ReprKind_Map     = '{' // still feel these should be string-only keys.  consider the Type.Fields behavior if it's an int-keyed map: insane?
	ReprKind_List    = '['
	ReprKind_Null    = '-'
	ReprKind_Link    = '/'
)

// Kind is our type level kind.  It includes "object", "union", and other
// advanced concepts.  "map" at this layer also contains additional constraints;
// it must be a single type of element.
type Kind uint8

const (
	Kind_Invalid = 0
	Kind_Bool    = 'b'
	Kind_String  = 's'
	Kind_Bytes   = 'x'
	Kind_Int     = 'i'
	Kind_Float   = 'f'
	Kind_Map     = '{'
	Kind_List    = '['
	Kind_Null    = '-'
	Kind_Link    = '/'
	Kind_Union   = 'u'
	Kind_Obj     = 'o'
	Kind_Enum    = 'e'
)
