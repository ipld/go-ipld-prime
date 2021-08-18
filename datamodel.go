package ipld

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

type (
	Kind          = datamodel.Kind
	Node          = datamodel.Node
	NodeAssembler = datamodel.NodeAssembler
	NodeBuilder   = datamodel.NodeBuilder
	NodePrototype = datamodel.NodePrototype
	MapIterator   = datamodel.MapIterator
	MapAssembler  = datamodel.MapAssembler
	ListIterator  = datamodel.ListIterator
	ListAssembler = datamodel.ListAssembler

	Link          = datamodel.Link
	LinkPrototype = datamodel.LinkPrototype

	Path        = datamodel.Path
	PathSegment = datamodel.PathSegment
)
