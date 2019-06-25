package schema

import (
	"github.com/ipld/go-ipld-prime"
)

/* cookie-cutter standard interface stuff */

func (anyType) _Type()                    {}
func (t anyType) TypeSystem() *TypeSystem { return t.universe }
func (t anyType) Name() TypeName          { return t.name }

func (TypeBool) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Bool
}
func (TypeString) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (TypeBytes) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Bytes
}
func (TypeInt) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Int
}
func (TypeFloat) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Float
}
func (TypeMap) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (TypeList) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_List
}
func (TypeLink) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Link
}
func (t TypeUnion) ReprKind() ipld.ReprKind {
	if t.style == UnionStyle_Kinded {
		return ipld.ReprKind_Invalid
	} else {
		return ipld.ReprKind_Map
	}
}
func (t TypeStruct) ReprKind() ipld.ReprKind {
	if t.tupleStyle {
		return ipld.ReprKind_List
	} else {
		return ipld.ReprKind_Map
	}
}
func (TypeEnum) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}

/* interesting methods per Type type */

// IsAnonymous is returns true if the type was unnamed.  Unnamed types will
// claim to have a Name property like `{Foo:Bar}`, and this is not guaranteed
// to be a unique string for all types in the universe.
func (t TypeMap) IsAnonymous() bool {
	return t.anonymous
}

// KeyType returns the Type of the map keys.
//
// Note that map keys will must always be some type which is representable as a
// string in the IPLD Data Model (e.g. either TypeString or TypeEnum).
func (t TypeMap) KeyType() Type {
	return t.keyType
}

// ValueType returns to the Type of the map values.
func (t TypeMap) ValueType() Type {
	return t.valueType
}

// ValueIsNullable returns a bool describing if the map values are permitted
// to be null.
func (t TypeMap) ValueIsNullable() bool {
	return t.valueNullable
}

// IsAnonymous is returns true if the type was unnamed.  Unnamed types will
// claim to have a Name property like `[Foo]`, and this is not guaranteed
// to be a unique string for all types in the universe.
func (t TypeList) IsAnonymous() bool {
	return t.anonymous
}

// ValueType returns to the Type of the list values.
func (t TypeList) ValueType() Type {
	return t.valueType
}

// ValueIsNullable returns a bool describing if the list values are permitted
// to be null.
func (t TypeList) ValueIsNullable() bool {
	return t.valueNullable
}

// UnionMembers returns a set of all the types that can inhabit this Union.
func (t TypeUnion) UnionMembers() map[Type]struct{} {
	m := make(map[Type]struct{}, len(t.values)+len(t.valuesKinded))
	switch t.style {
	case UnionStyle_Kinded:
		for _, v := range t.valuesKinded {
			m[v] = struct{}{}
		}
	default:
		for _, v := range t.values {
			m[v] = struct{}{}
		}
	}
	return m
}

// Fields returns a slice of descriptions of the object's fields.
func (t TypeStruct) Fields() []StructField {
	a := make([]StructField, len(t.fields))
	for i := range t.fields {
		a[i] = t.fields[i]
	}
	return a
}

// Name returns the string name of this field.  The name is the string that
// will be used as a map key if the structure this field is a member of is
// serialized as a map representation.
func (f StructField) Name() string { return f.name }

// Type returns the Type of this field's value.  Note the field may
// also be unset if it is either Optional or Nullable.
func (f StructField) Type() Type { return f.typ }

// IsOptional returns true if the field is allowed to be absent from the object.
// If IsOptional is false, the field may be absent from the serial representation
// of the object entirely.
//
// Note being optional is different than saying the value is permitted to be null!
// A field may be both nullable and optional simultaneously, or either, or neither.
func (f StructField) IsOptional() bool { return f.optional }

// IsNullable returns true if the field value is allowed to be null.
//
// If is Nullable is false, note that it's still possible that the field value
// will be absent if the field is Optional!  Being nullable is unrelated to
// whether the field's presence is optional as a whole.
//
// Note that a field may be both nullable and optional simultaneously,
// or either, or neither.
func (f StructField) IsNullable() bool { return f.nullable }

// Members returns a slice the strings which are valid inhabitants of this enum.
func (t TypeEnum) Members() []string {
	a := make([]string, len(t.members))
	for i := range t.members {
		a[i] = t.members[i]
	}
	return a
}
