// Package bindnode provides a datamodel.Node implementation via Go reflection.
//
// This package is EXPERIMENTAL; its behavior and API might change as it's still
// in development.
package bindnode

import (
	"reflect"

	"github.com/ipfs/go-cid"
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

		if schemaType == nil {
			schemaType = inferSchema(goType, 0)
		} else {
			verifyCompatibility(cfg, make(map[seenEntry]bool), goType, schemaType)
		}
	}

	return &_prototype{cfg: cfg, schemaType: schemaType, goType: goType}
}

// scalar kinds excluding Null

// CustomFromBool is a custom converter function that takes a bool and returns a
// custom type
type CustomFromBool func(bool) (interface{}, error)

// CustomToBool is a custom converter function that takes a custom type and
// returns a bool
type CustomToBool func(interface{}) (bool, error)

// CustomFromInt is a custom converter function that takes an int and returns a
// custom type
type CustomFromInt func(int64) (interface{}, error)

// CustomToInt is a custom converter function that takes a custom type and
// returns an int
type CustomToInt func(interface{}) (int64, error)

// CustomFromFloat is a custom converter function that takes a float and returns
// a custom type
type CustomFromFloat func(float64) (interface{}, error)

// CustomToFloat is a custom converter function that takes a custom type and
// returns a float
type CustomToFloat func(interface{}) (float64, error)

// CustomFromString is a custom converter function that takes a string and
// returns custom type
type CustomFromString func(string) (interface{}, error)

// CustomToString is a custom converter function that takes a custom type and
// returns a string
type CustomToString func(interface{}) (string, error)

// CustomFromBytes is a custom converter function that takes a byte slice and
// returns a custom type
type CustomFromBytes func([]byte) (interface{}, error)

// CustomToBytes is a custom converter function that takes a custom type and
// returns a byte slice
type CustomToBytes func(interface{}) ([]byte, error)

// CustomFromLink is a custom converter function that takes a cid.Cid and
// returns a custom type
type CustomFromLink func(cid.Cid) (interface{}, error)

// CustomToLink is a custom converter function that takes a custom type and
// returns a cid.Cid
type CustomToLink func(interface{}) (cid.Cid, error)

// CustomFromAny is a custom converter function that takes a datamodel.Node and
// returns a custom type
type CustomFromAny func(datamodel.Node) (interface{}, error)

// CustomToAny is a custom converter function that takes a custom type and
// returns a datamodel.Node
type CustomToAny func(interface{}) (datamodel.Node, error)

type converter struct {
	kind schema.TypeKind

	customFromBool CustomFromBool
	customToBool   CustomToBool

	customFromInt CustomFromInt
	customToInt   CustomToInt

	customFromFloat CustomFromFloat
	customToFloat   CustomToFloat

	customFromString CustomFromString
	customToString   CustomToString

	customFromBytes CustomFromBytes
	customToBytes   CustomToBytes

	customFromLink CustomFromLink
	customToLink   CustomToLink

	customFromAny CustomFromAny
	customToAny   CustomToAny
}

type config map[reflect.Type]converter

func (c config) converterFor(val reflect.Value) (converter, bool) {
	if len(c) == 0 {
		return converter{}, false
	}
	conv, ok := c[nonPtrType(val)]
	return conv, ok
}

func (c config) converterForType(typ reflect.Type) (converter, bool) {
	if len(c) == 0 {
		return converter{}, false
	}
	conv, ok := c[typ]
	return conv, ok
}

// Option is able to apply custom options to the bindnode API
type Option func(config)

// AddCustomTypeBoolConverter adds custom converter functions for a particular
// type as identified by a pointer in the first argument.
// The fromFunc is of the form: func(bool) (interface{}, error)
// and toFunc is of the form: func(interface{}) (bool, error)
// where interface{} is a pointer form of the type we are converting.
//
// AddCustomTypeBoolConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeBoolConverter(ptrVal interface{}, from CustomFromBool, to CustomToBool) Option {
	customType := nonPtrType(reflect.ValueOf(ptrVal))
	return func(cfg config) {
		cfg[customType] = converter{
			kind:           schema.TypeKind_Bool,
			customFromBool: from,
			customToBool:   to,
		}
	}
}

