package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeList struct {
	ts            *TypeSystem
	name          schema.TypeName
	valueTypeRef  schema.TypeReference
	valueNullable bool
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeList)(nil)

func (t *TypeList) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeList) TypeKind() schema.TypeKind {
	return schema.TypeKind_List
}

func (t *TypeList) Name() schema.TypeName {
	return t.name
}

func (t TypeList) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_List
}

// -- specific to TypeList -->

// ValueType returns the Type of the list values.
func (t *TypeList) ValueType() schema.Type {
	return t.ts.types[schema.TypeReference(t.valueTypeRef)]
}

// ValueIsNullable returns a bool describing if the list values are permitted to be null.
func (t *TypeList) ValueIsNullable() bool {
	return t.valueNullable
}
