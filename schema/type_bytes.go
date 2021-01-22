package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeBytes struct {
	ts   *TypeSystem
	name TypeName
}

// -- Type interface satisfaction -->

var _ Type = (*TypeBytes)(nil)

func (t *TypeBytes) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeBytes) TypeKind() TypeKind {
	return TypeKind_Bytes
}

func (t *TypeBytes) Name() TypeName {
	return t.name
}

func (t *TypeBytes) Reference() TypeReference {
	return TypeReference(t.name)
}

func (t TypeBytes) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Bytes
}
