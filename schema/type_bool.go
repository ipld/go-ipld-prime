package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeBool struct {
	ts   *TypeSystem
	name TypeName
}

// -- Type interface satisfaction -->

var _ Type = (*TypeBool)(nil)

func (t *TypeBool) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeBool) TypeKind() TypeKind {
	return TypeKind_Bool
}

func (t *TypeBool) Name() TypeName {
	return t.name
}

func (t *TypeBool) Reference() TypeReference {
	return TypeReference(t.name)
}

func (t TypeBool) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Bool
}
