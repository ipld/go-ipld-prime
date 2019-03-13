package traversal

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
)

// An example use of LinkContext in LinkLoader logic might be inspecting the
// LinkNode, and if it's using the type system, inspecting its Type property;
// then deciding on whether or not we want to load objects of that Type.
// This might be used to do a traversal which looks at all directory objects,
// but not file contents, for example.
type LinkContext struct {
	LinkPath   Path      // whoops, nice cycle
	LinkNode   ipld.Node // has the cid again, but also might have type info // always zero for writing new nodes, for obvi reasons.
	ParentNode ipld.Node
}

type LinkLoader func(context.Context, cid.Cid, LinkContext) (ipld.Node, error)

// One presumes implementations will also *save* the content somewhere.
// The LinkContext parameter is a nod to this -- none of those parameters
// are relevant to the generation of the Cid itself, but perhaps one might
// want to use them when deciding where to store some files, etc.
type LinkBuilder func(
	ctx context.Context,
	node ipld.Node, lnkCtx LinkContext,
	multicodecType uint64, multihashType uint64, multihashLength int,
) (cid.Cid, error)

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
	Path      Path     // Path is how we reached the current point in the traversal.
	LastBlock struct { // LastBlock stores the Path and CID of the last block edge we had to load.  (It will always be zero in traversals with no linkloader.)
		Path
		cid.Cid
	}
}

type TraversalConfig struct {
	Ctx        context.Context // Context carried through a traversal.  Optional; use it if you need cancellation.
	LinkLoader LinkLoader
	// `blockWriter func(Context, Node, multicodec(?)) (CID, error)` probably belongs here.
}
