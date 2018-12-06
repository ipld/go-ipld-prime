package typed

// Kind is our type level kind.  It includes "object", "union", and other
// advanced concepts.  "map" at this layer also contains additional constraints;
// it must be a single type of element.
type Kind uint8

// REVIEW: unclear if this is needed.  Can switch on `Type.(type)`.

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
