package ipld

import (
	"fmt"
)

// ErrWrongKind may be returned from functions on the Node interface when
// a method is invoked which doesn't make sense for the Kind and/or ReprKind
// that node concretely contains.
//
// For example, calling AsString on a map will return ErrWrongKind.
// Calling TraverseField on an int will similarly return ErrWrongKind.
type ErrWrongKind struct {
	// MethodName is literally the string for the operation attempted, e.g.
	// "AsString".
	MethodName string

	// ApprorpriateKind is used to describe the Kind which the erroring method
	// would make sense for.
	//
	// In the case of typed nodes, this will typically refer to the 'natural'
	// data-model kind for such a type (e.g., structs will say 'map' here).
	AppropriateKind ReprKind

	ActualKind ReprKind // FIXME okay just no, this really needs to say what it knows.

	// REVIEW this is almost certainly wrong.  Maybe you need some short enums
	//  for things like "the kinds you can traverse by name" e.g. map+struct.
	//  And I'm really sparse for reasons not to put schema-level Kind in the root package.
	//   All the counter-arguments are similar to the ones about Path, and they
	//    just flat out lose when up against "errors matter".
}

func (e ErrWrongKind) Error() string {
	return fmt.Sprintf("func called on wrong kind: %s called on a %s node, but only makes sense on %s", e.MethodName, e.ActualKind, e.AppropriateKind)
}
