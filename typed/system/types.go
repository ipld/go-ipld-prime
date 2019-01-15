package typesystem

import (
	ipld "github.com/ipld/go-ipld-prime"
)

// FIXME make everything exported in this file a read-only accessor.

// TODO make most references to `Type` into `*Type`.

// TODO rename `TypeObject` into `TypeStruct`.

type TypeName string

type Type interface {
	// Name() TypeName // annoying name collision.
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
	_ Type = TypeObject{}
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
	KeyType       Type // must be ReprKind==string (e.g. Type==String|Enum).
	ValueType     Type
	ValueNullable bool
}
type TypeList struct {
	Name          TypeName
	Anon          bool
	ValueType     Type
	ValueNullable bool
}
type TypeLink struct {
	Name TypeName
	// ...?
}

type TypeUnion struct {
	Name         TypeName
	Style        UnionStyle
	ValuesKinded map[ipld.ReprKind]Type // for Style==Kinded
	Values       map[TypeName]Type      // for Style!=Kinded
	TypeHintKey  string                 // for Style==Envelope|Inline
	ContentKey   string                 // for Style==Envelope
}
type UnionStyle struct{ x string }

var (
	UnionStyle_Kinded   = UnionStyle{"kinded"}
	UnionStyle_Keyed    = UnionStyle{"keyed"}
	UnionStyle_Envelope = UnionStyle{"envelope"}
	UnionStyle_Inline   = UnionStyle{"inline"}
)

type TypeObject struct {
	Name       TypeName
	TupleStyle bool // if true, ReprKind=Array instead of map (and optional fields are invalid!)
	Fields     []ObjectField
}
type ObjectField struct {
	Name     string
	Type     Type
	Optional bool
	Nullable bool
}

type TypeEnum struct {
	Name   string
	Values []string
}
