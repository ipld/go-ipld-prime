package schema

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeStruct struct {
	name TypeName
	dmt  schemadmt.TypeStruct
	ts   *TypeSystem
}

type StructField struct {
	parent *TypeStruct
	name   schemadmt.FieldName
	dmt    schemadmt.StructField
}

type StructRepresentation interface{ _StructRepresentation() }

func (StructRepresentation_Map) _StructRepresentation()         {}
func (StructRepresentation_Tuple) _StructRepresentation()       {}
func (StructRepresentation_Stringpairs) _StructRepresentation() {}
func (StructRepresentation_Stringjoin) _StructRepresentation()  {}
func (StructRepresentation_Listpairs) _StructRepresentation()   {}

type StructRepresentation_Map struct {
	parent *TypeStruct // this one needs a pointer back up to figure out its defaults.
	dmt    schemadmt.StructRepresentation_Map
}
type StructRepresentation_Tuple struct {
	dmt schemadmt.StructRepresentation_Tuple
}
type StructRepresentation_Stringpairs struct {
	dmt schemadmt.StructRepresentation_Stringpairs
}
type StructRepresentation_Stringjoin struct {
	dmt schemadmt.StructRepresentation_Stringjoin
}
type StructRepresentation_Listpairs struct {
	dmt schemadmt.StructRepresentation_Listpairs
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeStruct)(nil)

func (t *TypeStruct) _Type() {}

func (t *TypeStruct) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeStruct) Kind() Kind {
	return Kind_Struct
}

func (t *TypeStruct) Name() TypeName {
	return t.name
}

func (t TypeStruct) RepresentationBehavior() ipld.ReprKind {
	switch t.dmt.FieldRepresentation().AsInterface().(type) {
	case schemadmt.StructRepresentation_Map:
		return ipld.ReprKind_Map
	case schemadmt.StructRepresentation_Tuple:
		return ipld.ReprKind_List
	case schemadmt.StructRepresentation_Stringpairs:
		return ipld.ReprKind_String
	case schemadmt.StructRepresentation_Stringjoin:
		return ipld.ReprKind_String
	case schemadmt.StructRepresentation_Listpairs:
		return ipld.ReprKind_List
	default:
		panic("unreachable")
	}
}

// -- specific to TypeStruct -->

// Fields returns a slice of descriptions of the object's fields.
func (t *TypeStruct) Fields() []StructField {
	a := make([]StructField, 0, t.dmt.FieldFields().Length())
	for itr := t.dmt.FieldFields().Iterator(); itr.Done(); {
		k, v := itr.Next()
		a = append(a, StructField{t, k, v})
	}
	return a
}

// Field looks up a StructField by name, or returns nil if no such field.
func (t *TypeStruct) Field(name string) *StructField {
	fndmt, err := schemadmt.Type.FieldName.FromString(name)
	if err != nil {
		panic(fmt.Errorf("invalid fieldname: %w", err))
	}
	fdmt := t.dmt.FieldFields().Lookup(fndmt)
	if fdmt == nil {
		return nil
	}
	return &StructField{t, fndmt, fdmt}
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
func (f *StructField) Name() string { return f.name.String() }

// Type returns the Type of this field's value.  Note the field may
// also be unset if it is either Optional or Nullable.
func (f *StructField) Type() Type { return f.parent.ts.types[f.dmt.FieldType().TypeReference()] }

// IsOptional returns true if the field is allowed to be absent from the object.
// If IsOptional is false, the field may be absent from the serial representation
// of the object entirely.
//
// Note being optional is different than saying the value is permitted to be null!
// A field may be both nullable and optional simultaneously, or either, or neither.
func (f *StructField) IsOptional() bool { return f.dmt.FieldOptional().Bool() }

// IsNullable returns true if the field value is allowed to be null.
//
// If is Nullable is false, note that it's still possible that the field value
// will be absent if the field is Optional!  Being nullable is unrelated to
// whether the field's presence is optional as a whole.
//
// Note that a field may be both nullable and optional simultaneously,
// or either, or neither.
func (f *StructField) IsNullable() bool { return f.dmt.FieldNullable().Bool() }

// IsMaybe returns true if the field value is allowed to be either null or absent.
//
// This is a simple "or" of the two properties,
// but this method is a shorthand that turns out useful often.
func (f *StructField) IsMaybe() bool { return f.IsNullable() || f.IsOptional() }

func (t *TypeStruct) RepresentationStrategy() StructRepresentation {
	switch x := t.dmt.FieldRepresentation().AsInterface().(type) {
	case schemadmt.StructRepresentation_Map:
		return StructRepresentation_Map{t, x}
	case schemadmt.StructRepresentation_Tuple:
		return StructRepresentation_Tuple{x}
	case schemadmt.StructRepresentation_Stringpairs:
		return StructRepresentation_Stringpairs{x}
	case schemadmt.StructRepresentation_Stringjoin:
		return StructRepresentation_Stringjoin{x}
	case schemadmt.StructRepresentation_Listpairs:
		return StructRepresentation_Listpairs{x}
	default:
		panic("unreachable")
	}
}

// GetFieldKey returns the string that should be the key when serializing this field.
// For some fields, it's the same as the field name; for others, a rename directive may provide a different value.
func (r StructRepresentation_Map) GetFieldKey(field StructField) string {
	maybeOverrides := r.dmt.FieldFields()
	if !maybeOverrides.Exists() {
		return field.Name()
	}
	fieldInfo := maybeOverrides.Must().Lookup(field.name)
	if fieldInfo == nil {
		return field.Name()
	}
	maybeRename := fieldInfo.FieldRename()
	if !maybeRename.Exists() {
		return field.Name()
	}
	return maybeRename.Must().String()
}

func (r StructRepresentation_Stringjoin) GetJoinDelim() string {
	return r.dmt.FieldJoin().String()
}
