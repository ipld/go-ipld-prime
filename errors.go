package ipld

import (
	"fmt"
)

// ErrWrongKind may be returned from functions on the Node interface when
// a method is invoked which doesn't make sense for the Kind and/or ReprKind
// that node concretely contains.
//
// For example, calling AsString on a map will return ErrWrongKind.
// Calling Lookup on an int will similarly return ErrWrongKind.
type ErrWrongKind struct {
	// TypeName may optionally indicate the named type of a node the function
	// was called on (if the node was typed!), or, may be the empty string.
	TypeName string

	// MethodName is literally the string for the operation attempted, e.g.
	// "AsString".
	//
	// For methods on nodebuilders, we say e.g. "NodeBuilder.CreateMap".
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
	if e.TypeName == "" {
		return fmt.Sprintf("func called on wrong kind: %s called on a %s node, but only makes sense on %s", e.MethodName, e.ActualKind, e.AppropriateKind)
	} else {
		return fmt.Sprintf("func called on wrong kind: %s called on a %s node (kind: %s), but only makes sense on %s", e.MethodName, e.TypeName, e.ActualKind, e.AppropriateKind)
	}
}

// ErrNotExists may be returned from the lookup functions of the Node interface
// to indicate a missing value.
//
// Note that typed.ErrNoSuchField is another type of error which sometimes
// occurs in similar places as ErrNotExists.  ErrNoSuchField is preferred
// when handling data with constraints provided by a schema that mean that
// a field can *never* exist (as differentiated from a map key which is
// simply absent in some data).
type ErrNotExists struct {
	Segment PathSegment
}

func (e ErrNotExists) Error() string {
	return fmt.Sprintf("key not found: %q", e.Segment)
}

// ErrInvalidKey may be returned from lookup functions on the Node interface
// when a key is invalid.
//
// Common examples of this are when `Lookup(Node)` is used with a non-string Node;
// typed nodes also introduce other reasons a key may be invalid.
type ErrInvalidKey struct {
	Reason string

	// Perhaps typed.ErrNoSuchField could be folded into this?
	// Perhaps Reason could be replaced by an enum of "NoSuchField"|"NotAString"|"ConstraintRejected"?
	// Might be hard to get rid of the freetext field entirely -- constraints may be nontrivial to describe.
}

func (e ErrInvalidKey) Error() string {
	return fmt.Sprintf("invalid key: %s", e.Reason)
}

// ErrIteratorOverread is returned when calling 'Next' on a MapIterator or
// ListIterator when it is already done.
type ErrIteratorOverread struct{}

func (e ErrIteratorOverread) Error() string {
	return "iterator overread"
}
