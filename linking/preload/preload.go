package preload

import (
	"context"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
)

type Context struct {

	// Ctx is the familiar golang Context pattern.
	// Use this for cancellation, or attaching additional info
	// (for example, perhaps to pass auth tokens through to the storage functions).
	Ctx context.Context

	// Path where the link was encountered.  May be zero.
	//
	// Functions in the traversal package will set this automatically.
	BasePath datamodel.Path

	// Parent of the LinkNode.  May be zero.
	//
	// Functions in the traversal package will set this automatically.
	ParentNode datamodel.Node
}

type Link struct {
	Segment  datamodel.PathSegment
	LinkNode datamodel.Node
	Link     ipld.Link
}

type Loader func(Context, []Link)
