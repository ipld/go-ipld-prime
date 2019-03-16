package traversal

import (
	"fmt"
	"strconv"

	ipld "github.com/ipld/go-ipld-prime"
)

// Focus is a shortcut for kicking off
// TraversalProgress.Focus with an empty initial state
// (e.g. the Node given here is the "root" node of your operation).
func Focus(n ipld.Node, p ipld.Path, fn VisitFn) error {
	return TraversalProgress{}.Focus(n, p, fn)
}

// FocusedTransform is a shortcut for kicking off
// TraversalProgress.FocusedTransform with an empty initial state
// (e.g. the Node given here is the "root" node of your operation).
func FocusedTransform(n ipld.Node, p ipld.Path, fn TransformFn) (ipld.Node, error) {
	return TraversalProgress{}.FocusedTransform(n, p, fn)
}

// Focus traverses an ipld.Node graph, reaches a single Node,
// and applies a function to the reached node.
//
// Focus is a read-only traversal.
// See FocusedTransform if looking for a way to do an "update" to a Node.
//
// Focus can be used again again inside the applied VisitFn!
// By using the TraversalProgress handed to the VisitFn, the traversal Path
// so far will continue to be extended, so continued nested uses of Focus
// will see a fully contextualized Path.
func (tp TraversalProgress) Focus(n ipld.Node, p ipld.Path, fn VisitFn) error {
	segments := p.Segments()
	for i, seg := range segments {
		switch n.Kind() {
		case ipld.ReprKind_Invalid:
			return fmt.Errorf("cannot traverse node at %q: it is undefined", p.Truncate(i))
		case ipld.ReprKind_Map:
			next, err := n.TraverseField(seg)
			if err != nil {
				return fmt.Errorf("error traversing node at %q: %s", p.Truncate(i), err)
			}
			n = next
		case ipld.ReprKind_List:
			intSeg, err := strconv.Atoi(seg)
			if err != nil {
				return fmt.Errorf("cannot traverse node at %q: the next path segment (%q) cannot be parsed as a number and the node is a list", p.Truncate(i), seg)
			}
			next, err := n.TraverseIndex(intSeg)
			if err != nil {
				return fmt.Errorf("error traversing node at %q: %s", p.Truncate(i), err)
			}
			n = next
		case ipld.ReprKind_Link:
			panic("NYI link loading") // TODO
			// this would set a progress marker in `tp` as well
		default:
			return fmt.Errorf("error traversing node at %q: %s", p.Truncate(i), fmt.Errorf("cannot traverse terminals"))
		}
	}
	tp.Path = tp.Path.Join(p)
	return fn(tp, n)
}

// FocusedTransform traverses an ipld.Node graph, reaches a single Node,
// and applies a function to the reached node which make return a new Node.
//
// If the TransformFn returns a Node which is the same as the original
// reached node, the transform is a no-op, and the Node returned from the
// FocusedTransform call as a whole will also be the same as its starting Node.
//
// Otherwise, the reached node will be "replaced" with the new Node -- meaning
// that new intermediate nodes will be constructed to also replace each
// parent Node that was traversed to get here, thus propagating the changes in
// a copy-on-write fashion -- and the FocusedTransform function as a whole will
// return a new Node containing identical children except for those replaced.
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
//
// Note that anything you can do with the Transform function, you can also
// do with regular Node and NodeBuilder usage directly.  Transform just
// does a large amount of the intermediate bookkeeping that's useful when
// creating new values which are partial updates to existing values.
func (tp TraversalProgress) FocusedTransform(n ipld.Node, p ipld.Path, fn TransformFn) (ipld.Node, error) {
	panic("TODO") // TODO surprisingly different from Focus -- need to store nodes we traversed, and able do building.
}
