package traversal

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func Traverse(n ipld.Node, s selector.Selector, fn VisitFn) error {
	return TraversalProgress{}.Traverse(n, s, fn)
}

func TraverseInformatively(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	return TraversalProgress{}.TraverseInformatively(n, s, fn)
}

func TraverseTransform(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	return TraversalProgress{}.TraverseTransform(n, s, fn)
}

func (tp TraversalProgress) Traverse(n ipld.Node, s selector.Selector, fn VisitFn) error {
	tp.init()
	return tp.traverseInformatively(n, s, func(tp TraversalProgress, n ipld.Node, tr TraversalReason) error {
		if tr != TraversalReason_SelectionMatch {
			return nil
		}
		return fn(tp, n)
	})
}

func (tp TraversalProgress) TraverseInformatively(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	tp.init()
	return tp.traverseInformatively(n, s, fn)
}

func (tp TraversalProgress) traverseInformatively(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	if s.Decide(n) {
		if err := fn(tp, n, TraversalReason_SelectionMatch); err != nil {
			return err
		}
	} else {
		if err := fn(tp, n, TraversalReason_SelectionCandidate); err != nil {
			return err
		}
	}
	nk := n.ReprKind()
	switch nk {
	case ipld.ReprKind_Map, ipld.ReprKind_List: // continue
	default:
		return nil
	}
	// TODO: should only do this full loop if high-cardinality indicated.
	//   attn := s.Interests()
	//   if attn == nil {
	// FIXME need another kind switch here, and list support!
	for itr := n.MapIterator(); !itr.Done(); {
		k, v, err := itr.Next()
		if err != nil {
			return err
		}
		kstr, _ := k.AsString()
		sNext := s.Explore(n, selector.PathSegmentString{kstr})
		if sNext != nil {
			// TODO when link load is implemented, it should go roughly here.
			tpNext := tp
			tpNext.Path = tp.Path.AppendSegment(kstr)
			if err := tpNext.traverseInformatively(v, sNext, fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tp TraversalProgress) TraverseTransform(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	panic("TODO")
}
