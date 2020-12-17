package schema

import (
	"github.com/ipld/go-ipld-prime"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeFloat struct {
	name TypeName
	dmt  schemadmt.TypeFloat
	ts   *TypeSystem
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeFloat)(nil)

func (t *TypeFloat) _Type() {}

func (t *TypeFloat) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeFloat) TypeKind() TypeKind {
	return TypeKind_Float
}

func (t *TypeFloat) Name() TypeName {
	return t.name
}

func (t TypeFloat) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Float
}
