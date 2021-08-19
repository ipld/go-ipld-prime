package traversal

import (
	"context"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
)

// This file defines interfaces for things users provide,
//  plus a few of the parameters they'll need to receieve.
//--------------------------------------------------------

// VisitFn is a read-only visitor.
type VisitFn func(Progress, datamodel.Node) error

// TransformFn is like a visitor that can also return a new Node to replace the visited one.
type TransformFn func(Progress, datamodel.Node) (datamodel.Node, error)

// AdvVisitFn is like VisitFn, but for use with AdvTraversal: it gets additional arguments describing *why* this node is visited.
type AdvVisitFn func(Progress, datamodel.Node, VisitReason) error

// VisitReason provides additional information to traversals using AdvVisitFn.
type VisitReason byte

const (
	VisitReason_SelectionMatch     VisitReason = 'm' // Tells AdvVisitFn that this node was explicitly selected.  (This is the set of nodes that VisitFn is called for.)
	VisitReason_SelectionParent    VisitReason = 'p' // Tells AdvVisitFn that this node is a parent of one that will be explicitly selected.  (These calls only happen if the feature is enabled -- enabling parent detection requires a different algorithm and adds some overhead.)
	VisitReason_SelectionCandidate VisitReason = 'x' // Tells AdvVisitFn that this node was visited while searching for selection matches.  It is not necessarily implied that any explicit match will be a child of this node; only that we had to consider it.  (Merkle-proofs generally need to include any node in this group.)
)

type Progress struct {
	Cfg       *Config
	Path      datamodel.Path // Path is how we reached the current point in the traversal.
	LastBlock struct {       // LastBlock stores the Path and Link of the last block edge we had to load.  (It will always be zero in traversals with no linkloader.)
		Path datamodel.Path
		Link datamodel.Link
	}
}

type Config struct {
	Ctx                            context.Context                // Context carried through a traversal.  Optional; use it if you need cancellation.
	LinkSystem                     linking.LinkSystem             // LinkSystem used for automatic link loading, and also any storing if mutation features (e.g. traversal.Transform) are used.
	LinkTargetNodePrototypeChooser LinkTargetNodePrototypeChooser // Chooser for Node implementations to produce during automatic link traversal.
}

// LinkTargetNodePrototypeChooser is a function that returns a NodePrototype based on
// the information in a Link and/or its LinkContext.
//
// A LinkTargetNodePrototypeChooser can be used in a traversal.Config to be clear about
// what kind of Node implementation to use when loading a Link.
// In a simple example, it could constantly return a `basicnode.Prototype.Any`.
// In a more complex example, a program using `bind` over native Go types
// could decide what kind of native type is expected, and return a
// `bind.NodeBuilder` for that specific concrete native type.
type LinkTargetNodePrototypeChooser func(datamodel.Link, linking.LinkContext) (datamodel.NodePrototype, error)

// SkipMe is a signalling "error" which can be used to tell traverse to skip some data.
//
// SkipMe can be returned by the Config.LinkLoader to skip entire blocks without aborting the walk.
// (This can be useful if you know you don't have data on hand,
// but want to continue the walk in other areas anyway;
// or, if you're doing a way where you know that it's valid to memoize seen
// areas based on Link alone.)
type SkipMe struct{}

func (SkipMe) Error() string {
	return "skip"
}
