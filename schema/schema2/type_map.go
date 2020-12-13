package schema

import (
	"github.com/ipld/go-ipld-prime"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeMap struct {
	name TypeName
	dmt  schemadmt.TypeMap
	ts   *TypeSystem
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeMap)(nil)

func (t *TypeMap) _Type() {}

func (t *TypeMap) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeMap) Kind() Kind {
	return Kind_Map
}

func (t *TypeMap) Name() TypeName {
	return t.name
}

func (t TypeMap) RepresentationBehavior() ipld.ReprKind {
	return ipld.ReprKind_Map
}

// -- specific to TypeMap -->

// KeyType returns the Type of the map keys.
//
// Note that map keys will must always be some type which is representable as a
// string in the IPLD Data Model (e.g. either TypeString or TypeEnum).
func (t *TypeMap) KeyType() Type {
	return t.ts.types[t.dmt.FieldKeyType().TypeReference()]
}

// ValueType returns the Type of the map values.
func (t *TypeMap) ValueType() Type {
	return t.ts.types[t.dmt.FieldValueType().TypeReference()]
}

// ValueIsNullable returns a bool describing if the map values are permitted
// to be null.
func (t *TypeMap) ValueIsNullable() bool {
	return t.dmt.FieldValueNullable().Bool()
}
