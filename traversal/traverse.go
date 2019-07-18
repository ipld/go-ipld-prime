package traversal

import (
	"fmt"

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
	attn := s.Interests()
	if attn == nil {
		return tp.traverseAll(n, s, fn)
	}
	return tp.traverseSelective(n, attn, s, fn)

}

func (tp TraversalProgress) traverseAll(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	for itr := selector.NewSegmentIterator(n); !itr.Done(); {
		ps, v, err := itr.Next()
		if err != nil {
			return err
		}
		sNext := s.Explore(n, ps)
		if sNext != nil {
			tpNext := tp
			tpNext.Path = tp.Path.AppendSegment(ps.String())
			if v.ReprKind() == ipld.ReprKind_Link {
				v, err = tpNext.loadLink(v, n)
				if err != nil {
					return err
				}
			}

			err = tpNext.traverseInformatively(v, sNext, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tp TraversalProgress) traverseSelective(n ipld.Node, attn []selector.PathSegment, s selector.Selector, fn AdvVisitFn) error {
	for _, ps := range attn {
		// TODO: Probably not the most efficient way to be doing this...
		v, err := n.TraverseField(ps.String())
		if err != nil {
			continue
		}
		sNext := s.Explore(n, ps)
		if sNext != nil {
			tpNext := tp
			tpNext.Path = tp.Path.AppendSegment(ps.String())
			if v.ReprKind() == ipld.ReprKind_Link {
				v, err = tpNext.loadLink(v, n)
				if err != nil {
					return err
				}
			}

			err = tpNext.traverseInformatively(v, sNext, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tp TraversalProgress) loadLink(v ipld.Node, parent ipld.Node) (ipld.Node, error) {
	lnk, err := v.AsLink()
	if err != nil {
		return nil, err
	}
	// Assemble the LinkContext in case the Loader or NBChooser want it.
	lnkCtx := ipld.LinkContext{
		LinkPath:   tp.Path,
		LinkNode:   v,
		ParentNode: parent,
	}
	// Load link!
	v, err = lnk.Load(
		tp.Cfg.Ctx,
		lnkCtx,
		tp.Cfg.LinkNodeBuilderChooser(lnk, lnkCtx),
		tp.Cfg.LinkLoader,
	)
	if err != nil {
		return nil, fmt.Errorf("error traversing node at %q: could not load link %q: %s", tp.Path, lnk, err)
	}
	return v, nil
}

func (tp TraversalProgress) TraverseTransform(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	panic("TODO")
}
