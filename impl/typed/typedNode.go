package typed

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

// typed.Node is a superset of the ipld.Node interface, and has additional behaviors.
//
// A typed.Node can be inspected for its schema.Type and schema.Kind,
// which conveys much more and richer information than the Data Model layer
// ipld.ReprKind.
//
// There are many different implementations of typed.Node.
// One implementation can wrap any other existing ipld.Node (i.e., it's zero-copy)
// and promises that it has *already* been validated to match the typesystem.Type;
// another implementation similarly wraps any other existing ipld.Node, but
// defers to the typesystem validation checking to fields that are accessed;
// and when using code generation tools, all of the generated native Golang
// types produced by the codegen will each individually implement typed.Node.
//
// Note that typed.Node can wrap *other* typed.Node instances.
// Imagine you have two parts of a very large code base which have codegen'd
// components which are from different versions of a schema.  Smooth migrations
// and zero-copy type-safe data sharing between them: We can accommodate that!
//
// Typed nodes sometimes have slightly different behaviors than plain nodes:
// For example, when looking up fields on a typed node that's a struct,
// the error returned for a lookup with a key that's not a field name will
// be ErrNoSuchField (instead of ErrNotExists).
// These behaviors apply to the typed.Node only and not their representations;
// continuing the example, the .Representation().LookupString() method on
// that same node for the same key as plain `.LookupString()` will still
// return ErrNotExists, because the representation isn't a typed.Node!
type Node interface {
	// typed.Node acts just like a regular Node for almost all purposes;
	// which ReprKind it acts as is determined by the TypeKind.
	// (Note that the representation strategy of the type does *not* affect
	// the ReprKind of typed.Node -- rather, the representation strategy
	// affects the `.Representation().ReprKind()`.)
	//
	// For example: if the `.Type().Kind()` of this node is "struct",
	// it will act like ReprKind() == "map"
	// (even if Type().(Struct).ReprStrategy() is "tuple").
	ipld.Node

	// Type returns a reference to the reified schema.Type value.
	Type() schema.Type

	// Representation returns an ipld.Node which sees the data in this node
	// in its representation form.
	//
	// For example: if the `.Type().Kind()` of this node is "struct",
	// `.Representation().Kind()` may vary based on its representation strategy:
	// if the representation strategy is "map", then it will be ReprKind=="map";
	// if the streatgy is "tuple", then it will be ReprKind=="list".
	Representation() ipld.Node
}

// unboxing is... ugh, we probably should codegen an unbox method per concrete type.
//  (or, attach them to the non-pointer type, which would namespace in an alloc-free way, but i don't know if that's anything but confusing.)
//  there are notes about this from way back at 2019.01; reread to see if any remain relevant and valid.
// main important point is: it's not gonna be casting.
//  if casting was sufficient to unbox, it'd mean every method on the Node interface would be difficult to use as a field name on a struct type.  undesirable.
//   okay, or, alternative, we flip this to `superapi.Footype{}.Fields().FrobFieldName()`.  that strikes me as unlikely to be pleasing, though.
//    istm we can safely expect direct use of field names much, much more often that flipping back and forth to hypergeneric node; so we should optimize syntax for that accordingly.
