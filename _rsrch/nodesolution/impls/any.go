package impls

import (
	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

var (
	//_ ipld.Node          = &anyNode{}
	_ ipld.NodeStyle   = Style__Any{}
	_ ipld.NodeBuilder = &anyBuilder{}
	//_ ipld.NodeAssembler = &anyAssembler{}
)

// anyNode is a union meant for alloc amortization; see anyAssembler.
// Note that anyBuilder doesn't use anyNode, because it's not aiming to amortize anything.
//
// REVIEW: if there's any point in keeping this around.  It's here for completeness,
// but not currently used anywhere in package, and also not currently exported.
type anyNode struct {
	plainMap
	plainString
	plainInt
	// TODO: more of these embeds for the remaining kinds.
}

// -- Node interface methods -->

// Unimplemented at present -- see "REVIEW" comment on anyNode.

// -- NodeStyle -->

type Style__Any struct{}

func (Style__Any) NewBuilder() ipld.NodeBuilder {
	return &anyBuilder{}
}

// -- NodeBuilder -->

// anyBuilder is a builder for any kind of node.
//
// anyBuilder is a little unusual in its internal workings:
// unlike most builders, it doesn't embed the corresponding assembler,
// nor will it end up using anyNode,
// but instead embeds a builder for each of the kinds it might contain.
// This is because we want a more granular return at the end:
// if we used anyNode, and returned a pointer to just the relevant part of it,
// we'd have all the extra bytes of anyNode still reachable in GC terms
// for as long as that handle to the interior of it remains live.
type anyBuilder struct {
	kind ipld.ReprKind // used to select which builder to delegate 'Build' to!  Set on first interaction.

	// Only one of the following ends up being used...
	//  but we don't know in advance which one, so all are embeded here.
	//   This uses excessive space, but amortizes allocations, and all will be
	//    freed as soon as the builder is done.

	mapBuilder    plainMap__Builder
	stringBuilder plainString__Builder
	intBuilder    plainInt__Builder
	// TODO: more of these embeds for the remaining kinds.
}

func (nb *anyBuilder) Reset() {
	*nb = anyBuilder{}
}

func (nb *anyBuilder) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error) {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Map
	return nb.mapBuilder.BeginMap(sizeHint)
}
func (nb *anyBuilder) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) {
	panic("soon")
}
func (nb *anyBuilder) AssignNull() error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Null
	return nil
}
func (nb *anyBuilder) AssignBool(v bool) error {
	panic("soon")
}
func (nb *anyBuilder) AssignInt(v int) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Int
	return nb.intBuilder.AssignInt(v)
}
func (nb *anyBuilder) AssignFloat(v float64) error {
	panic("soon")
}
func (nb *anyBuilder) AssignString(v string) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_String
	return nb.stringBuilder.AssignString(v)
}
func (nb *anyBuilder) AssignBytes(v []byte) error {
	panic("soon")
}
func (nb *anyBuilder) AssignLink(v ipld.Link) error {
	panic("soon")
}
func (nb *anyBuilder) AssignNode(v ipld.Node) error {
	// TODO what to do here?  should we just... keep it, in another `Node` field?
	panic("soon")
}
func (anyBuilder) Style() ipld.NodeStyle {
	return Style__Any{}
}

func (nb *anyBuilder) Build() ipld.Node {
	switch nb.kind {
	case ipld.ReprKind_Invalid:
		panic("misuse")
	case ipld.ReprKind_Map:
		return nb.mapBuilder.Build()
	case ipld.ReprKind_List:
		panic("soon")
	case ipld.ReprKind_Null:
		//return ipld.Null
		panic("soon")
	case ipld.ReprKind_Bool:
		panic("soon")
	case ipld.ReprKind_Int:
		return nb.intBuilder.Build()
	case ipld.ReprKind_Float:
		panic("soon")
	case ipld.ReprKind_String:
		return nb.stringBuilder.Build()
	case ipld.ReprKind_Bytes:
		panic("soon")
	case ipld.ReprKind_Link:
		panic("soon")
	default:
		panic("unreachable")
	}
}

// -- NodeAssembler -->

// ... oddly enough, we seem to be able to put off implementing this
//  until we also implement something that goes full-hog on amortization
//   and actually has a slab of `anyNode`.  Which so far, nothing does.
//    See "REVIEW" comment on anyNode.
type anyAssembler struct {
	w *anyNode
}

// -- Additional typedefs for maintaining 'any' style property -->

type anyInhabitedByString plainString

func (anyInhabitedByString) Style() ipld.NodeStyle {
	return Style__Any{}
}

type anyInhabitedByInt plainInt

func (anyInhabitedByInt) Style() ipld.NodeStyle {
	return Style__Any{}
}

// TODO: more of these typedefs for the remaining kinds.
