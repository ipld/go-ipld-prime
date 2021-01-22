package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeStruct struct {
	ts        *TypeSystem
	name      TypeName
	fields    []StructField
	fieldsMap map[StructFieldName]*StructField // same content, indexed for lookup.
	rstrat    StructRepresentation
}

type StructField struct {
	parent   *TypeStruct // a pointer back up is used so we can provide the method that gives a reified type instead of just the TypeReference.
	name     StructFieldName
	typeRef  TypeReference
	optional bool
	nullable bool
}

type StructFieldName string

func (t *TypeStruct) Representation() StructRepresentation {
	return t.rstrat
}

type StructRepresentation interface{ _StructRepresentation() }

func (StructRepresentation_Map) _StructRepresentation()         {}
func (StructRepresentation_Tuple) _StructRepresentation()       {}
func (StructRepresentation_Stringpairs) _StructRepresentation() {}
func (StructRepresentation_Stringjoin) _StructRepresentation()  {}
func (StructRepresentation_Listpairs) _StructRepresentation()   {}

type StructRepresentation_Map struct {
	parent       *TypeStruct // this one needs a pointer back up to figure out its defaults.
	fieldDetails map[StructFieldName]StructRepresentation_Map_FieldDetails
}
type StructRepresentation_Map_FieldDetails struct {
	Rename   string
	Implicit interface{}
}

type StructRepresentation_Tuple struct {
	fieldOrder []StructFieldName
}

type StructRepresentation_Stringpairs struct {
	innerDelim string
	entryDelim string
}

type StructRepresentation_Stringjoin struct {
	delim      string
	fieldOrder []StructFieldName
}

type StructRepresentation_Listpairs struct {
}

// -- Type interface satisfaction -->

var _ Type = (*TypeStruct)(nil)

func (t *TypeStruct) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeStruct) TypeKind() TypeKind {
	return TypeKind_Struct
}

func (t *TypeStruct) Name() TypeName {
	return t.name
}

func (t TypeStruct) RepresentationBehavior() ipld.Kind {
	switch t.rstrat.(type) {
	case StructRepresentation_Map:
		return ipld.Kind_Map
	case StructRepresentation_Tuple:
		return ipld.Kind_List
	case StructRepresentation_Stringpairs:
		return ipld.Kind_String
	case StructRepresentation_Stringjoin:
		return ipld.Kind_String
	case StructRepresentation_Listpairs:
		return ipld.Kind_List
	default:
		panic("unreachable")
	}
}

// -- specific to TypeStruct -->

// Fields returns a slice of descriptions of the object's fields.
func (t *TypeStruct) Fields() []StructField {
	// A defensive copy to preserve immutability is performed.
	a := make([]StructField, len(t.fields))
	copy(a, t.fields)
	return a
}

// Field looks up a StructField by name, or returns nil if no such field.
func (t *TypeStruct) Field(name string) *StructField {
	return t.fieldsMap[StructFieldName(name)]
}

// Parent returns the type information that this field describes a part of.
//
// While in many cases, you may know the parent already from context,
// there may still be situations where want to pass around a field and
// not need to continue passing down the parent type with it; this method
// helps your code be less redundant in such a situation.
// (You'll find this useful for looking up any rename directives, for example,
// when holding onto a field, since that requires looking up information from
// the representation strategy, which is a property of the type as a whole.)
func (f *StructField) Parent() *TypeStruct { return f.parent }

// Name returns the string name of this field.  The name is the string that
// will be used as a map key if the structure this field is a member of is
// serialized as a map representation.
func (f *StructField) Name() StructFieldName { return f.name }

// Type returns the Type of this field's value.  Note the field may
// also be unset if it is either Optional or Nullable.
func (f *StructField) Type() Type { return f.parent.ts.types[f.typeRef] }

// IsOptional returns true if the field is allowed to be absent from the object.
// If IsOptional is false, the field may be absent from the serial representation
// of the object entirely.
//
// Note being optional is different than saying the value is permitted to be null!
// A field may be both nullable and optional simultaneously, or either, or neither.
func (f *StructField) IsOptional() bool { return f.optional }

// IsNullable returns true if the field value is allowed to be null.
//
// If is Nullable is false, note that it's still possible that the field value
// will be absent if the field is Optional!  Being nullable is unrelated to
// whether the field's presence is optional as a whole.
//
// Note that a field may be both nullable and optional simultaneously,
// or either, or neither.
func (f *StructField) IsNullable() bool { return f.nullable }

// IsMaybe returns true if the field value is allowed to be either null or absent.
//
// This is a simple "or" of the two properties,
// but this method is a shorthand that turns out useful often.
func (f *StructField) IsMaybe() bool { return f.IsNullable() || f.IsOptional() }

func (t *TypeStruct) RepresentationStrategy() StructRepresentation {
	return t.rstrat
}

// GetFieldKey returns the string that should be the key when serializing this field.
// For some fields, it's the same as the field name; for others, a rename directive may provide a different value.
func (r StructRepresentation_Map) GetFieldKey(field StructField) string {
	details, exists := r.fieldDetails[field.name]
	if !exists {
		return string(field.Name())
	}
	if details.Rename == "" {
		return string(field.Name())
	}
	return details.Rename
}

func (r StructRepresentation_Stringjoin) GetJoinDelim() string {
	return r.delim
}
