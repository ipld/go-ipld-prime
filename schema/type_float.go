package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeFloat struct {
	ts   *TypeSystem
	name TypeName
}

// -- Type interface satisfaction -->

var _ Type = (*TypeFloat)(nil)

func (t *TypeFloat) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeFloat) TypeKind() TypeKind {
	return TypeKind_Float
}

func (t *TypeFloat) Name() TypeName {
	return t.name
}

func (t *TypeFloat) Reference() TypeReference {
	return TypeReference(t.name)
}

func (t TypeFloat) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Float
}
