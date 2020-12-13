package schema

import (
	"github.com/ipld/go-ipld-prime"
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

func (TypeBytes) Kind() Kind {
	return Kind_Bytes
}

func (t *TypeBytes) Name() TypeName {
	return t.name
}

func (t TypeBytes) RepresentationBehavior() ipld.ReprKind {
	return ipld.ReprKind_Bytes
}
