package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeString struct {
	ts   *TypeSystem
	name TypeName
}

// -- Type interface satisfaction -->

var _ Type = (*TypeString)(nil)

func (t *TypeString) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeString) TypeKind() TypeKind {
	return TypeKind_String
}

func (t *TypeString) Name() TypeName {
	return t.name
}

func (t TypeString) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_String
}
