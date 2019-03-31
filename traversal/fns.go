package traversal

import (
	"context"

	ipld "github.com/ipld/go-ipld-prime"
)

// This file defines interfaces for things users provide.
//------------------------------------------------

// VisitFn is a read-only visitor.
type VisitFn func(TraversalProgress, ipld.Node) error

// TransformFn is like a visitor that can also return a new Node to replace the visited one.
type TransformFn func(TraversalProgress, ipld.Node) (ipld.Node, error)

// AdvVisitFn is like VisitFn, but for use with AdvTraversal: it gets additional arguments describing *why* this node is visited.
type AdvVisitFn func(TraversalProgress, ipld.Node, TraversalReason) error

// TraversalReason provides additional information to traversals using AdvVisitFn.
type TraversalReason byte // enum = SelectionMatch | SelectionParent | SelectionCandidate // probably only pointful for block edges?

type TraversalProgress struct {
	*TraversalConfig
	Path      ipld.Path // Path is how we reached the current point in the traversal.
	LastBlock struct {  // LastBlock stores the Path and Link of the last block edge we had to load.  (It will always be zero in traversals with no linkloader.)
		ipld.Path
		ipld.Link
	}
}

type TraversalConfig struct {
	Ctx                    context.Context    // Context carried through a traversal.  Optional; use it if you need cancellation.
	LinkLoader             ipld.Loader        // Loader used for automatic link traversal.
	LinkNodeBuilderChooser NodeBuilderChooser // Chooser for Node implementations to produce during automatic link traversal.
	LinkStorer             ipld.Storer        // Storer used if any mutation features (e.g. traversal.Transform) are used.
}

// NodeBuilderChooser is a function that returns a NodeBuilder based on
// the information in a Link its LinkContext.
//
// A NodeBuilderChooser can be used in a TraversalConfig to be clear about
// what kind of Node implementation to use when loading a Link.
// In a simple example, it could constantly return an `ipldfree.NodeBuilder`.
// In a more complex example, a program using `bind` over native Go types
// could decide what kind of native type is expected, and return a
// `bind.NodeBuilder` for that specific concrete native type.
type NodeBuilderChooser func(ipld.Link, ipld.LinkContext) ipld.NodeBuilder
