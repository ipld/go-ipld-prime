package schema

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeBool struct {
	name TypeName
	dmt  schemadmt.TypeBool
	ts   *TypeSystem
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeBool)(nil)

func (t *TypeBool) _Type() {}

func (t *TypeBool) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeBool) TypeKind() TypeKind {
	return TypeKind_Bool
}

func (t *TypeBool) Name() TypeName {
	return t.name
}

func (t TypeBool) RepresentationBehavior() datamodel.Kind {
	return datamodel.Kind_Bool
}
