package schema

import (
	"github.com/ipld/go-ipld-prime"
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

func (TypeInt) Kind() Kind {
	return Kind_Int
}

func (t *TypeInt) Name() TypeName {
	return t.name
}

func (t TypeInt) RepresentationBehavior() ipld.ReprKind {
	return ipld.ReprKind_Int
}
