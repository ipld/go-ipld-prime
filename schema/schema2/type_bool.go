package schema

import (
	"github.com/ipld/go-ipld-prime"
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

func (TypeBool) Kind() Kind {
	return Kind_Bool
}

func (t *TypeBool) Name() TypeName {
	return t.name
}

func (t TypeBool) RepresentationBehavior() ipld.ReprKind {
	return ipld.ReprKind_Bool
}
