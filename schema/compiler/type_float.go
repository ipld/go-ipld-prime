package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeFloat struct {
	ts   *TypeSystem
	name schema.TypeName
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeFloat)(nil)

func (t *TypeFloat) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeFloat) TypeKind() schema.TypeKind {
	return schema.TypeKind_Float
}

func (t *TypeFloat) Name() schema.TypeName {
	return t.name
}

func (t TypeFloat) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Float
}
