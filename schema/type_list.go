package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeList struct {
	ts            *TypeSystem
	name          TypeName      // may be empty if this is an anon type
	ref           TypeReference // may be dup of name, if this is a named type
	valueTypeRef  TypeReference
	valueNullable bool
}

// -- Type interface satisfaction -->

var _ Type = (*TypeList)(nil)

func (t *TypeList) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeList) TypeKind() TypeKind {
	return TypeKind_List
}

func (t *TypeList) Name() TypeName {
	return t.name
}

func (t *TypeList) Reference() TypeReference {
	return t.ref
}

func (t TypeList) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_List
}

// -- specific to TypeList -->

// ValueType returns the Type of the list values.
func (t *TypeList) ValueType() Type {
	return t.ts.types[TypeReference(t.valueTypeRef)]
}

// ValueIsNullable returns a bool describing if the list values are permitted to be null.
func (t *TypeList) ValueIsNullable() bool {
	return t.valueNullable
}
