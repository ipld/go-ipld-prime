package traversal

import "github.com/ipld/go-ipld-prime/datamodel"

// Focus traverses a Node graph according to a path, reaches a single Node,
// and calls the given VisitFn on that reached node.
//
// This function is a helper function which starts a new traversal with default configuration.
// It cannot cross links automatically (since this requires configuration).
// Use the equivalent Focus function on the Progress structure
// for more advanced and configurable walks.
func Focus(n datamodel.Node, p datamodel.Path, fn VisitFn) error {
	return Progress{}.Focus(n, p, fn)
}

// Get is the equivalent of Focus, but returns the reached node (rather than invoking a callback at the target),
// and does not yield Progress information.
//
// This function is a helper function which starts a new traversal with default configuration.
// It cannot cross links automatically (since this requires configuration).
// Use the equivalent Get function on the Progress structure
// for more advanced and configurable walks.
func Get(n datamodel.Node, p datamodel.Path) (datamodel.Node, error) {
	return Progress{}.Get(n, p)
}

// FocusedTransform traverses a datamodel.Node graph, reaches a single Node,
// and calls the given TransformFn to decide what new node to replace the visited node with.
// A new Node tree will be returned (the original is unchanged).
//
// This function is a helper function which starts a new traversal with default configuration.
// It cannot cross links automatically (since this requires configuration).
// Use the equivalent FocusedTransform function on the Progress structure
// for more advanced and configurable walks.
func FocusedTransform(n datamodel.Node, p datamodel.Path, fn TransformFn, createParents bool) (datamodel.Node, error) {
	return Progress{}.FocusedTransform(n, p, fn, createParents)
}

// Focus traverses a Node graph according to a path, reaches a single Node,
// and calls the given VisitFn on that reached node.
//
// Focus is a read-only traversal.
// See FocusedTransform if looking for a way to do an "update" to a Node.
//
// Provide configuration to this process using the Config field in the Progress object.
//
// This walk will automatically cross links, but requires some configuration
// with link loading functions to do so.
//
// Focus (and the other traversal functions) can be used again again inside the VisitFn!
// By using the traversal.Progress handed to the VisitFn,
// the Path recorded of the traversal so far will continue to be extended,
// and thus continued nested uses of Walk and Focus will see the fully contextualized Path.
func (prog Progress) Focus(n datamodel.Node, p datamodel.Path, fn VisitFn) error {
	n, err := prog.get(n, p, true)
	if err != nil {
		return err
	}
	return fn(prog, n)
}

// Get is the equivalent of Focus, but returns the reached node (rather than invoking a callback at the target),
// and does not yield Progress information.
//
// Provide configuration to this process using the Config field in the Progress object.
//
// This walk will automatically cross links, but requires some configuration
// with link loading functions to do so.
//
// If doing several traversals which are nested, consider using the Focus funcion in preference to Get;
// the Focus functions provide updated Progress objects which can be used to do nested traversals while keeping consistent track of progress,
// such that continued nested uses of Walk or Focus or Get will see the fully contextualized Path.
func (prog Progress) Get(n datamodel.Node, p datamodel.Path) (datamodel.Node, error) {
	return prog.get(n, p, false)
}

// get is the internal implementation for Focus and Get.
// It *mutates* the Progress object it's called on, and returns reached nodes.
// For Get calls, trackProgress=false, which avoids some allocations for state tracking that's not needed by that call.
func (prog *Progress) get(n datamodel.Node, p datamodel.Path, trackProgress bool) (datamodel.Node, error) {
	prog.init()
	return NewAmender(n).Fetch(prog, p, trackProgress)
}

// FocusedTransform traverses a datamodel.Node graph, reaches a single Node,
// and calls the given TransformFn to decide what new node to replace the visited node with.
// A new Node tree will be returned (the original is unchanged).
//
// If the TransformFn returns the same Node which it was called with,
// then the transform is a no-op, and the Node returned from the
// FocusedTransform call as a whole will also be the same as its starting Node.
//
// Otherwise, the reached node will be "replaced" with the new Node -- meaning
// that new intermediate nodes will be constructed to also replace each
// parent Node that was traversed to get here, thus propagating the changes in
// a copy-on-write fashion -- and the FocusedTransform function as a whole will
// return a new Node containing identical children except for those replaced.
//
// Returning nil from the TransformFn as the replacement node means "remove this".
//
// FocusedTransform can be used again inside the applied function!
// This kind of composition can be useful for doing batches of updates.
// E.g. if have a large Node graph which contains a 100-element list, and
// you want to replace elements 12, 32, and 95 of that list:
// then you should FocusedTransform to the list first, and inside that
// TransformFn's body, you can replace the entire list with a new one
// that is composed of copies of everything but those elements -- including
// using more TransformFn calls as desired to produce the replacement elements
// if it so happens that those replacement elements are easiest to construct
// by regarding them as incremental updates to the previous values.
// (This approach can also be used when doing other modifications like insertion
// or reordering -- which would otherwise be tricky to define, since
// each operation could change the meaning of subsequently used indexes.)
//
// As a special case, list appending is supported by using the path segment "-".
// (This is determined by the node it applies to -- if that path segment
// is applied to a map, it's just a regular map key of the string of dash.)
//
// Note that anything you can do with the Transform function, you can also
// do with regular Node and NodeBuilder usage directly.  Transform just
// does a large amount of the intermediate bookkeeping that's useful when
// creating new values which are partial updates to existing values.
//
func (prog Progress) FocusedTransform(n datamodel.Node, p datamodel.Path, fn TransformFn, createParents bool) (datamodel.Node, error) {
	prog.init()
	a := NewAmender(n)
	if _, err := a.Transform(&prog, p, fn, createParents); err != nil {
		return nil, err
	}
	return a.Build(), nil
}
