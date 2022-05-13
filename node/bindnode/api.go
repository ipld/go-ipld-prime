// Package bindnode provides a datamodel.Node implementation via Go reflection.
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
// IPLD schema type. Note that the result is also a datamodel.NodePrototype.
//
// If both the Go type and schema type are supplied, it is assumed that they are
// compatible with one another.
//
// If either the Go type or schema type are nil, we infer the missing type from
// the other provided type. For example, we can infer an unnamed Go struct type
// for a schema struct type, and we can infer a schema Int type for a Go int64
// type. The inferring logic is still a work in progress and subject to change.
// At this time, inferring IPLD Unions and Enums from Go types is not supported.
//
// When supplying a non-nil ptrType, Prototype only obtains the Go pointer type
// from it, so its underlying value will typically be nil. For example:
//
//     proto := bindnode.Prototype((*goType)(nil), schemaType)
func Prototype(ptrType interface{}, schemaType schema.Type, options ...Option) schema.TypedPrototype {
	if ptrType == nil && schemaType == nil {
		panic("bindnode: either ptrType or schemaType must not be nil")
	}

	cfg := applyOptions(options...)

	// TODO: if both are supplied, verify that they are compatible

	var goType reflect.Type
	if ptrType == nil {
		goType = inferGoType(schemaType, make(map[schema.TypeName]inferredStatus), 0)
	} else {
		goPtrType := reflect.TypeOf(ptrType)
		if goPtrType.Kind() != reflect.Ptr {
			panic("bindnode: ptrType must be a pointer")
		}
		goType = goPtrType.Elem()
		if goType.Kind() == reflect.Ptr {
			panic("bindnode: ptrType must not be a pointer to a pointer")
		}

		if schemaType == nil {
			schemaType = inferSchema(goType, 0)
		} else {
			verifyCompatibility(cfg, make(map[seenEntry]bool), goType, schemaType)
		}
	}

	return &_prototype{cfg: cfg, schemaType: schemaType, goType: goType}
}

type converter struct {
	kindType reflect.Type
	fromFunc reflect.Value // func(Kind) (*Custom, error)
	toFunc   reflect.Value // func(*Custom) (Kind, error)
}

type config map[reflect.Type]converter

// Option is able to apply custom options to the bindnode API
type Option func(config)

var byteSliceType = reflect.TypeOf([]byte{})
var errorType = reflect.TypeOf((*error)(nil)).Elem()

// AddCustomTypeBytesConverter adds custom converter functions for a particular
// type. The fromBytesFunc is of the form: func([]byte) (interface{}, error)
// and toBytesFunc is of the form: func(interface{}) ([]byte, error)
// where interface{} is a pointer form of the type we are converting.
//
// AddCustomTypeBytesConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeBytesConverter(fromBytesFunc interface{}, toBytesFunc interface{}) Option {
	fbfVal := reflect.ValueOf(fromBytesFunc)
	fbfType := fbfVal.Type()
	if fbfType == nil || fbfType.Kind() != reflect.Func || fbfType.IsVariadic() {
		panic("fromBytesFunc must be a (non-varadic) function")
	}
	ni, no := fbfType.NumIn(), fbfType.NumOut()
	if ni != 1 || no != 2 || fbfType.In(0) != byteSliceType || fbfType.Out(1) != errorType {
		panic("fromBytesFunc must be of the form func([]byte) (interface{}, error)")
	}

	tbfVal := reflect.ValueOf(toBytesFunc)
	tbfType := tbfVal.Type()
	if tbfType == nil || tbfType.Kind() != reflect.Func || tbfType.IsVariadic() {
		panic("toBytesFunc must be a (non-varadic) function")
	}
	ni, no = tbfType.NumIn(), tbfType.NumOut()
	if ni != 1 || no != 2 || tbfType.Out(0) != byteSliceType || tbfType.Out(1) != errorType {
		panic("toBytesFunc must be of the form func(interface{}) ([]byte, error)")
	}

	if fbfType.Out(0) != tbfType.In(0) {
		panic("toBytesFunc must be of the form func(interface{}) ([]byte, error)")
	}

	customType := fbfType.Out(0)
	if customType.Kind() == reflect.Ptr {
		customType = customType.Elem()
	}
	return func(cfg config) {
		cfg[customType] = converter{byteSliceType, fbfVal, tbfVal}
	}
}

func applyOptions(opt ...Option) config {
	cfg := make(map[reflect.Type]converter)
	for _, o := range opt {
		o(cfg)
	}
	return cfg
}

// Wrap implements a schema.TypedNode given a non-nil pointer to a Go value and an
// IPLD schema type. Note that the result is also a datamodel.Node.
//
// Wrap is meant to be used when one already has a Go value with data.
// As such, ptrVal must not be nil.
//
// Similar to Prototype, if schemaType is non-nil it is assumed to be compatible
// with the Go type, and otherwise it's inferred from the Go type.
func Wrap(ptrVal interface{}, schemaType schema.Type, options ...Option) schema.TypedNode {
	if ptrVal == nil {
		panic("bindnode: ptrVal must not be nil")
	}
	goPtrVal := reflect.ValueOf(ptrVal)
	if goPtrVal.Kind() != reflect.Ptr {
		panic("bindnode: ptrVal must be a pointer")
	}
	if goPtrVal.IsNil() {
		// Note that this can happen if ptrVal was a typed nil.
		panic("bindnode: ptrVal must not be nil")
	}
	cfg := applyOptions(options...)
	goVal := goPtrVal.Elem()
	if goVal.Kind() == reflect.Ptr {
		panic("bindnode: ptrVal must not be a pointer to a pointer")
	}
	if schemaType == nil {
		schemaType = inferSchema(goVal.Type(), 0)
	} else {
		verifyCompatibility(cfg, make(map[seenEntry]bool), goVal.Type(), schemaType)
	}
	return &_node{cfg: cfg, val: goVal, schemaType: schemaType}
}

// TODO: consider making our own Node interface, like:
//
// type WrappedNode interface {
//     datamodel.Node
//     Unwrap() (ptrVal interface)
// }
//
// Pros: API is easier to understand, harder to mix up with other datamodel.Nodes.
// Cons: One usually only has a datamodel.Node, and type assertions can be weird.

// Unwrap takes a datamodel.Node implemented by Prototype or Wrap,
// and returns a pointer to the inner Go value.
//
// Unwrap returns nil if the node isn't implemented by this package.
func Unwrap(node datamodel.Node) (ptrVal interface{}) {
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
		panic("bindnode: didn't expect val to be a pointer")
	}
	return val.Addr().Interface()
}
