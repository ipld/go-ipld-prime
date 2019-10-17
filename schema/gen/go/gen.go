package gengo

import (
	"fmt"
	"io"
)

// there are roughly *seven* categories of API to generate per type:
// - 1: the readonly thing a native caller uses
// - 2: the builder thing a native caller uses
// - 3: the readonly typed node
// - 4: the builder for typed node
// - 5: the readonly representation node
// - 6: the builder via representation
// - 7: and a maybe wrapper
//
// 1, 2, 3, and 7 are emitted from `typedNodeGenerator` (3 from the embedded `nodeGenerator`).
// 5 is emitted from another `nodeGenerator` instance.
// 4 and 6 are emitted from two distinct `nodebuilderGenerator` instances.

// file patterns in this package:
//
// - for each kind:
//   - `genKind{Kind}.go` -- has emitters for the native type parts (1, 2, 7).
//   - `genKind{Kind}Node.go` -- has emitters for the typed node parts (3, 4), and the entrypoint to (5).
//   - for each representation that kind can have:
//      - `genKind{Kind}Repr{ReprStrat}.go` -- has emitters for (5, 6).

// typedNodeGenerator declares a standard names for a bunch of methods for generating
// code for our schema types.  There's still numerous places where other casts
// to more specific interfaces will be required (so, technically, it's not a
// very powerful interface; it's not so much that the abstractions leak as that
// the floodgates are outright open), but this at least forces consistency onto
// the parts where we can.
//
// All Emit{foo} methods should emit one trailing and one leading linebreak, or,
// nothing (e.g. string kinds don't need to produce a dummy map iterator, so
// such a method can just emit nothing, and the extra spacing between sections
// shouldn't accumulate).
//
// None of these methods return error values because we panic in this package.
//
type typedNodeGenerator interface {

	// -- the natively-typed apis -->
	//   (might be more readable to group these in another interface and have it
	//     return a `typedNodeGenerator` with the rest?  but structurally same.)

	EmitNativeType(io.Writer)
	EmitNativeAccessors(io.Writer) // depends on the kind -- field accessors for struct, typed iterators for map, etc.
	EmitNativeBuilder(io.Writer)   // typically emits some kind of struct that has a Build method.
	EmitNativeMaybe(io.Writer)     // a pointer-free 'maybe' mechanism is generated for all types.

	// -- the typed.Node.Type method and vars -->

	EmitTypedNodeMethodType(io.Writer) // these emit dummies for now

	// -- all node methods -->
	//   (and note that the nodeBuilder for this one should be the "semantic" one,
	//     e.g. it *always* acts like a map for structs, even if the repr is different.)

	nodeGenerator

	// -- and the representation and its node and nodebuilder -->

	EmitTypedNodeMethodRepresentation(io.Writer)
	GetRepresentationNodeGen() nodeGenerator // includes transitively the matched nodebuilderGenerator
}

type nodeGenerator interface {
	EmitNodeType(io.Writer)
	EmitNodeMethodReprKind(io.Writer)
	EmitNodeMethodLookupString(io.Writer)
	EmitNodeMethodLookup(io.Writer)
	EmitNodeMethodLookupIndex(io.Writer)
	EmitNodeMethodLookupSegment(io.Writer)
	EmitNodeMethodMapIterator(io.Writer)  // also iterator itself
	EmitNodeMethodListIterator(io.Writer) // also iterator itself
	EmitNodeMethodLength(io.Writer)
	EmitNodeMethodIsUndefined(io.Writer)
	EmitNodeMethodIsNull(io.Writer)
	EmitNodeMethodAsBool(io.Writer)
	EmitNodeMethodAsInt(io.Writer)
	EmitNodeMethodAsFloat(io.Writer)
	EmitNodeMethodAsString(io.Writer)
	EmitNodeMethodAsBytes(io.Writer)
	EmitNodeMethodAsLink(io.Writer)
	EmitNodeMethodNodeBuilder(io.Writer)

	GetNodeBuilderGen() nodebuilderGenerator
}

type nodebuilderGenerator interface {
	EmitNodebuilderType(io.Writer)
	EmitNodebuilderConstructor(io.Writer)

	EmitNodebuilderMethodCreateMap(io.Writer) // also mapbuilder itself
	EmitNodebuilderMethodAmendMap(io.Writer)  // also listbuilder itself
	EmitNodebuilderMethodCreateList(io.Writer)
	EmitNodebuilderMethodAmendList(io.Writer)
	EmitNodebuilderMethodCreateNull(io.Writer)
	EmitNodebuilderMethodCreateBool(io.Writer)
	EmitNodebuilderMethodCreateInt(io.Writer)
	EmitNodebuilderMethodCreateFloat(io.Writer)
	EmitNodebuilderMethodCreateString(io.Writer)
	EmitNodebuilderMethodCreateBytes(io.Writer)
	EmitNodebuilderMethodCreateLink(io.Writer)
}

func emitFileHeader(w io.Writer) {
	fmt.Fprintf(w, "package whee\n\n")
	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "\tipld \"github.com/ipld/go-ipld-prime\"\n")
	fmt.Fprintf(w, "\t\"github.com/ipld/go-ipld-prime/impl/typed\"\n")
	fmt.Fprintf(w, "\t\"github.com/ipld/go-ipld-prime/schema\"\n")
	fmt.Fprintf(w, ")\n\n")
}
