package schemadmt

import "fmt"

// TypeReference is a string that's either a TypeName or a computed string from an InlineDefn.
// This string is often useful as a map key.
//
// The computed string for an InlineDefn happens to match the IPLD Schema DSL syntax,
// but it would be very odd for any code to depend on that detail.
type TypeReference string

func (x TypeNameOrInlineDefn) TypeReference() TypeReference {
	switch y := x.AsInterface().(type) {
	case TypeName:
		return TypeReference(y.String())
	case TypeDefnInline:
		return y.TypeReference()
	default:
		panic("unreachable")
	}
}

func (x TypeDefnInline) TypeReference() TypeReference {
	switch y := x.AsInterface().(type) {
	case TypeMap:
		if y.FieldValueNullable().Bool() {
			return TypeReference(fmt.Sprintf("{%s : nullable %s}", y.FieldKeyType(), y.FieldValueType().TypeReference()))
		}
		return TypeReference(fmt.Sprintf("{%s:%s}", y.FieldKeyType(), y.FieldValueType().TypeReference()))
	case TypeList:
		if y.FieldValueNullable().Bool() {
			return TypeReference(fmt.Sprintf("[nullable %s]", y.FieldValueType().TypeReference()))
		}
		return TypeReference(fmt.Sprintf("[%s]", y.FieldValueType().TypeReference()))
	default:
		panic("unreachable")
	}
}

func (x TypeName) TypeReference() TypeReference {
	return TypeReference(x.String())
}
