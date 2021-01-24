package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeEnum struct {
	name    TypeName
	members []string // a map in the dmt, but really an ordered set.  easier as a slice in golang.
	ts      *TypeSystem
	rstrat  EnumRepresentation
}

func (t *TypeEnum) Representation() EnumRepresentation {
	return t.rstrat
}

type EnumRepresentation interface{ _EnumRepresentation() }

func (EnumRepresentation_String) _EnumRepresentation() {}
func (EnumRepresentation_Int) _EnumRepresentation()    {}

type EnumRepresentation_String struct {
	labels map[string]string // member:label
}
type EnumRepresentation_Int struct {
	labels map[string]int // member:label
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeEnum)(nil)

func (t *TypeEnum) _Type() {}

func (t *TypeEnum) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeEnum) TypeKind() TypeKind {
	return TypeKind_Enum
}

func (t *TypeEnum) Name() TypeName {
	return t.name
}

func (t *TypeEnum) Reference() TypeReference {
	return TypeReference(t.name)
}

func (t TypeEnum) RepresentationBehavior() ipld.Kind {
	switch t.rstrat.(type) {
	case EnumRepresentation_String:
		return ipld.Kind_String
	case EnumRepresentation_Int:
		return ipld.Kind_Int
	default:
		panic("unreachable")
	}
}

// -- specific to TypeEnum -->

func (t *TypeEnum) RepresentationStrategy() EnumRepresentation {
	return t.rstrat
}
