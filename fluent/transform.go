package fluent

import (
	ipld "github.com/ipld/go-ipld-prime"
)

// Transform traverses an ipld.Node graph and applies a function
// to the reached node.
//
// The applicative function must return a new Node; if the returned value is
// not equal to the original reached node, the reached node will be replaced
// with the new Node, and the Transform function as a whole will return a new
// node which is in the comparable graph position to the node the Transform
// was originally launched from.
// If this update takes place deep in the graph, new intermediate nodes will
// be constructed as necessary to propagate the changes in a copy-on-write
// fashion.
// (If the returned value is identical to the original reached node, there
// is no update; and the final value returned from Transform will also be
// identical to the starting value.)
//
// Transform can be used again inside the applicative function!
// This kind of composition can be useful for doing batches of updates.
// E.g. if have a large Node graph which contains a 100-element list, and
// you want to replace elements 12, 32, and 95 of that list:
// then you should Transform to the list first, and inside that applicative
// function's body, you can replace the entire list with a new one
// that is composed of copies of everything but those elements -- including
// using more Transform calls as desired to produce the replacement elements
// if it so happens that those replacement elements are easiest to construct
// by regarding them as incremental updates to the previous values.
//
// Transform will panic with a fluent.Error if any intermediate operations
// error.  (We are in the fluent package, after all.)
func Transform(
	node ipld.Node,
	path ipld.Path,
	applicative func(targetNode ipld.Node) (targetReplacement ipld.Node),
) (nodeReplacement ipld.Node) {
	return TransformUsingTraversal(node, path.Traverse, applicative)
}

// TransformUsingTraversal is identical to Transform, but accepts a generic
// ipld.Traversal function instead of an ipld.Path for guiding its selection
// of a node to transform.
func TransformUsingTraversal(
	node ipld.Node,
	traversal ipld.Traversal,
	applicative func(targetNode ipld.Node) (targetReplacement ipld.Node),
) (nodeReplacement ipld.Node) {
	return nil // TODO
}
