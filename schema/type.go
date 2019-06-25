package schema

import (
	ipld "github.com/ipld/go-ipld-prime"
)

type TypeName string // = ast.TypeName

// typesystem.Type is an union interface; each of the `Type*` concrete types
// in this package are one of its members.
//
// Specifically,
//
// 	TypeBool
// 	TypeString
// 	TypeBytes
// 	TypeInt
// 	TypeFloat
// 	TypeMap
// 	TypeList
// 	TypeLink
// 	TypeUnion
// 	TypeStruct
// 	TypeEnum
//
// are all of the kinds of Type.
//
// This is a closed union; you can switch upon the above members without
// including a default case.  The membership is closed by the unexported
// '_Type' method; you may use the BurntSushi/go-sumtype tool to check
// your switches for completeness.
//
// Many interesting properties of each Type are only defined for that specific
// type, so it's typical to use a type switch to handle each type of Type.
// (Your humble author is truly sorry for the word-mash that results from
// attempting to describe the types that describe the typesystem.Type.)
//
// For example, to inspect the kind of fields in a struct: you might
// cast a `Type` interface into `TypeStruct`, and then the `Fields()` on
// that `TypeStruct` can be inspected.  (`Fields()` isn't defined for any
// other kind of Type.)
type Type interface {
	// Unexported marker method to force the union closed.
	_Type()

	// Returns a pointer to the typesystem.Universe this type is a member of.
	TypeSystem() *TypeSystem

	// Returns the string name of the Type.  This name is unique within the
	// universe this type is a member of, *unless* this type is Anonymous,
	// in which case a string describing the type will still be returned, but
	// that string will not be required to be unique.
	Name() TypeName

	// Returns the Representation Kind in the IPLD Data Model that this type
	// is expected to be serialized as.
	//
	// Note that in one case, this will return `ipld.ReprKind_Invalid` --
	// TypeUnion with Style=Kinded may be serialized as different kinds
	// depending on their value, so we can't say from the type definition
	// alone what kind we expect.
	ReprKind() ipld.ReprKind
}

var (
	_ Type = TypeBool{}
	_ Type = TypeString{}
	_ Type = TypeBytes{}
	_ Type = TypeInt{}
	_ Type = TypeFloat{}
	_ Type = TypeMap{}
	_ Type = TypeList{}
	_ Type = TypeLink{}
	_ Type = TypeUnion{}
	_ Type = TypeStruct{}
	_ Type = TypeEnum{}
)

type anyType struct {
	name     TypeName
	universe *TypeSystem
}

type TypeBool struct {
	anyType
}

type TypeString struct {
	anyType
}

type TypeBytes struct {
	anyType
}

type TypeInt struct {
	anyType
}

type TypeFloat struct {
	anyType
}

type TypeMap struct {
	anyType
	anonymous     bool
	keyType       Type // must be ReprKind==string (e.g. Type==String|Enum).
	valueType     Type
	valueNullable bool
}

type TypeList struct {
	anyType
	anonymous     bool
	valueType     Type
	valueNullable bool
}

type TypeLink struct {
	anyType
	// ...?
}

type TypeUnion struct {
	anyType
	style        UnionStyle
	valuesKinded map[ipld.ReprKind]Type // for Style==Kinded
	values       map[string]Type        // for Style!=Kinded (note, key is freetext, not necessarily TypeName of the value)
	typeHintKey  string                 // for Style==Envelope|Inline
	contentKey   string                 // for Style==Envelope
}

type UnionStyle struct{ x string }

var (
	UnionStyle_Kinded   = UnionStyle{"kinded"}
	UnionStyle_Keyed    = UnionStyle{"keyed"}
	UnionStyle_Envelope = UnionStyle{"envelope"}
	UnionStyle_Inline   = UnionStyle{"inline"}
)

type TypeStruct struct {
	anyType
	tupleStyle bool // if true, ReprKind=Array instead of map (and optional fields are invalid!)
	fields     []StructField
}
type StructField struct {
	name     string
	typ      Type
	optional bool
	nullable bool
}

type TypeEnum struct {
	anyType
	members []string
}
