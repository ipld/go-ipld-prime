package compiler

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeLink struct {
	ts              *TypeSystem
	name            schema.TypeName
	expectedTypeRef schema.TypeName // can be empty
}

// -- schema.Type interface satisfaction -->

var _ schema.Type = (*TypeLink)(nil)

func (t *TypeLink) TypeSystem() schema.TypeSystem {
	return t.ts
}

func (TypeLink) TypeKind() schema.TypeKind {
	return schema.TypeKind_Link
}

func (t *TypeLink) Name() schema.TypeName {
	return t.name
}

func (t TypeLink) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Link
}

// -- specific to TypeLink -->

// HasExpectedType returns true if the link has a hint about the type it references.
func (t *TypeLink) HasExpectedType() bool {
	return t.expectedTypeRef != ""
}

// ExpectedType returns the type which is expected for the node on the other side of the link.
// Nil is returned if there is no information about the expected type
// (which may be interpreted as "any").
func (t *TypeLink) ExpectedType() schema.Type {
	if !t.HasExpectedType() {
		return nil
	}
	return t.ts.types[schema.TypeReference(t.expectedTypeRef)]
}
