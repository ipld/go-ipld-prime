package traversal

import (
	"context"

	ipld "github.com/ipld/go-ipld-prime"
)

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
	Ctx        context.Context // Context carried through a traversal.  Optional; use it if you need cancellation.
	LinkLoader ipld.Loader     // Loader used for automatic link traversal.
	LinkStorer ipld.Storer     // Storer used if any mutation features (e.g. traversal.Transform) are used.
}
