package schema

// Kind is an enum of kind in the IPLD Schema system.
//
// Note that schema.Kind is distinct from ipld.ReprKind!
// Schema kinds include concepts such as "struct" and "enum", which are
// concepts only introduced by the Schema layer, and not present in the
// Data Model layer.
type Kind uint8

const (
	Kind_Invalid Kind = 0
	Kind_Map     Kind = '{'
	Kind_List    Kind = '['
	Kind_Unit    Kind = '1'
	Kind_Bool    Kind = 'b'
	Kind_Int     Kind = 'i'
	Kind_Float   Kind = 'f'
	Kind_String  Kind = 's'
	Kind_Bytes   Kind = 'x'
	Kind_Link    Kind = '/'
	Kind_Union   Kind = '^'
	Kind_Struct  Kind = '$'
	Kind_Enum    Kind = '%'
	// FUTURE: Kind_Any = '?'?
)

func (k Kind) String() string {
	switch k {
	case Kind_Invalid:
		return "Invalid"
	case Kind_Map:
		return "Map"
	case Kind_List:
		return "List"
	case Kind_Unit:
		return "Unit"
	case Kind_Bool:
		return "Bool"
	case Kind_Int:
		return "Int"
	case Kind_Float:
		return "Float"
	case Kind_String:
		return "String"
	case Kind_Bytes:
		return "Bytes"
	case Kind_Link:
		return "Link"
	case Kind_Union:
		return "Union"
	case Kind_Struct:
		return "Struct"
	case Kind_Enum:
		return "Enum"
	default:
		panic("invalid enumeration value!")
	}
}
