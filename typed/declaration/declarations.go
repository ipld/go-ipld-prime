package typedeclaration

import (
	ipld "github.com/ipld/go-ipld-prime"
)

// Interesting jibblybit: docs actually shouldn't be part of the schema.
// Do the affect the cardinality of any of the members?  Change any behavior?  No.
// Therefore I shouldn't want them to change the CID of the schema as a whole.
//
// Now, obviously this is in dispute with what I want for the AST.
// What shall we do about this?
//
// It's exercise for two schemas applying to the same data, I guess.
// Hopefully we can get it to do something sane when it comes to codegen:
// having double the code for this would be undesirable.
// (I think this is probably a gimme, though.  We'll just use a generic
// typed.Node wrapper (with limited view) around the even-more-typed codegen types.)

type TypeName string

// typedeclaration.Type is a union interface; each of the `Type*` concrete types
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
// In contrast to typesystem.Type, the typedeclaration.Type members are meant
// to represent something more like the AST (Abstract Syntax Tree) that
// *describes* a typesystem.  It's meant to be serializable, so that
// we can be self-describing.  As such, none of the members of
// typedeclaration.Type are permitted to have pointers nor cycles (whereas
// typesystem.Type can and definitely does).
//
// Since we do all of the consistency checking during the final conversion
// and reification of typedeclaration.Type -> typesystem.Type and
// typesystem.Universe, we don't bother much with consistency nor immutablity
// in the design of typedeclaration types; you can use these as builders
// fairly haphazardly, and everything will become nicely immutable after you
// submit it all to typesystem.Construct.
type Type interface {
	// Unexported marker method to force the union closed.
	_Type()
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

type TypeBool struct{}

type TypeString struct{}

type TypeBytes struct{}

type TypeInt struct{}

type TypeFloat struct{}

type TypeMap struct {
	KeyType       TypeName
	ValueType     TypeTerm
	ValueNullable bool
}

type TypeList struct {
	ValueType     TypeTerm
	ValueNullable bool
}

type TypeLink struct {
	// ...?
}

type TypeUnion struct {
	Representation UnionRepresentation
}

type UnionRepresentation interface {
	_UnionRepresentation()
}

type UnionRepresentation_Kinded map[ipld.ReprKind]TypeName

type UnionRepresentation_Keyed map[string]TypeName

type UnionRepresentation_Envelope struct {
	DiscriminatorKey  string
	ContentKey        string
	DiscriminantTable map[string]TypeName
}

type UnionRepresentation_Inline struct {
	DiscriminatorKey  string
	DiscriminantTable map[string]TypeName
}

func (UnionRepresentation_Kinded) _UnionRepresentation()   {}
func (UnionRepresentation_Keyed) _UnionRepresentation()    {}
func (UnionRepresentation_Envelope) _UnionRepresentation() {}
func (UnionRepresentation_Inline) _UnionRepresentation()   {}

type UnionStyle struct{ x string }

var (
	UnionStyle_Kinded   = UnionStyle{"kinded"}
	UnionStyle_Keyed    = UnionStyle{"keyed"}
	UnionStyle_Envelope = UnionStyle{"envelope"}
	UnionStyle_Inline   = UnionStyle{"inline"}
)

type TypeStruct struct {
	Fields         map[string]StructField
	Representation StructRepresentation
}

type StructField struct {
	Type     TypeTerm
	Optional bool
	Nullable bool
}

type TypeTerm interface{} // TODO finish.  not sure where we'll hit this in bootstrap.

type StructRepresentation interface{} // TODO finish.  will be a union.  can assume 'map' for now.

type TypeEnum struct {
	Members map[string]struct{}
}

func (TypeBool) _Type()   {}
func (TypeString) _Type() {}
func (TypeBytes) _Type()  {}
func (TypeInt) _Type()    {}
func (TypeFloat) _Type()  {}
func (TypeMap) _Type()    {}
func (TypeList) _Type()   {}
func (TypeLink) _Type()   {}
func (TypeUnion) _Type()  {}
func (TypeStruct) _Type() {}
func (TypeEnum) _Type()   {}