// AddCustomTypeIntConverter adds custom converter functions for a particular
// type as identified by a pointer in the first argument.
// The fromFunc is of the form: func(int64) (interface{}, error)
// and toFunc is of the form: func(interface{}) (int64, error)
// where interface{} is a pointer form of the type we are converting.
//
// AddCustomTypeIntConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeIntConverter(ptrVal interface{}, from CustomFromInt, to CustomToInt) Option {
	customType := nonPtrType(reflect.ValueOf(ptrVal))
	return func(cfg config) {
		cfg[customType] = converter{
			kind:          schema.TypeKind_Int,
			customFromInt: from,
			customToInt:   to,
		}
	}
}

// AddCustomTypeFloatConverter adds custom converter functions for a particular
// type as identified by a pointer in the first argument.
// The fromFunc is of the form: func(float64) (interface{}, error)
// and toFunc is of the form: func(interface{}) (float64, error)
// where interface{} is a pointer form of the type we are converting.
//
// AddCustomTypeFloatConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeFloatConverter(ptrVal interface{}, from CustomFromFloat, to CustomToFloat) Option {
	customType := nonPtrType(reflect.ValueOf(ptrVal))
	return func(cfg config) {
		cfg[customType] = converter{
			kind:            schema.TypeKind_Float,
			customFromFloat: from,
			customToFloat:   to,
		}
	}
}

// AddCustomTypeStringConverter adds custom converter functions for a particular
// type as identified by a pointer in the first argument.
// The fromFunc is of the form: func(string) (interface{}, error)
// and toFunc is of the form: func(interface{}) (string, error)
// where interface{} is a pointer form of the type we are converting.
//
// AddCustomTypeStringConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeStringConverter(ptrVal interface{}, from CustomFromString, to CustomToString) Option {
	customType := nonPtrType(reflect.ValueOf(ptrVal))
	return func(cfg config) {
		cfg[customType] = converter{
			kind:             schema.TypeKind_String,
			customFromString: from,
			customToString:   to,
		}
	}
}

// AddCustomTypeBytesConverter adds custom converter functions for a particular
// type as identified by a pointer in the first argument.
// The fromFunc is of the form: func([]byte) (interface{}, error)
// and toFunc is of the form: func(interface{}) ([]byte, error)
// where interface{} is a pointer form of the type we are converting.
//
// AddCustomTypeBytesConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeBytesConverter(ptrVal interface{}, from CustomFromBytes, to CustomToBytes) Option {
	customType := nonPtrType(reflect.ValueOf(ptrVal))
	return func(cfg config) {
		cfg[customType] = converter{
			kind:            schema.TypeKind_Bytes,
			customFromBytes: from,
			customToBytes:   to,
		}
	}
}

// AddCustomTypeLinkConverter adds custom converter functions for a particular
// type as identified by a pointer in the first argument.
// The fromFunc is of the form: func([]byte) (interface{}, error)
// and toFunc is of the form: func(interface{}) ([]byte, error)
// where interface{} is a pointer form of the type we are converting.
//
// Beware that this API is only compatible with cidlink.Link types in the data
// model and may result in errors if attempting to convert from other
// datamodel.Link types.
//
// AddCustomTypeLinkConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeLinkConverter(ptrVal interface{}, from CustomFromLink, to CustomToLink) Option {
	customType := nonPtrType(reflect.ValueOf(ptrVal))
	return func(cfg config) {
		cfg[customType] = converter{
			kind:           schema.TypeKind_Link,
			customFromLink: from,
			customToLink:   to,
		}
	}
}

// AddCustomTypeAnyConverter adds custom converter functions for a particular
// type as identified by a pointer in the first argument.
// The fromFunc is of the form: func(datamodel.Node) (interface{}, error)
// and toFunc is of the form: func(interface{}) (datamodel.Node, error)
// where interface{} is a pointer form of the type we are converting.
//
// This method should be able to deal with all forms of Any and return an error
// if the expected data forms don't match the expected.
//
// AddCustomTypeAnyConverter is an EXPERIMENTAL API and may be removed or
// changed in a future release.
func AddCustomTypeAnyConverter(ptrVal interface{}, from CustomFromAny, to CustomToAny) Option {
	customType := nonPtrType(reflect.ValueOf(ptrVal))
	return func(cfg config) {
		cfg[customType] = converter{
			kind:          schema.TypeKind_Any,
			customFromAny: from,
			customToAny:   to,
		}
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
