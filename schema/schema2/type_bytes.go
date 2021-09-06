package schema

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeBytes struct {
	name TypeName
	dmt  schemadmt.TypeBytes
	ts   *TypeSystem
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeBytes)(nil)

func (t *TypeBytes) _Type() {}

func (t *TypeBytes) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeBytes) TypeKind() TypeKind {
	return TypeKind_Bytes
}

func (t *TypeBytes) Name() TypeName {
	return t.name
}

func (t TypeBytes) RepresentationBehavior() datamodel.Kind {
	return datamodel.Kind_Bytes
}
