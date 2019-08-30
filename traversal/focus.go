package traversal

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
)

// Focus is a shortcut for kicking off
// traversal.Progress.Focus with an empty initial state
// (e.g. the Node given here is the "root" node of your operation).
func Focus(n ipld.Node, p ipld.Path, fn VisitFn) error {
	return Progress{}.Focus(n, p, fn)
}

// FocusedTransform is a shortcut for kicking off
// traversal.Progress.FocusedTransform with an empty initial state
// (e.g. the Node given here is the "root" node of your operation).
func FocusedTransform(n ipld.Node, p ipld.Path, fn TransformFn) (ipld.Node, error) {
	return Progress{}.FocusedTransform(n, p, fn)
}

// Focus traverses an ipld.Node graph, reaches a single Node,
// and applies a function to the reached node.
//
// Focus is a read-only traversal.
// See FocusedTransform if looking for a way to do an "update" to a Node.
//
// Focus can be used again again inside the applied VisitFn!
// By using the traversal.Progress handed to the VisitFn, the Path recorded
// of the traversal so far will continue to be extended, and thus continued
// nested uses of Focus will see a fully contextualized Path.
func (prog Progress) Focus(n ipld.Node, p ipld.Path, fn VisitFn) error {
	prog.init()
	segments := p.Segments()
	var prev ipld.Node // for LinkContext
	for i, seg := range segments {
		// Traverse the segment.
		switch n.ReprKind() {
		case ipld.ReprKind_Invalid:
			return fmt.Errorf("cannot traverse node at %q: it is undefined", p.Truncate(i))
		case ipld.ReprKind_Map:
			next, err := n.LookupString(seg.String())
			if err != nil {
				return fmt.Errorf("error traversing segment %q on node at %q: %s", seg, p.Truncate(i), err)
			}
			prev, n = n, next
		case ipld.ReprKind_List:
			intSeg, err := seg.Index()
			if err != nil {
				return fmt.Errorf("error traversing segment %q on node at %q: the segment cannot be parsed as a number and the node is a list", seg, p.Truncate(i))
			}
			next, err := n.LookupIndex(intSeg)
			if err != nil {
				return fmt.Errorf("error traversing segment %q on node at %q: %s", seg, p.Truncate(i), err)
			}
			prev, n = n, next
		default:
			return fmt.Errorf("cannot traverse node at %q: %s", p.Truncate(i), fmt.Errorf("cannot traverse terminals"))
		}
		// Dereference any links.
		for n.ReprKind() == ipld.ReprKind_Link {
			lnk, _ := n.AsLink()
			// Assemble the LinkContext in case the Loader or NBChooser want it.
			lnkCtx := ipld.LinkContext{
				LinkPath:   p.Truncate(i),
				LinkNode:   n,
				ParentNode: prev,
			}
			// Load link!
			next, err := lnk.Load(
				prog.Cfg.Ctx,
				lnkCtx,
				prog.Cfg.LinkNodeBuilderChooser(lnk, lnkCtx),
				prog.Cfg.LinkLoader,
			)
			if err != nil {
				return fmt.Errorf("error traversing node at %q: could not load link %q: %s", p.Truncate(i+1), lnk, err)
			}
			prog.LastBlock.Path = p.Truncate(i + 1)
			prog.LastBlock.Link = lnk
			prev, n = n, next
		}
	}
	prog.Path = prog.Path.Join(p)
	return fn(prog, n)
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
func (prog Progress) FocusedTransform(n ipld.Node, p ipld.Path, fn TransformFn) (ipld.Node, error) {
	panic("TODO") // TODO surprisingly different from Focus -- need to store nodes we traversed, and able do building.
}
