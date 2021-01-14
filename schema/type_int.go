package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeInt struct {
	ts   *TypeSystem
	name TypeName
}

// -- Type interface satisfaction -->

var _ Type = (*TypeInt)(nil)

func (t *TypeInt) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeInt) TypeKind() TypeKind {
	return TypeKind_Int
}

func (t *TypeInt) Name() TypeName {
	return t.name
}

func (t TypeInt) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Int
}
