package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeBool struct {
	ts   *TypeSystem
	name schema.TypeName
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeBool)(nil)

func (t *TypeBool) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeBool) TypeKind() schema.TypeKind {
	return schema.TypeKind_Bool
}

func (t *TypeBool) Name() schema.TypeName {
	return t.name
}

func (t TypeBool) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Bool
}
