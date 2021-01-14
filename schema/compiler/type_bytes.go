package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeBytes struct {
	ts   *TypeSystem
	name schema.TypeName
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeBytes)(nil)

func (t *TypeBytes) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeBytes) TypeKind() schema.TypeKind {
	return schema.TypeKind_Bytes
}

func (t *TypeBytes) Name() schema.TypeName {
	return t.name
}

func (t TypeBytes) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Bytes
}
