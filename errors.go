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

	// ApprorpriateKind describes which ReprKinds the erroring method would
	// make sense for.
	AppropriateKind ReprKindSet

	// ActualKind describes the ReprKind of the node the method was called on.
	//
	// In the case of typed nodes, this will typically refer to the 'natural'
	// data-model kind for such a type (e.g., structs will say 'map' here).
	ActualKind ReprKind
}

func (e ErrWrongKind) Error() string {
	return fmt.Sprintf("func called on wrong kind: %s called on a %s node, but only makes sense on %s", e.MethodName, e.ActualKind, e.AppropriateKind)
}
