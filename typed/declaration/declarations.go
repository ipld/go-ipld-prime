package typedeclaration

import (
	ipld "github.com/ipld/go-ipld-prime"
)

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

type TypeBool struct {
	Name TypeName
}

type TypeString struct {
	Name TypeName
}

type TypeBytes struct {
	Name TypeName
}

type TypeInt struct {
	Name TypeName
}

type TypeFloat struct {
	Name TypeName
}

type TypeMap struct {
	Name          TypeName
	Anon          bool
	KeyType       TypeName
	ValueType     TypeName
	ValueNullable bool
}

type TypeList struct {
	Name          TypeName
	Anon          bool
	ValueType     TypeName
	ValueNullable bool
}

type TypeLink struct {
	Name TypeName
	// ...?
}

type TypeUnion struct {
	Name         TypeName
	Style        UnionStyle
	ValuesKinded map[ipld.ReprKind]TypeName // for Style==Kinded
	Values       map[string]TypeName        // for Style!=Kinded
	TypeHintKey  string                     // for Style==Envelope|Inline
	ContentKey   string                     // for Style==Envelope
}

type UnionStyle struct{ x string }

var (
	UnionStyle_Kinded   = UnionStyle{"kinded"}
	UnionStyle_Keyed    = UnionStyle{"keyed"}
	UnionStyle_Envelope = UnionStyle{"envelope"}
	UnionStyle_Inline   = UnionStyle{"inline"}
)

type TypeStruct struct {
	Name       TypeName
	TupleStyle bool // if true, ReprKind=Array instead of map (and optional fields are invalid!)
	Fields     []StructField
}

type StructField struct {
	Name     string
	Type     TypeName
	Optional bool
	Nullable bool
}

type TypeEnum struct {
	Name   TypeName
	Values []string
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
