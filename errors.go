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
	// CONSIDER: if we should add a `TypeName string` here as well?
	// It seems to be useful information, and in many places we've shoved it
	// along with the MethodName; but while that's fine for the printed message,
	// we could do better internally (and it would enable `typed.wrapnode*` to
	// touch this field on its way out in a nice reusable way).

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

// ErrNotExists may be returned from the traversal functions of the Node
// interface to indicate a missing value.
//
// Note that typed.ErrNoSuchField is another type of error which sometimes
// occurs in similar places as ErrNotExists.  ErrNoSuchField is preferred
// when handling data with constraints provided by a schema that mean that
// a field can *never* exist (as differentiated from a map key which is
// simply absent in some data).
type ErrNotExists struct {
	Segment string // REVIEW: might be better to use PathSegment, but depends on another refactor.
}

func (e ErrNotExists) Error() string {
	return fmt.Sprintf("key not found: %q", e.Segment)
}

// ErrIteratorOverread is returned when calling 'Next' on a MapIterator or
// ListIterator when it is already done.
type ErrIteratorOverread struct{}

func (e ErrIteratorOverread) Error() string {
	return "iterator overread"
}
