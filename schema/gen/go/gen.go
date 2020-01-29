package gengo

import (
	"fmt"
	"io"
)

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

type typedLinkNodeGenerator interface {
	// all methods in typedNodeGenerator
	typedNodeGenerator

	// as typed.LinkNode.ReferencedNodeBuilder generator
	EmitTypedLinkNodeMethodReferencedLinkBuilder(io.Writer)
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

// EmitFileHeader emits a baseline package header that will
// allow a file with a generated type to compile
func EmitFileHeader(packageName string, w io.Writer) {
	fmt.Fprintf(w, "package %s\n\n", packageName)
	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "\tipld \"github.com/ipld/go-ipld-prime\"\n")
	fmt.Fprintf(w, "\t\"github.com/ipld/go-ipld-prime/impl/typed\"\n")
	fmt.Fprintf(w, "\t\"github.com/ipld/go-ipld-prime/schema\"\n")
	fmt.Fprintf(w, ")\n\n")
	fmt.Fprintf(w, "// Code generated go-ipld-prime DO NOT EDIT.\n\n")
}

// EmitEntireType outputs every possible type of code generation for a
// typedNodeGenerator
func EmitEntireType(tg typedNodeGenerator, w io.Writer) {
	tg.EmitNativeType(w)
	tg.EmitNativeAccessors(w)
	tg.EmitNativeBuilder(w)
	tg.EmitNativeMaybe(w)
	tg.EmitNodeType(w)
	tg.EmitTypedNodeMethodType(w)
	tg.EmitNodeMethodReprKind(w)
	tg.EmitNodeMethodLookupString(w)
	tg.EmitNodeMethodLookup(w)
	tg.EmitNodeMethodLookupIndex(w)
	tg.EmitNodeMethodLookupSegment(w)
	tg.EmitNodeMethodMapIterator(w)
	tg.EmitNodeMethodListIterator(w)
	tg.EmitNodeMethodLength(w)
	tg.EmitNodeMethodIsUndefined(w)
	tg.EmitNodeMethodIsNull(w)
	tg.EmitNodeMethodAsBool(w)
	tg.EmitNodeMethodAsInt(w)
	tg.EmitNodeMethodAsFloat(w)
	tg.EmitNodeMethodAsString(w)
	tg.EmitNodeMethodAsBytes(w)
	tg.EmitNodeMethodAsLink(w)

	tg.EmitNodeMethodNodeBuilder(w)
	tnbg := tg.GetNodeBuilderGen()
	tnbg.EmitNodebuilderType(w)
	tnbg.EmitNodebuilderConstructor(w)
	tnbg.EmitNodebuilderMethodCreateMap(w)
	tnbg.EmitNodebuilderMethodAmendMap(w)
	tnbg.EmitNodebuilderMethodCreateList(w)
	tnbg.EmitNodebuilderMethodAmendList(w)
	tnbg.EmitNodebuilderMethodCreateNull(w)
	tnbg.EmitNodebuilderMethodCreateBool(w)
	tnbg.EmitNodebuilderMethodCreateInt(w)
	tnbg.EmitNodebuilderMethodCreateFloat(w)
	tnbg.EmitNodebuilderMethodCreateString(w)
	tnbg.EmitNodebuilderMethodCreateBytes(w)
	tnbg.EmitNodebuilderMethodCreateLink(w)

	tlg, ok := tg.(typedLinkNodeGenerator)
	if ok {
		tlg.EmitTypedLinkNodeMethodReferencedLinkBuilder(w)
	}

	tg.EmitTypedNodeMethodRepresentation(w)
	rng := tg.GetRepresentationNodeGen()
	if rng == nil { // FIXME: hack to save me from stubbing tons right now, remove when done
		return
	}
	rng.EmitNodeType(w)
	rng.EmitNodeMethodReprKind(w)
	rng.EmitNodeMethodLookupString(w)
	rng.EmitNodeMethodLookup(w)
	rng.EmitNodeMethodLookupIndex(w)
	rng.EmitNodeMethodLookupSegment(w)
	rng.EmitNodeMethodMapIterator(w)
	rng.EmitNodeMethodListIterator(w)
	rng.EmitNodeMethodLength(w)
	rng.EmitNodeMethodIsUndefined(w)
	rng.EmitNodeMethodIsNull(w)
	rng.EmitNodeMethodAsBool(w)
	rng.EmitNodeMethodAsInt(w)
	rng.EmitNodeMethodAsFloat(w)
	rng.EmitNodeMethodAsString(w)
	rng.EmitNodeMethodAsBytes(w)
	rng.EmitNodeMethodAsLink(w)

	rng.EmitNodeMethodNodeBuilder(w)
	rnbg := rng.GetNodeBuilderGen()
	rnbg.EmitNodebuilderType(w)
	rnbg.EmitNodebuilderConstructor(w)
	rnbg.EmitNodebuilderMethodCreateMap(w)
	rnbg.EmitNodebuilderMethodAmendMap(w)
	rnbg.EmitNodebuilderMethodCreateList(w)
	rnbg.EmitNodebuilderMethodAmendList(w)
	rnbg.EmitNodebuilderMethodCreateNull(w)
	rnbg.EmitNodebuilderMethodCreateBool(w)
	rnbg.EmitNodebuilderMethodCreateInt(w)
	rnbg.EmitNodebuilderMethodCreateFloat(w)
	rnbg.EmitNodebuilderMethodCreateString(w)
	rnbg.EmitNodebuilderMethodCreateBytes(w)
	rnbg.EmitNodebuilderMethodCreateLink(w)
}
