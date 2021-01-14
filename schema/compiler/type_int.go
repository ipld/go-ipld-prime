package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeInt struct {
	ts   *TypeSystem
	name schema.TypeName
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeInt)(nil)

func (t *TypeInt) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeInt) TypeKind() schema.TypeKind {
	return schema.TypeKind_Int
}

func (t *TypeInt) Name() schema.TypeName {
	return t.name
}

func (t TypeInt) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Int
}
