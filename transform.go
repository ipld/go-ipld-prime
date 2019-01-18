package ipld

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
// Note that anything you can do with the Transform function, you can also
// do with regular Node and NodeBuilder usage directly.  Transform just
// does a large amount of the intermediate bookkeeping that's useful when
// creating new values which are partial updates to existing values.
func Transform(
	node Node,
	path Path,
	applicative func(reachedNode Node, reachedPath Path) (reachedNodeReplacement Node),
) (nodeReplacement Node, err error) {
	return TransformUsingTraversal(node, path.Traverse, applicative)
}

// TransformUsingTraversal is identical to Transform, but accepts a generic
// ipld.Traversal function instead of an ipld.Path for guiding its selection
// of a node to transform.
func TransformUsingTraversal(
	node Node,
	traversal Traversal,
	applicative func(reachedNode Node, reachedPath Path) (reachedNodeReplacement Node),
) (nodeReplacement Node, err error) {
	panic("TODO") // TODO
}

// ContinueTransform is similar to Transform, but takes an additional parameter
// in order to keep Path information complete when doing nested Transforms.
//
// Use ContinueTransform in the body of the applicative function of a Transform
// (or ContinueTransform) call: providing the so-far reachedPath as the nodePath
// to the ContinueTransform call will make sure the next, deeper applicative
// call will get a reachedPath which is the complete path all the way from the
// root node of the outermost transform.
//
// (Or, ignore all this and use Transform nested bare.  It's your own Path
// information you're messing with; if you feel like your algorithm would
// work better seeing a locally scoped path rather than a more globally
// rooted one, that's absolutely fine.)
func FurtherTransform(
	node Node,
	nodePath Path,
	path Path,
	applicative func(reachedNode Node, reachedPath Path) (reachedNodeReplacement Node),
) (nodeReplacement Node, err error) {
	tfn := ContinuedTraversal(path.Traverse, nodePath)
	// REVIEW: so do we *need* ContinuedTraversal?  Much less exported?
	//  I don't have a huge argument against it other than that it seems this
	//   is the only place we'd use it, and it does *all* the work, surprisingly.
	return TransformUsingTraversal(node, tfn, applicative)
}

// Note that another way to write this would be using some sort of
// `StartTransform() TransformController` pattern...
// But the arity and placement of things that have to carry information
// is basically identical.  The only win would be keeping the number of
// parameters of type Path down to one instead of sometimes two.
// Is that worth it?  Maybe.  I'm unconvinced.
