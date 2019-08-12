package traversal

import (
	"context"

	ipld "github.com/ipld/go-ipld-prime"
)

// This file defines interfaces for things users provide,
//  plus a few of the parameters they'll need to receieve.
//--------------------------------------------------------

// VisitFn is a read-only visitor.
type VisitFn func(Progress, ipld.Node) error

// TransformFn is like a visitor that can also return a new Node to replace the visited one.
type TransformFn func(Progress, ipld.Node) (ipld.Node, error)

// AdvVisitFn is like VisitFn, but for use with AdvTraversal: it gets additional arguments describing *why* this node is visited.
type AdvVisitFn func(Progress, ipld.Node, VisitReason) error

// VisitReason provides additional information to traversals using AdvVisitFn.
type VisitReason byte

const (
	VisitReason_SelectionMatch     VisitReason = 'm' // Tells AdvVisitFn that this node was explicitly selected.  (This is the set of nodes that VisitFn is called for.)
	VisitReason_SelectionParent    VisitReason = 'p' // Tells AdvVisitFn that this node is a parent of one that will be explicitly selected.  (These calls only happen if the feature is enabled -- enabling parent detection requires a different algorithm and adds some overhead.)
	VisitReason_SelectionCandidate VisitReason = 'x' // Tells AdvVisitFn that this node was visited while searching for selection matches.  It is not necessarily implied that any explicit match will be a child of this node; only that we had to consider it.  (Merkle-proofs generally need to include any node in this group.)
)

type Progress struct {
	Cfg       *Config
	Path      ipld.Path // Path is how we reached the current point in the traversal.
	LastBlock struct {  // LastBlock stores the Path and Link of the last block edge we had to load.  (It will always be zero in traversals with no linkloader.)
		Path ipld.Path
		Link ipld.Link
	}
}

type Config struct {
	Ctx                    context.Context    // Context carried through a traversal.  Optional; use it if you need cancellation.
	LinkLoader             ipld.Loader        // Loader used for automatic link traversal.
	LinkNodeBuilderChooser NodeBuilderChooser // Chooser for Node implementations to produce during automatic link traversal.
	LinkStorer             ipld.Storer        // Storer used if any mutation features (e.g. traversal.Transform) are used.
}

// NodeBuilderChooser is a function that returns a NodeBuilder based on
// the information in a Link its LinkContext.
//
// A NodeBuilderChooser can be used in a traversal.Config to be clear about
// what kind of Node implementation to use when loading a Link.
// In a simple example, it could constantly return an `ipldfree.NodeBuilder`.
// In a more complex example, a program using `bind` over native Go types
// could decide what kind of native type is expected, and return a
// `bind.NodeBuilder` for that specific concrete native type.
type NodeBuilderChooser func(ipld.Link, ipld.LinkContext) ipld.NodeBuilder
