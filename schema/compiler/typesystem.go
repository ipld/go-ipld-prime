package compiler

import (
	"github.com/ipld/go-ipld-prime/schema"
)

type TypeSystem struct {
	// Mind the key type here: TypeReference, not TypeName.
	// The key might be a computed anon "name" which is not actually a valid type name itself.
	types map[schema.TypeReference]schema.Type

	// List of types, retained in the original order they were specified,
	// including only those which are named (not any computed anonymous types).
	// This is kept so we can do any listing in the order the user expects,
	// report any errors during rule validation in the same order as the input, etc.
	list []schema.Type
}
