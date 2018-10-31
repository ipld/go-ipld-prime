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

type Type struct {
	Name     string // must be unique when scoped to this universe
	Kind            // the type kind.
	ReprKind        // the representation kind.  if wildly different than the type kind, implies a transform func.

	Elem         *Type // only valid kind Kind==Map or Kind==List
	ElemNullable bool  // only valid kind Kind==Map or Kind==List

	Fields map[string]ObjField // only valid if Kind==Obj.

	Union map[string]*Type // only valid if Kind==Union.

	EnumVals []interface{} // only valid if Kind==Enum.  can be ReprKind int, str.  (maybe bytes?  discuss.)
}

type ObjField struct {
	Name     string // must be unique in scope of this object
	Type     *Type
	Required bool // required means the field must be present for the obj to be recognized (but no comment on if null is valid).
	Nullable bool // nullable means that when present, the field may be null.
}

type Universe map[string]*Type

/*

type Foo {
  f1: String
  f2: [String!]!
  ?f3: String!
}

*/
var (
	// Prelude types
	tString = &Type{
		Name:     "String",
		Kind:     Kind_String,
		ReprKind: ReprKind_String,
	}
	// User's types
	tFoo = &Type{
		Name:     "Foo",
		Kind:     Kind_Obj,
		ReprKind: ReprKind_Map,
		Fields: map[string]ObjField{
			"f1": {"f1", tString, true, true},
			"f2": {"f2", &Type{
				Name:         "",
				Kind:         Kind_List,
				ReprKind:     Kind_List,
				Elem:         tString,
				ElemNullable: false,
			}, true, false},
			"f3": {"f3", tString, false, false},
		},
	}
	// The Universe
	example = Universe{
		tFoo.Name: tFoo,
	}
)
