package gengo

import (
	"fmt"
	"io"
)

// Not yet handled:
//   - the Style access singleton (and this one requires some aggregated input, so may trigger some refactoring).
//   - reprBuilder

// TypeGenerator gathers all the info for generating all code related to one
// type in the schema.
type TypeGenerator interface {
	// -- the natively-typed apis -->

	EmitNativeType(io.Writer)
	EmitNativeAccessors(io.Writer) // depends on the kind -- field accessors for struct, typed iterators for map, etc.
	EmitNativeBuilder(io.Writer)   // typically emits some kind of struct that has a Build method.
	EmitNativeMaybe(io.Writer)     // a pointer-free 'maybe' mechanism is generated for all types.

	// -- the schema.TypedNode.Type method and vars -->

	EmitTypeConst(io.Writer)           // these emit dummies for now
	EmitTypedNodeMethodType(io.Writer) // these emit dummies for now

	// -- all node methods -->
	//   (and note that the nodeBuilder for this one should be the "semantic" one,
	//     e.g. it *always* acts like a map for structs, even if the repr is different.)

	NodeGenerator

	// -- and the representation and its node and nodebuilder -->
	//    (these vary!)

	EmitTypedNodeMethodRepresentation(io.Writer)
	GetRepresentationNodeGen() NodeGenerator // includes transitively the matched nodebuilderGenerator
}

type NodeGenerator interface {
	EmitNodeType(io.Writer)           // usually already covered by EmitNativeType for the primary node, but has a nonzero body for the repr node
	EmitNodeTypeAssertions(io.Writer) // optional to include this content
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
	EmitNodeMethodStyle(io.Writer)
	EmitNodeStyleType(io.Writer)
	EmitNodeBuilder(io.Writer)   // the whole thing
	EmitNodeAssembler(io.Writer) // the whole thing
}

// EmitFileHeader emits a baseline package header that will
// allow a file with a generated type to compile.
// (Fortunately, there are no variations in this.)
func EmitFileHeader(packageName string, w io.Writer) {
	fmt.Fprintf(w, "package %s\n\n", packageName)
	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "\tipld \"github.com/ipld/go-ipld-prime\"\n")
	fmt.Fprintf(w, "\t\"github.com/ipld/go-ipld-prime/node/mixins\"\n")
	fmt.Fprintf(w, "\t\"github.com/ipld/go-ipld-prime/schema\"\n")
	fmt.Fprintf(w, ")\n\n")
	fmt.Fprintf(w, "// Code generated go-ipld-prime DO NOT EDIT.\n\n")
}

// EmitEntireType calls all methods of TypeGenerator and streams
// all results into a single writer.
func EmitEntireType(tg TypeGenerator, w io.Writer) {
	tg.EmitNativeType(w)
	tg.EmitNativeAccessors(w)
	tg.EmitNativeBuilder(w)
	tg.EmitNativeMaybe(w)
	EmitNode(tg, w)
	tg.EmitTypedNodeMethodType(w)
	tg.EmitTypedNodeMethodRepresentation(w)

	rng := tg.GetRepresentationNodeGen()
	if rng == nil { // FIXME: hack to save me from stubbing tons right now, remove when done
		return
	}
	EmitNode(rng, w)
}

func EmitNode(ng NodeGenerator, w io.Writer) {
	ng.EmitNodeType(w)
	ng.EmitNodeTypeAssertions(w)
	ng.EmitNodeMethodReprKind(w)
	ng.EmitNodeMethodLookupString(w)
	ng.EmitNodeMethodLookup(w)
	ng.EmitNodeMethodLookupIndex(w)
	ng.EmitNodeMethodLookupSegment(w)
	ng.EmitNodeMethodMapIterator(w)
	ng.EmitNodeMethodListIterator(w)
	ng.EmitNodeMethodLength(w)
	ng.EmitNodeMethodIsUndefined(w)
	ng.EmitNodeMethodIsNull(w)
	ng.EmitNodeMethodAsBool(w)
	ng.EmitNodeMethodAsInt(w)
	ng.EmitNodeMethodAsFloat(w)
	ng.EmitNodeMethodAsString(w)
	ng.EmitNodeMethodAsBytes(w)
	ng.EmitNodeMethodAsLink(w)
	ng.EmitNodeMethodStyle(w)

	ng.EmitNodeStyleType(w)

	ng.EmitNodeBuilder(w)
	ng.EmitNodeAssembler(w)
}