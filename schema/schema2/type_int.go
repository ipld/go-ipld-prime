package schema

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeInt struct {
	name TypeName
	dmt  schemadmt.TypeInt
	ts   *TypeSystem
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeInt)(nil)

func (t *TypeInt) _Type() {}

func (t *TypeInt) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeInt) TypeKind() TypeKind {
	return TypeKind_Int
}

func (t *TypeInt) Name() TypeName {
	return t.name
}

func (t TypeInt) RepresentationBehavior() datamodel.Kind {
	return datamodel.Kind_Int
}
