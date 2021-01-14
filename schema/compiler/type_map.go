package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeMap struct {
	ts            *TypeSystem
	name          schema.TypeName
	keyTypeRef    schema.TypeName // is a TypeName and not a TypeReference because it can't be an anon.
	valueTypeRef  schema.TypeReference
	valueNullable bool
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeMap)(nil)

func (t *TypeMap) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeMap) TypeKind() schema.TypeKind {
	return schema.TypeKind_Map
}

func (t *TypeMap) Name() schema.TypeName {
	return t.name
}

func (t TypeMap) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Map
}

// -- specific to TypeMap -->

// KeyType returns the Type of the map keys.
//
// Note that map keys must always be some type which is representable as a
// string in the IPLD Data Model (e.g. any string type is valid,
// but something with enum typekind and a string representation is also valid,
// and a struct typekind with a representation that has a string kind is also valid, etc).
func (t *TypeMap) KeyType() schema.Type {
	return t.ts.types[schema.TypeReference(t.keyTypeRef)]
}

// ValueType returns the Type of the map values.
func (t *TypeMap) ValueType() schema.Type {
	return t.ts.types[schema.TypeReference(t.valueTypeRef)]
}

// ValueIsNullable returns a bool describing if the map values are permitted to be null.
func (t *TypeMap) ValueIsNullable() bool {
	return t.valueNullable
}
