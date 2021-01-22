package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeLink struct {
	ts              *TypeSystem
	name            TypeName      // may be empty if this is an anon type
	ref             TypeReference // may be dup of name, if this is a named type
	expectedTypeRef TypeName      // can be empty
}

// -- Type interface satisfaction -->

var _ Type = (*TypeLink)(nil)

func (t *TypeLink) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeLink) TypeKind() TypeKind {
	return TypeKind_Link
}

func (t *TypeLink) Name() TypeName {
	return t.name
}

func (t *TypeLink) Reference() TypeReference {
	return t.ref
}

func (t TypeLink) RepresentationBehavior() ipld.Kind {
	return ipld.Kind_Link
}

// -- specific to TypeLink -->

// HasExpectedType returns true if the link has a hint about the type it references.
func (t *TypeLink) HasReferencedType() bool {
	return t.expectedTypeRef != ""
}

// ExpectedType returns the type which is expected for the node on the other side of the link.
// Nil is returned if there is no information about the expected type
// (which may be interpreted as "any").
func (t *TypeLink) ReferencedType() Type {
	if !t.HasReferencedType() {
		return nil
	}
	return t.ts.types[TypeReference(t.expectedTypeRef)]
}
