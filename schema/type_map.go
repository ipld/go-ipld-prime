package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeMap struct {
	ts            *TypeSystem
	name          TypeName
	keyTypeRef    TypeName // is a TypeName and not a TypeReference because it can't be an anon.
	valueTypeRef  TypeReference
	valueNullable bool
	rstrat        MapRepresentation
}

func (t *TypeMap) Representation() MapRepresentation {
	return t.rstrat
}

type MapRepresentation interface{ _MapRepresentation() }

func (MapRepresentation_Map) _MapRepresentation()         {}
func (MapRepresentation_Listpairs) _MapRepresentation()   {}
func (MapRepresentation_Stringpairs) _MapRepresentation() {}

type MapRepresentation_Map struct{}
type MapRepresentation_Listpairs struct{}
type MapRepresentation_Stringpairs struct {
	innerDelim string
	entryDelim string
}

// -- Type interface satisfaction -->

var _ Type = (*TypeMap)(nil)

func (t *TypeMap) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeMap) TypeKind() TypeKind {
	return TypeKind_Map
}

func (t *TypeMap) Name() TypeName {
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
func (t *TypeMap) KeyType() Type {
	return t.ts.types[TypeReference(t.keyTypeRef)]
}

// ValueType returns the Type of the map values.
func (t *TypeMap) ValueType() Type {
	return t.ts.types[TypeReference(t.valueTypeRef)]
}

// ValueIsNullable returns a bool describing if the map values are permitted to be null.
func (t *TypeMap) ValueIsNullable() bool {
	return t.valueNullable
}
