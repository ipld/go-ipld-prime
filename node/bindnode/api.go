// Package bindnode provides an datamodel.Node implementation via Go reflection.
//
// This package is EXPERIMENTAL; its behavior and API might change as it's still
// in development.
package bindnode

import (
	"reflect"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

// Prototype implements a schema.TypedPrototype given a Go pointer type and an
// IPLD schema type. Note that the result is also an datamodel.NodePrototype.
//
// If both the Go type and schema type are supplied, it is assumed that they are
// compatible with one another.
//
// If either the Go type or schema type are nil, we infer the missing type from
// the other provided type. For example, we can infer an unnamed Go struct type
// for a schema struct tyep, and we can infer a schema Int type for a Go int64
// type. The inferring logic is still a work in progress and subject to change.
// At this time, inferring IPLD Unions and Enums from Go types is not supported.
//
// When supplying a non-nil ptrType, Prototype only obtains the Go pointer type
// from it, so its underlying value will typically be nil. For example:
//
//     proto := bindnode.Prototype((*goType)(nil), schemaType)
func Prototype(ptrType interface{}, schemaType schema.Type) schema.TypedPrototype {
	if ptrType == nil && schemaType == nil {
		panic("either ptrType or schemaType must not be nil")
	}

	// TODO: if both are supplied, verify that they are compatible

	var goType reflect.Type
	if ptrType == nil {
		goType = inferGoType(schemaType)
	} else {
		goPtrType := reflect.TypeOf(ptrType)
		if goPtrType.Kind() != reflect.Ptr {
			panic("ptrType must be a pointer")
		}
		goType = goPtrType.Elem()
	}

	if schemaType == nil {
		schemaType = inferSchema(goType)
	}

	return &_prototype{schemaType: schemaType, goType: goType}
}

// Wrap implements a schema.TypedNode given a non-nil pointer to a Go value and an
// IPLD schema type. Note that the result is also an datamodel.Node.
//
// Wrap is meant to be used when one already has a Go value with data.
// As such, ptrVal must not be nil.
//
// Similar to Prototype, if schemaType is non-nil it is assumed to be compatible
// with the Go type, and otherwise it's inferred from the Go type.
func Wrap(ptrVal interface{}, schemaType schema.Type) schema.TypedNode {
	if ptrVal == nil {
		panic("ptrVal must not be nil")
	}
	goPtrVal := reflect.ValueOf(ptrVal)
	if goPtrVal.Kind() != reflect.Ptr {
		panic("ptrVal must be a pointer")
	}
	if goPtrVal.IsNil() {
		panic("ptrVal must not be nil")
	}
	goVal := goPtrVal.Elem()
	if schemaType == nil {
		schemaType = inferSchema(goVal.Type())
	}
	return &_node{val: goVal, schemaType: schemaType}
}

// Unwrap takes an datamodel.Node implemented by Prototype or Wrap,
// and returns a pointer to the inner Go value.
//
// Unwrap returns nil if the node isn't implemented by this package.
func Unwrap(node datamodel.Node) (ptr interface{}) {
	var val reflect.Value
	switch node := node.(type) {
	case *_node:
		val = node.val
	case *_nodeRepr:
		val = node.val
	default:
		return nil
	}
	if val.Kind() == reflect.Ptr {
		panic("didn't expect val to be a pointer")
	}
	return val.Addr().Interface()
}
