package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeString struct {
	ts   *TypeSystem
	name schema.TypeName
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeString)(nil)

func (t *TypeString) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeString) TypeKind() schema.TypeKind {
	return schema.TypeKind_String
}

func (t *TypeString) Name() schema.TypeName {
	return t.name
}

func (t TypeString) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_String
}
