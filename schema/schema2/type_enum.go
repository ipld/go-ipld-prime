package schema

import (
	"github.com/ipld/go-ipld-prime"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeEnum struct {
	name TypeName
	dmt  schemadmt.TypeEnum
	ts   *TypeSystem
}

type EnumRepresentation interface{ _EnumRepresentation() }

func (EnumRepresentation_String) _EnumRepresentation() {}
func (EnumRepresentation_Int) _EnumRepresentation()    {}

type EnumRepresentation_String struct {
	dmt schemadmt.EnumRepresentation_String
}
type EnumRepresentation_Int struct {
	dmt schemadmt.EnumRepresentation_Int
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeEnum)(nil)

func (t *TypeEnum) _Type() {}

func (t *TypeEnum) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeEnum) Kind() Kind {
	return Kind_Struct
}

func (t *TypeEnum) Name() TypeName {
	return t.name
}

func (t TypeEnum) RepresentationBehavior() ipld.ReprKind {
	switch t.dmt.FieldRepresentation().AsInterface().(type) {
	case schemadmt.EnumRepresentation_String:
		return ipld.ReprKind_String
	case schemadmt.EnumRepresentation_Int:
		return ipld.ReprKind_Int
	default:
		panic("unreachable")
	}
}

// -- specific to TypeEnum -->

func (t *TypeEnum) RepresentationStrategy() EnumRepresentation {
	switch x := t.dmt.FieldRepresentation().AsInterface().(type) {
	case schemadmt.EnumRepresentation_String:
		return EnumRepresentation_String{x}
	case schemadmt.EnumRepresentation_Int:
		return EnumRepresentation_Int{x}
	default:
		panic("unreachable")
	}
}
