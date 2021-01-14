package schemadmt

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/schema"
)

func (x TypeNameOrInlineDefn) TypeReference() schema.TypeReference {
	switch y := x.AsInterface().(type) {
	case TypeName:
		return schema.TypeReference(y.String())
	case TypeDefnInline:
		return y.TypeReference()
	default:
		panic("unreachable")
	}
}

func (x TypeDefnInline) TypeReference() schema.TypeReference {
	switch y := x.AsInterface().(type) {
	case TypeMap:
		if y.FieldValueNullable().Bool() {
			return schema.TypeReference(fmt.Sprintf("{%s : nullable %s}", y.FieldKeyType(), y.FieldValueType().TypeReference()))
		}
		return schema.TypeReference(fmt.Sprintf("{%s:%s}", y.FieldKeyType(), y.FieldValueType().TypeReference()))
	case TypeList:
		if y.FieldValueNullable().Bool() {
			return schema.TypeReference(fmt.Sprintf("[nullable %s]", y.FieldValueType().TypeReference()))
		}
		return schema.TypeReference(fmt.Sprintf("[%s]", y.FieldValueType().TypeReference()))
	default:
		panic("unreachable")
	}
}

func (x TypeName) TypeReference() schema.TypeReference {
	return schema.TypeReference(x.String())
}
