package schema

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeLink struct {
	name TypeName
	dmt  schemadmt.TypeLink
	ts   *TypeSystem
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeLink)(nil)

func (t *TypeLink) _Type() {}

func (t *TypeLink) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeLink) TypeKind() TypeKind {
	return TypeKind_Link
}

func (t *TypeLink) Name() TypeName {
	return t.name
}

func (t TypeLink) RepresentationBehavior() datamodel.Kind {
	return datamodel.Kind_Link
}

// -- specific to TypeLink -->

// HasReferencedType returns true if the link has a hint about the type it references.
func (t *TypeLink) HasReferencedType() bool {
	return t.dmt.FieldExpectedType().Exists()
}

// ReferencedType returns the type which is expected for the node on the other side of the link.
// Nil is returned if there is no information about the expected type
// (which may be interpreted as "any").
func (t *TypeLink) ReferencedType() Type {
	if !t.dmt.FieldExpectedType().Exists() {
		return nil
	}
	return t.ts.types[t.dmt.FieldExpectedType().Must().TypeReference()]
}
