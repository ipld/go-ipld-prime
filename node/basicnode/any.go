package basicnode

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
)

var (
	//_ datamodel.Node          = &anyNode{}
	_ datamodel.NodePrototype                = Prototype__Any{}
	_ datamodel.NodePrototypeSupportingAmend = Prototype__Any{}
	_ datamodel.NodeBuilder                  = &anyBuilder{}
	//_ datamodel.NodeAssembler = &anyAssembler{}
)

// Note that we don't use a "var _" declaration to assert that Chooser
// implements traversal.LinkTargetNodePrototypeChooser, to keep basicnode's
// dependencies fairly light.

// Chooser implements traversal.LinkTargetNodePrototypeChooser.
//
// It can be used directly when loading links into the "any" prototype,
// or with another chooser layer on top, such as:
//
//	prototypeChooser := dagpb.AddSupportToChooser(basicnode.Chooser)
func Chooser(_ datamodel.Link, _ linking.LinkContext) (datamodel.NodePrototype, error) {
	return Prototype.Any, nil
}

// -- Node interface methods -->

// Unimplemented at present -- see "REVIEW" comment on anyNode.

// -- NodePrototype -->

type Prototype__Any struct{}

func (p Prototype__Any) NewBuilder() datamodel.NodeBuilder {
	return p.AmendingBuilder(nil)
}

// -- NodePrototypeSupportingAmend -->

func (p Prototype__Any) AmendingBuilder(base datamodel.Node) datamodel.NodeAmender {
	ab := &anyBuilder{}
	if base != nil {
		ab.kind = base.Kind()
		if npa, castOk := base.Prototype().(datamodel.NodePrototypeSupportingAmend); castOk {
			ab.amender = npa.AmendingBuilder(base)
		} else {
			// This node could be either scalar or recursive
			ab.baseNode = base
		}
	}
	return ab
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
	// kind is set on first interaction, and used to select which builder to delegate 'Build' to!
	// As soon as it's been set to a value other than zero (being "Invalid"), all other Assign/Begin calls will fail since something is already in progress.
	// May also be set to the magic value '99', which means "i dunno, I'm just carrying another node of unknown prototype".
	kind datamodel.Kind

	// Only one of the following ends up being used...
	//  but we don't know in advance which one, so both are embedded here.
	//   This uses excessive space, but amortizes allocations, and all will be
	//    freed as soon as the builder is done.
	// An amender is only used for amendable nodes, while all non-amendable nodes (both recursives and scalars) are
	//  stored directly.
	// 'baseNode' may also hold another Node of unknown prototype (possibly not even from this package),
	//  in which case this is indicated by 'kind==99'.

	amender  datamodel.NodeAmender
	baseNode datamodel.Node
}

func (nb *anyBuilder) Reset() {
	*nb = anyBuilder{}
}

func (nb *anyBuilder) BeginMap(sizeHint int64) (datamodel.MapAssembler, error) {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_Map
	mapBuilder := Prototype.Map.NewBuilder().(*plainMap__Builder)
	nb.amender = mapBuilder
	return mapBuilder.BeginMap(sizeHint)
}
func (nb *anyBuilder) BeginList(sizeHint int64) (datamodel.ListAssembler, error) {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_List
	listBuilder := Prototype.List.NewBuilder().(*plainList__Builder)
	nb.amender = listBuilder
	return listBuilder.BeginList(sizeHint)
}
func (nb *anyBuilder) AssignNull() error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_Null
	return nil
}
func (nb *anyBuilder) AssignBool(v bool) error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_Bool
	nb.baseNode = NewBool(v)
	return nil
}
func (nb *anyBuilder) AssignInt(v int64) error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_Int
	nb.baseNode = NewInt(v)
	return nil
}
func (nb *anyBuilder) AssignFloat(v float64) error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_Float
	nb.baseNode = NewFloat(v)
	return nil
}
func (nb *anyBuilder) AssignString(v string) error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_String
	nb.baseNode = NewString(v)
	return nil
}
func (nb *anyBuilder) AssignBytes(v []byte) error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_Bytes
	nb.baseNode = NewBytes(v)
	return nil
}
func (nb *anyBuilder) AssignLink(v datamodel.Link) error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = datamodel.Kind_Link
	nb.baseNode = NewLink(v)
	return nil
}
func (nb *anyBuilder) AssignNode(v datamodel.Node) error {
	if nb.kind != datamodel.Kind_Invalid {
		panic("misuse")
	}
	nb.kind = 99
	nb.baseNode = v
	return nil
}
func (anyBuilder) Prototype() datamodel.NodePrototype {
	return Prototype.Any
}

func (nb *anyBuilder) Build() datamodel.Node {
	if nb.amender != nil {
		return nb.amender.Build()
	}
	switch nb.kind {
	case datamodel.Kind_Invalid:
		panic("misuse")
	case datamodel.Kind_Map:
		return nb.baseNode
	case datamodel.Kind_List:
		return nb.baseNode
	case datamodel.Kind_Null:
		return datamodel.Null
	case datamodel.Kind_Bool:
		return nb.baseNode
	case datamodel.Kind_Int:
		return nb.baseNode
	case datamodel.Kind_Float:
		return nb.baseNode
	case datamodel.Kind_String:
		return nb.baseNode
	case datamodel.Kind_Bytes:
		return nb.baseNode
	case datamodel.Kind_Link:
		return nb.baseNode
	case 99:
		return nb.baseNode
	default:
		panic("unreachable")
	}
}

// -- NodeAmender -->

func (nb *anyBuilder) Transform(path datamodel.Path, transform datamodel.AmendFn) (datamodel.Node, error) {
	// If `baseNode` is set and supports amendment, apply the transformation. If it doesn't, and the root is being
	// replaced, replace it. If the transformation is for a nested node in a non-amendable recursive object, panic.
	if nb.amender != nil {
		return nb.amender.Transform(path, transform)
	}
	// `Transform` should never be called for a non-amendable node
	panic("misuse")
}

// -- NodeAssembler -->

// ... oddly enough, we seem to be able to put off implementing this
//  until we also implement something that goes full-hog on amortization
//   and actually has a slab of `anyNode`.  Which so far, nothing does.
//    See "REVIEW" comment on anyNode.
// type anyAssembler struct {
// 	w *anyNode
// }
