package traversal

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func WalkMatching(n ipld.Node, s selector.Selector, fn VisitFn) error {
	return Progress{}.WalkMatching(n, s, fn)
}

func WalkAdv(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	return Progress{}.WalkAdv(n, s, fn)
}

func WalkTransforming(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	return Progress{}.WalkTransforming(n, s, fn)
}

func (tp Progress) WalkMatching(n ipld.Node, s selector.Selector, fn VisitFn) error {
	tp.init()
	return tp.walkAdv(n, s, func(tp Progress, n ipld.Node, tr VisitReason) error {
		if tr != VisitReason_SelectionMatch {
			return nil
		}
		return fn(tp, n)
	})
}

func (tp Progress) WalkAdv(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	tp.init()
	return tp.walkAdv(n, s, fn)
}

func (tp Progress) walkAdv(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	if s.Decide(n) {
		if err := fn(tp, n, VisitReason_SelectionMatch); err != nil {
			return err
		}
	} else {
		if err := fn(tp, n, VisitReason_SelectionCandidate); err != nil {
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
		return tp.walkAdv_iterateAll(n, s, fn)
	}
	return tp.walkAdv_iterateSelective(n, attn, s, fn)

}

func (tp Progress) walkAdv_iterateAll(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
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
				lnk, _ := v.AsLink()
				tpNext.LastBlock.Path = tpNext.Path
				tpNext.LastBlock.Link = lnk
				v, err = tpNext.loadLink(v, n)
				if err != nil {
					return err
				}
			}

			err = tpNext.walkAdv(v, sNext, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tp Progress) walkAdv_iterateSelective(n ipld.Node, attn []selector.PathSegment, s selector.Selector, fn AdvVisitFn) error {
	for _, ps := range attn {
		// TODO: Probably not the most efficient way to be doing this...
		v, err := n.LookupString(ps.String())
		if err != nil {
			continue
		}
		sNext := s.Explore(n, ps)
		if sNext != nil {
			tpNext := tp
			tpNext.Path = tp.Path.AppendSegment(ps.String())
			if v.ReprKind() == ipld.ReprKind_Link {
				lnk, _ := v.AsLink()
				tpNext.LastBlock.Path = tpNext.Path
				tpNext.LastBlock.Link = lnk
				v, err = tpNext.loadLink(v, n)
				if err != nil {
					return err
				}
			}

			err = tpNext.walkAdv(v, sNext, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tp Progress) loadLink(v ipld.Node, parent ipld.Node) (ipld.Node, error) {
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

func (tp Progress) WalkTransforming(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	panic("TODO")
}
