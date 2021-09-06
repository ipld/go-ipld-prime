package schema

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeString struct {
	name TypeName
	dmt  schemadmt.TypeString
	ts   *TypeSystem
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeString)(nil)

func (t *TypeString) _Type() {}

func (t *TypeString) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeString) TypeKind() TypeKind {
	return TypeKind_String
}

func (t *TypeString) Name() TypeName {
	return t.name
}

func (t TypeString) RepresentationBehavior() datamodel.Kind {
	return datamodel.Kind_String
}
